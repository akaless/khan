#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════
# خان (Khan) — Installer for Linux (generic)
# نصب: ./install.sh [--service] [--port 1727]
# ═══════════════════════════════════════════════════════════
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/khan"
DATA_DIR="/var/lib/khan"
BIN="/usr/local/bin/khan"
PORT=1727
AS_SERVICE=false

# ─── پارس آرگومان‌ها ───
while [[ $# -gt 0 ]]; do
  case "$1" in
    --service) AS_SERVICE=true; shift ;;
    --port) PORT="$2"; shift 2 ;;
    *) echo "ناشناخته: $1"; exit 1 ;;
  esac
done

echo -e "${GREEN}🏠 نصب خان${NC}"

# ─── کپی باینری ───
echo -e "${YELLOW}→ کپی باینری...${NC}"
mkdir -p "$INSTALL_DIR" "$DATA_DIR"
cp khan-linux "$INSTALL_DIR/khan-linux"
chmod +x "$INSTALL_DIR/khan-linux"
ln -sf "$INSTALL_DIR/khan-linux" "$BIN"

# ─── فایروال ───
echo -e "${YELLOW}→ باز کردن پورت $PORT...${NC}"
if command -v ufw >/dev/null 2>&1; then
  ufw allow "$PORT/tcp" 2>/dev/null && echo -e "  ${GREEN}✅ ufw: پورت $PORT باز شد${NC}" || echo "  (ufw موجود نیست — دستی باز کنید)"
elif command -v firewall-cmd >/dev/null 2>&1; then
  firewall-cmd --permanent --add-port="$PORT/tcp" >/dev/null 2>&1 && firewall-cmd --reload >/dev/null 2>&1 && echo -e "  ${GREEN}✅ firewalld: پورت $PORT باز شد${NC}"
fi

# ─── سرویس systemd ───
if [ "$AS_SERVICE" = true ]; then
  echo -e "${YELLOW}→ ساخت سرویس systemd...${NC}"
  cat > /etc/systemd/system/khan.service <<EOF
[Unit]
Description=Khan Chat Server
After=network.target

[Service]
ExecStart=$INSTALL_DIR/khan-linux -host 0.0.0.0 -port $PORT -data $DATA_DIR
WorkingDirectory=$INSTALL_DIR
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable khan
  systemctl start khan
  echo -e "  ${GREEN}✅ سرویس فعال شد${NC}"
fi

# ─── خلاصه ───
echo ""
echo -e "${GREEN}══════════════════════════════════${NC}"
echo -e "🎉 خان نصب شد!"
echo -e "   باینری: $INSTALL_DIR/khan-linux"
echo -e "   دستور:  $BIN"
echo -e "   داده:   $DATA_DIR"
echo -e "   پورت:   $PORT"
if [ "$AS_SERVICE" = true ]; then
  echo -e "   سرویس:  systemctl status khan"
fi
echo ""
echo -e "🚀 اجرا:"
echo -e "   ${YELLOW}systemctl start khan${NC}   (سرویس)"
echo -e "   ${YELLOW}khan${NC}                   (مستقیم)"
echo ""
echo -e "🌐 مرورگر: ${YELLOW}http://localhost:$PORT${NC}"
echo -e "══════════════════════════════════${NC}"
