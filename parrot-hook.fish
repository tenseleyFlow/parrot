# Parrot shell hook for fish - source this in your fish config
# Typically installed to ~/.config/fish/conf.d/parrot.fish

# Path to parrot binary - update this if needed
set -g PARROT_BIN "parrot"

# Function to check if parrot binary exists
function parrot_check
    if not command -v $PARROT_BIN >/dev/null 2>&1
        echo "⚠️  Parrot binary not found. Make sure 'parrot' is in your PATH."
        return 1
    end
    return 0
end

# Function called after each command completes
function parrot_postexec --on-event fish_postexec
    set -l exit_code $status
    set -l last_cmd (history max 1)

    # Only mock if command failed and we have a command
    if test $exit_code -ne 0; and test -n "$last_cmd"; and parrot_check
        # Run parrot in background to avoid blocking shell if PARROT_ASYNC is set
        if test "$PARROT_ASYNC" = "true"
            $PARROT_BIN mock "$last_cmd" "$exit_code" &
        else
            $PARROT_BIN mock "$last_cmd" "$exit_code"
        end
    end
end

# Setup hook (only show messages if not already initialized and in interactive mode)
if not set -q PARROT_INITIALIZED
    # Only show messages in interactive shells
    if status is-interactive
        echo "🦜 Parrot is now watching your fish commands..."

        # Show performance tip
        if test "$PARROT_ASYNC" != "true"
            echo "💡 Tip: Set PARROT_ASYNC=true to prevent terminal hangs on slow networks"
        end
    end

    # Mark as initialized for this shell session
    set -gx PARROT_INITIALIZED 1
end
