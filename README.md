# 🦜 Parrot
(noun) : birb  
_see also_ : birbb  

A sassy CLI tool that mocks your failed commands with intelligent insults.

## 🚀 Quick Start

### **Option 1: Complete Setup (Recommended)**
```bash
# Build parrot
go build -o parrot main.go

# Auto-guided setup
./parrot setup
# Interactive wizard walks you through:
# 1. Choose backend (API/Local/Fallback) 
# 2. Configure API keys or install models
# 3. Install shell hooks
# 4. Test your parrot!
```

### **Option 2: Manual Setup**
```bash
# Basic functionality (works immediately)
./parrot install              # Install shell hooks
source ~/.bashrc              # Restart shell

# Add intelligence later
./parrot configure            # Interactive config wizard
```

### **🎯 Result**
```bash
git push origin nonexistent   # Try failing a command
# 🦜 Git rejected your code harder than everyone rejects you.
```

## Manual Testing

Test parrot responses without shell hooks:
```bash
./parrot mock "git commit" "1"
./parrot mock "npm install" "1" 
./parrot mock "docker run" "125"
```

## Phase 2.5 Packaging Architecture - COMPLETE ✅

- ✅ **Multi-backend LLM support** (API primary, Ollama secondary, fallbacks)
- ✅ **Configuration system** with TOML and environment variable support
- ✅ **Standard installation paths** (RPM-friendly)
- ✅ **OpenAI-compatible API integration** as primary backend
- ✅ **Phi-3.5-mini model** integration for better local quality
- ✅ **Smart backend priority** (API → Local → Fallback)
- ✅ **Debug mode and status reporting**

## Commands Available

- `parrot mock "command" "exit_code"` - Test responses
- `parrot status` - Show backend status and configuration
- `parrot config init` - Create sample config file
- `parrot install` - Install shell hooks

## Phase 3: Personality & Polish - COMPLETE ✅

- ✅ **Three personality levels**: mild, sarcastic, savage
- ✅ **Personality-specific prompt templates** for better AI responses
- ✅ **Terminal colors and formatting** with personality-based themes
- ✅ **Enhanced output modes** with borders and emphasis
- ✅ **Environment variable overrides** (NO_COLOR, PARROT_CONFIG, etc.)
- ✅ **Demo command** to showcase personalities and colors

## 🎮 Available Commands

| Command | Description |
|---------|-------------|
| `parrot setup` | **🚀 Complete setup wizard** - guided installation |
| `parrot configure` | **⚙️ Interactive config** - customize all settings |
| `parrot install` | **🔗 Install shell hooks** - enable auto-roasting |
| `parrot status` | **📊 System status** - check backends & config |
| `parrot mock "cmd" "code"` | **🧪 Test responses** - try commands manually |
| `parrot demo` | **🎨 Personality showcase** - see all personalities |
| `parrot config init` | **📝 Create config file** - manual configuration |

## Configuration Examples

```bash
# Use custom config file
PARROT_CONFIG=./my-config.toml parrot mock "git push" "1"

# Override personality
PARROT_PERSONALITY=savage parrot mock "docker run" "125"  

# Disable colors
NO_COLOR=1 parrot mock "npm install" "1"

# Enable debug mode
PARROT_DEBUG=true parrot mock "curl api.com" "7"
```

## Next: Ready for Production!

Your parrot is now **feature-complete** with:
- **Intelligent AI responses** with 3 personality levels
- **Multiple backend support** (API → Local → Fallback)  
- **Beautiful terminal output** with colors and formatting
- **RPM-ready architecture** for easy packaging
- **Comprehensive configuration** system
