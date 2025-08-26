# 🦜 Parrot
(noun) : birb  
_see also_ : birbb  

A sassy CLI tool that mocks your failed commands with intelligent insults.

## Quick Start

1. Build the binary:
   ```bash
   go build -o parrot main.go
   ```

2. Install shell hooks:
   ```bash
   ./parrot install
   source ~/.bashrc  # or ~/.zshrc
   ```

3. Watch parrot roast your failures!
   ```bash
   git push  # fails
   # 🦜 Did you forget to pull again? Classic amateur move.
   ```

## Manual Testing

Test parrot responses without shell hooks:
```bash
./parrot mock "git commit" "1"
./parrot mock "npm install" "1" 
./parrot mock "docker run" "125"
```

Coming up: TOML configuration, multiple personality profiles, and response variety!
