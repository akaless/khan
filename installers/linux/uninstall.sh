#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════
# خان (Khan) — Uninstaller for Linux
# حذف: ./uninstall.sh
# ═══════════════════════════════════════════════════════════
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/khan"
DATA_DIR="/var/lib/khan"
BIN="/usr/local/bin/khan"

echo -e "${RED}⚠️ حذف خان...${NC}"

# ─── سرویس ───
if [ -f /etc/systemd/system/khan.service ]; then
  echo -e "${YELLOW}→ توقف و حذف سرویس...${NC}"
  systemctl stop khan 2>/dev/null || true
  systemctl disable khan 2>/dev/null || true
  rm -f /etc/systemd/system/khan.service
  systemctl daemon-reload
fi

# ─── فایل‌ها ───
echo -e "${YELLOW}→ حذف فایل‌ها...${NC}"
rm -f "$BIN"
rm -rf "$INSTALL_DIR"

# ─── داده (اختیاری) ───
if [ -d "$DATA_DIR" ]; then
  read -p "داده‌ها حذف شوند؟ (y/N): " ans
  if [[ "$ans" == "y" || "$ans" == "Y" ]]; then
    rm -rf "$DATA_DIR"
    echo -e "  ${RED}✅ داده‌ها حذف شد${NC}"
  else
    echo -e "  ${GREEN}داده‌ها حفظ شد: $DATA_DIR${NC}"
  fi
fi

echo -e "${GREEN}✅ خان حذف شد.${NC}"
