#!/bin/bash

# Parrot shell hook - source this in your .bashrc or .zshrc

# Path to parrot binary - update this if needed
PARROT_BIN="parrot"

# Function to check if parrot binary exists
parrot_check() {
    if ! command -v "$PARROT_BIN" &> /dev/null; then
        echo "⚠️  Parrot binary not found. Make sure 'parrot' is in your PATH."
        return 1
    fi
    return 0
}

# Function called after each command in bash
parrot_prompt_command() {
    local exit_code=$?
    local last_cmd=$(history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
    
    # Only mock if command failed and we have a command
    if [ $exit_code -ne 0 ] && [ -n "$last_cmd" ] && parrot_check; then
        # Run parrot in background to avoid blocking shell if PARROT_ASYNC is set
        if [ "${PARROT_ASYNC:-}" = "true" ]; then
            "$PARROT_BIN" mock "$last_cmd" "$exit_code" &
        else
            "$PARROT_BIN" mock "$last_cmd" "$exit_code"
        fi
    fi
}

# Function called before each command in zsh
parrot_preexec() {
    PARROT_LAST_CMD="$1"
}

# Function called after each command in zsh
parrot_precmd() {
    local exit_code=$?
    
    # Only mock if command failed and we have a command
    if [ $exit_code -ne 0 ] && [ -n "$PARROT_LAST_CMD" ] && parrot_check; then
        # Run parrot in background to avoid blocking shell if PARROT_ASYNC is set
        if [ "${PARROT_ASYNC:-}" = "true" ]; then
            "$PARROT_BIN" mock "$PARROT_LAST_CMD" "$exit_code" &
        else
            "$PARROT_BIN" mock "$PARROT_LAST_CMD" "$exit_code"
        fi
    fi
}

# Setup based on shell type
if [ -n "$BASH_VERSION" ]; then
    # Bash setup
    PROMPT_COMMAND="parrot_prompt_command${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
    echo "🦜 Parrot is now watching your bash commands..."
elif [ -n "$ZSH_VERSION" ]; then
    # Zsh setup
    autoload -Uz add-zsh-hook
    add-zsh-hook preexec parrot_preexec
    add-zsh-hook precmd parrot_precmd
    echo "🦜 Parrot is now watching your zsh commands..."
else
    echo "⚠️  Parrot: Unsupported shell. Only bash and zsh are supported."
fi

# Show performance tip
if [ "${PARROT_ASYNC:-}" != "true" ]; then
    echo "💡 Tip: Set PARROT_ASYNC=true to prevent terminal hangs on slow networks"
fi