#!/usr/bin/env bash
# Complete release workflow for parrot
set -euo pipefail

# Configuration
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
AUR_DIR="/tmp/parrot-cli"
OLD_VERSION=""
NEW_VERSION=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}▶${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}❌${NC} $1"
}

# Parse command line arguments
usage() {
    echo "Usage: $0 <new-version> [old-version]"
    echo ""
    echo "Examples:"
    echo "  $0 1.0.5                    # Auto-detect current version"
    echo "  $0 1.0.5 1.0.4              # Explicit old version"
    echo ""
    echo "This script will:"
    echo "  1. Update version in Makefile and parrot.spec"
    echo "  2. Update changelog with new entry"
    echo "  3. Build and test the new version"
    echo "  4. Deploy to RPM repository (~/src/repos-musicsian-com)"
    echo "  5. Update AUR package (/tmp/parrot-cli)"
    echo "  6. Create git tag and commit"
    exit 1
}

if [ $# -eq 0 ]; then
    usage
fi

NEW_VERSION="$1"
if [ $# -gt 1 ]; then
    OLD_VERSION="$2"
else
    # Auto-detect current version from Makefile
    OLD_VERSION=$(grep "^VERSION = " "$PROJECT_DIR/Makefile" | cut -d' ' -f3)
fi

echo "🚀 Parrot Release Workflow"
echo "=========================="
echo "Old version: $OLD_VERSION"
echo "New version: $NEW_VERSION"
echo ""

# Confirm with user
read -p "Continue with release? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

cd "$PROJECT_DIR"

# Step 1: Update Makefile version
log "Updating Makefile version..."
sed -i "s/VERSION = $OLD_VERSION/VERSION = $NEW_VERSION/" Makefile
success "Updated Makefile: $OLD_VERSION → $NEW_VERSION"

# Step 2: Update RPM spec version
log "Updating parrot.spec version..."
sed -i "s/Version:        $OLD_VERSION/Version:        $NEW_VERSION/" parrot.spec
success "Updated parrot.spec: $OLD_VERSION → $NEW_VERSION"

# Step 3: Add changelog entry (user will need to edit)
log "Adding changelog entry to parrot.spec..."
DATE=$(date '+%a %b %d %Y')

# Create a temporary file with the changelog entry
cat > /tmp/changelog_entry << EOF
* $DATE mfw <espadonne@outlook.com> - $NEW_VERSION-1
- [ADD YOUR CHANGELOG ENTRIES HERE]

EOF

# Insert the changelog entry after %changelog line
sed -i '/^%changelog/r /tmp/changelog_entry' parrot.spec
rm -f /tmp/changelog_entry

warn "Please edit parrot.spec and update the changelog entry with actual changes!"

# Step 4: Build and test
log "Building new version..."
make clean
make build

log "Running smoke tests..."
./parrot --version
./parrot --help >/dev/null
./parrot mock "test command" "1" >/dev/null

success "Build and tests completed"

# Step 5: Deploy to RPM repository
log "Deploying to RPM repository..."
echo "You'll need to run 'make deploy' manually after this script completes."
echo "The deploy requires sudo access for the web server."

# Step 6: Update AUR package
if [ -d "$AUR_DIR" ]; then
    log "Updating AUR package..."
    cd "$AUR_DIR"
    
    # Update PKGBUILD version
    sed -i "s/pkgver=$OLD_VERSION/pkgver=$NEW_VERSION/" PKGBUILD
    
    # Update git tag reference
    sed -i "s/#tag=v$OLD_VERSION/#tag=v$NEW_VERSION/" PKGBUILD
    
    # Reset pkgrel to 1 for new version
    sed -i "s/pkgrel=[0-9]*/pkgrel=1/" PKGBUILD
    
    success "Updated AUR PKGBUILD: $OLD_VERSION → $NEW_VERSION"
    
    # Generate new .SRCINFO
    log "Generating .SRCINFO..."
    makepkg --printsrcinfo > .SRCINFO
    
    success "Generated new .SRCINFO"
    
    warn "Don't forget to:"
    warn "  1. Review the PKGBUILD changes"
    warn "  2. Test build with: makepkg -si"
    warn "  3. Commit and push to AUR"
else
    warn "AUR directory not found at $AUR_DIR"
fi

# Step 7: Git operations
cd "$PROJECT_DIR"
log "Git operations..."

# Check if we have uncommitted changes
if ! git diff --quiet || ! git diff --cached --quiet; then
    warn "You have uncommitted changes. Commit them first:"
    warn "  git add -A"
    warn "  git commit -m \"Bump version to $NEW_VERSION\""
    warn "  git tag v$NEW_VERSION"
    warn "  git push origin trunk v$NEW_VERSION"
else
    log "Creating git commit and tag..."
    git add -A
    git commit -m "Bump version to $NEW_VERSION"
    git tag "v$NEW_VERSION"
    
    warn "Don't forget to push:"
    warn "  git push origin trunk v$NEW_VERSION"
fi

echo ""
echo "🎉 Release workflow completed!"
echo ""
echo "Next steps:"
echo "1. Edit parrot.spec changelog entry"
echo "2. Run 'make deploy' (requires sudo)"
echo "3. Test the RPM repository update"
echo "4. Review and test AUR package in $AUR_DIR"
echo "5. Push git changes: git push origin trunk v$NEW_VERSION"
echo "6. Commit and push AUR changes"