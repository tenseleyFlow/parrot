# Parrot Release Guide

## Quick Release Flow

### 1. Automated Release Script
```bash
./scripts/release.sh <new-version> [old-version]
```

**Example:**
```bash
./scripts/release.sh 1.0.5          # Auto-detect current version  
./scripts/release.sh 1.0.5 1.0.4    # Explicit old version
```

### 2. Manual Steps After Script

#### A. Edit Changelog
Edit `parrot.spec` and replace `[ADD YOUR CHANGELOG ENTRIES HERE]` with actual changes:
```
* Mon Aug 26 2024 mfw <espadonne@outlook.com> - 1.0.5-1
- Fixed sanitization bug with character counts
- Improved timeout handling for better responsiveness
```

#### B. Deploy to RPM Repository
```bash
make deploy
# (Requires sudo access - will prompt)
```

#### C. Test RPM Update
```bash
sudo dnf update parrot
```

#### D. Push Git Changes
```bash
git push origin trunk v1.0.5
```

#### E. Update and Test AUR Package
```bash
cd /tmp/parrot-cli
makepkg -si                    # Test build
git add -A
git commit -m "Update to 1.0.5"
git push origin master
```

## Manual Release Flow (Alternative)

If you prefer to do it step by step:

### 1. Update Versions
```bash
# Update Makefile
sed -i 's/VERSION = 1.0.4/VERSION = 1.0.5/' Makefile

# Update RPM spec
sed -i 's/Version:        1.0.4/Version:        1.0.5/' parrot.spec
```

### 2. Update Changelog in parrot.spec
Add new entry at the top of the `%changelog` section:
```
%changelog
* Mon Aug 26 2024 mfw <espadonne@outlook.com> - 1.0.5-1
- Your changes here
- Another change here

* Previous entries...
```

### 3. Build and Test
```bash
make clean
make build
./parrot --version       # Should show new version
./parrot mock "test" "1" # Quick functionality test
```

### 4. Deploy to RPM Repository
```bash
make deploy
```

### 5. Update AUR Package
```bash
cd /tmp/parrot-cli

# Update PKGBUILD
sed -i 's/pkgver=1.0.4/pkgver=1.0.5/' PKGBUILD
sed -i 's/#tag=v1.0.4/#tag=v1.0.5/' PKGBUILD
sed -i 's/pkgrel=[0-9]*/pkgrel=1/' PKGBUILD

# Generate new .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Test build
makepkg -si
```

### 6. Git Operations
```bash
# In parrot project
git add -A
git commit -m "Bump version to 1.0.5"
git tag v1.0.5
git push origin trunk v1.0.5

# In AUR package
cd /tmp/parrot-cli
git add -A
git commit -m "Update to 1.0.5"  
git push origin master
```

## Repository Locations

- **Main Project**: `~/src/parrot`
- **RPM Repository**: `~/src/repos-musicsian-com`
- **AUR Package**: `/tmp/parrot-cli`

## Key Files to Update

### Main Project (`~/src/parrot`)
- `Makefile` - VERSION variable
- `parrot.spec` - Version and %changelog

### AUR Package (`/tmp/parrot-cli`)
- `PKGBUILD` - pkgver, pkgrel, source tag
- `.SRCINFO` - Generated automatically

## Testing Checklist

- [ ] Build completes without errors
- [ ] `./parrot --version` shows new version
- [ ] Basic functionality test: `./parrot mock "test" "1"`
- [ ] RPM packages build and sign correctly
- [ ] Repository deployment succeeds
- [ ] `sudo dnf update parrot` works without `--nogpgcheck`
- [ ] AUR package builds: `makepkg -si`
- [ ] Git tags and pushes complete

## Common Issues

### RPM Deployment Fails
- **Issue**: Deployment stops at sudo steps
- **Fix**: Complete manually with the commands shown in script output

### AUR Build Fails  
- **Issue**: Missing dependencies or build errors
- **Fix**: Check PKGBUILD dependencies and test locally

### GPG Verification Fails
- **Issue**: DNF requires `--nogpgcheck`
- **Fix**: Ensure `/etc/yum.repos.d/musicsian.repo` has `repo_gpgcheck=0`

### Version Mismatch
- **Issue**: Forgot to update a version somewhere
- **Fix**: Check all locations: Makefile, parrot.spec, PKGBUILD