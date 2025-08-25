Name:           parrot
Version:        1.0.0
Release:        1%{?dist}
Summary:        Intelligent CLI command failure assistant with AI-powered responses

License:        MIT
URL:            https://github.com/tenseleyFlow/parrot
Source0:        %{name}-%{version}.tar.gz

BuildArch:      x86_64
BuildRequires:  golang >= 1.21
BuildRequires:  make
Requires:       bash
Suggests:       ollama
Suggests:       curl

%description
Parrot is an intelligent CLI assistant that listens for failed command executions
and provides witty, AI-powered responses with helpful suggestions. It supports
multiple backend modes including API integration (OpenAI-compatible), local
LLM models via Ollama, and fallback responses for guaranteed functionality.

Features:
- Multi-backend architecture (API → Local → Fallback)
- Three personality levels (mild, sarcastic, savage)
- Shell integration for bash and zsh
- Interactive configuration wizard
- Comprehensive setup automation
- Terminal color theming
- Zero external dependencies required for basic operation

%prep
%autosetup

%build
# Build Go binary with release optimizations
go mod download
go build -ldflags="-w -s" -o parrot .

%install
# Install main binary
install -d %{buildroot}%{_bindir}
install -m 755 parrot %{buildroot}%{_bindir}/parrot

# Install shell integration hooks
install -d %{buildroot}%{_datadir}/%{name}
install -m 644 parrot-hook.sh %{buildroot}%{_datadir}/%{name}/parrot-hook.sh

# Install configuration templates
install -d %{buildroot}%{_sysconfdir}/%{name}
install -m 644 config/parrot.toml.example %{buildroot}%{_sysconfdir}/%{name}/parrot.toml.example

# Install documentation
install -d %{buildroot}%{_docdir}/%{name}
install -m 644 README.md %{buildroot}%{_docdir}/%{name}/
install -m 644 INSTALLATION_FLOWS.md %{buildroot}%{_docdir}/%{name}/

%post
# Post-install setup guidance
echo "🦜 Parrot has been installed successfully!"
echo ""
echo "Next steps to complete setup:"
echo "1. Run 'parrot setup' to configure your preferred backend"
echo "2. Follow the interactive prompts to enable shell integration"
echo ""
echo "Backend options:"
echo "  • API Backend: Use OpenAI-compatible services (requires API key)"
echo "  • Local Backend: Use Ollama for privacy-focused local AI"
echo "  • Fallback: Built-in responses (no setup required)"
echo ""
echo "For detailed setup instructions, see: /usr/share/doc/%{name}/INSTALLATION_FLOWS.md"
echo "Example configuration: /etc/%{name}/parrot.toml.example"
echo ""
echo "Run 'parrot --help' to get started!"

%preun
# Clean up shell integrations on uninstall
if [ "$1" = "0" ]; then
    # Only on complete removal, not upgrade
    echo "Removing Parrot shell integrations..."
    # Note: Users should run 'parrot setup --remove' before uninstalling
    echo "If you have shell integration enabled, please restart your terminal sessions."
fi

%files
%doc README.md INSTALLATION_FLOWS.md
%{_bindir}/parrot
%{_datadir}/%{name}/parrot-hook.sh
%{_sysconfdir}/%{name}/parrot.toml.example
%{_docdir}/%{name}/

%changelog
* Sun Aug 25 2025 mfw <espadonne@outlook.com> - 1.0.0-1
- Initial RPM release
- Multi-backend architecture with API, Local, and Fallback support
- Interactive setup wizard with automated backend configuration
- Shell integration for bash and zsh
- Three personality levels with terminal color theming
- Comprehensive installation flows and setup automation