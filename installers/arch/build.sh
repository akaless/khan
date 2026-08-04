#!/bin/bash
set -e

echo "=== Building Khan v1.0.3 Arch package ==="

# Generate SHA256 checksums
echo "Generating checksums..."
KHAN_SHA=$(sha256sum khan-linux | cut -d' ' -f1)
SVC_SHA=$(sha256sum khan.service | cut -d' ' -f1)
CFG_SHA=$(sha256sum config.json | cut -d' ' -f1)
LIC_SHA=$(sha256sum LICENSE | cut -d' ' -f1)

echo "khan-linux: $KHAN_SHA"
echo "khan.service: $SVC_SHA"
echo "config.json: $CFG_SHA"
echo "LICENSE: $LIC_SHA"

# Update PKGBUILD with actual checksums (properly)
sed -i "s|'SKIP'|'$KHAN_SHA'|" PKGBUILD
sed -i "0,/'SKIP'/s||'$SVC_SHA'|" PKGBUILD
sed -i "0,/'SKIP'/s||'$CFG_SHA'|" PKGBUILD
sed -i "0,/'SKIP'/s||'$LIC_SHA'|" PKGBUILD

echo "Updated PKGBUILD:"
cat PKGBUILD

echo ""
echo "Package ready for makepkg on Arch system"
echo "Run: makepkg -sf"
