#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════
# خان (Khan) — macOS Installer script
# ساخت DMG: ./build-dmg.sh
# ═══════════════════════════════════════════════════════════
set -e

echo "🏠 ساخت DMG برای خان..."

# ─── پیش‌نیاز: Xcode Command Line Tools ───
if ! command -v hdiutil >/dev/null 2>&1; then
  echo "❌ hdiutil یافت نشد — Xcode Command Line Tools نصب کنید:"
  echo "   xcode-select --install"
  exit 1
fi

APP_NAME="خان"
APP_DIR="dist/Khan.app"
DMG_NAME="Khan-1.0.0-$(uname -m).dmg"
STAGING="dist/staging"

# ─── ساخت ساختار app ───
echo "→ ساخت Khan.app..."
mkdir -p "$APP_DIR/Contents/MacOS"
mkdir -p "$APP_DIR/Contents/Resources"

# باینری (arm64 یا x86_64 بسته به سیستم)
if [ "$(uname -m)" = "arm64" ]; then
  cp dist/khan-macos-arm64 "$APP_DIR/Contents/MacOS/khan"
else
  cp dist/khan-macos-intel "$APP_DIR/Contents/MacOS/khan"
fi
chmod +x "$APP_DIR/Contents/MacOS/khan"

# ─── Info.plist ───
cat > "$APP_DIR/Contents/Info.plist" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key><string>Khan</string>
  <key>CFBundleDisplayName</key><string>خان</string>
  <key>CFBundleIdentifier</key><string>ir.khanchat.app</string>
  <key>CFBundleVersion</key><string>1.0.0</string>
  <key>CFBundleShortVersionString</key><string>1.0.0</string>
  <key>CFBundleExecutable</key><string>khan</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>LSMinimumSystemVersion</key><string>12.0</string>
  <key>NSHighResolutionCapable</key><true/>
  <key>LSUIElement</key><false/>
</dict>
</plist>
EOF

# ─── ساخت DMG ───
echo "→ ساخت DMG..."
rm -rf "$STAGING"
mkdir -p "$STAGING"
cp -R "$APP_DIR" "$STAGING/"
ln -sf /Applications "$STAGING/Applications"

hdiutil create -volname "خان" -srcfolder "$STAGING" -ov -format UDZO "$DMG_NAME"

# ─── پاک‌سازی ───
rm -rf "$STAGING" "$APP_DIR"

echo "✅ DMG ساخته شد: $DMG_NAME"
echo "   برای نصب: DMG را باز کنید و خان را به Applications بکشید"
