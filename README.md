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

## Commands Available

- `parrot mock "command" "exit_code"` - Test responses
- `parrot status` - Show backend status and configuration
- `parrot config init` - Create sample config file
- `parrot install` - Install shell hooks

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
