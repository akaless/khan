#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════
# khan-security-check.sh — بررسی امنیتی ریپوی خان
# چک می‌کند: توکن، کلید خصوصی، فایل‌های حساس
# ═══════════════════════════════════════════════════════════
set -e

TOKEN="${KHAN_GIT_TOKEN:-}"
REPO_DIR="."
ENV_FILE="/data/.hermes/.env"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0

check() {
  local desc="$1" result="$2"
  if [ "$result" = "ok" ]; then
    echo -e "  ${GREEN}✅${NC} $desc"
    PASS=$((PASS+1))
  else
    echo -e "  ${RED}❌${NC} $desc — $result"
    FAIL=$((FAIL+1))
  fi
}

echo "🔐 بررسی امنیتی ریپوی خان"
echo "══════════════════════════"

# 1. Token در فایل‌های ریپو؟
echo ""
echo "1️⃣ توکن در فایل‌های commit شده؟"
if grep -rq "$TOKEN" "$REPO_DIR" --exclude-dir=.git 2>/dev/null; then
  check "توکن در فایل‌های ریپو نیست" "🔴 توکن لو رفته!"
else
  check "توکن در فایل‌های ریپو نیست" "ok"
fi

# 2. Token در .git/config؟
echo ""
echo "2️⃣ توکن در remote URL؟"
if grep -q "$TOKEN" "$REPO_DIR/.git/config" 2>/dev/null; then
  check "remote URL بدون توکن" "🔴 توکن در .git/config!"
else
  check "remote URL بدون توکن" "ok"
fi

# 3. Token در اسکریپت‌ها؟
echo ""
echo "3️⃣ توکن در اسکریپت‌ها؟"
if grep -rq "$TOKEN" /data/.hermes/scripts 2>/dev/null; then
  check "توکن در اسکریپت نیست" "🔴 توکن در اسکریپت!"
else
  check "توکن در اسکریپت نیست" "ok"
fi

# 4. Token در .env؟
echo ""
echo "4️⃣ توکن در .env (اینجا باید باشد)؟"
if grep -q "KHAN_GIT_TOKEN" "$ENV_FILE" 2>/dev/null; then
  check "توکن در .env ذخیره شده" "ok"
else
  check "توکن در .env ذخیره شده" "🔴 توکن در .env نیست!"
fi

# 5. کلید خصوصی در ریپو؟
echo ""
echo "5️⃣ کلید خصوصی لایسنس؟"
if find "$REPO_DIR" -name "private.key" 2>/dev/null | grep -q .; then
  check "کلید خصوصی در ریپو نیست" "🔴 private.key در ریپو!"
else
  check "کلید خصوصی در ریپو نیست" "ok"
fi

# 6. کلید خصوصی در backup repo؟
echo ""
echo "6️⃣ کلید خصوصی در hermesbackup؟"
if find /data/workspace/hermesbackup -name "private.key" 2>/dev/null | grep -q .; then
  check "کلید خصوصی در backup نیست" "🔴 private.key در backup!"
else
  check "کلید خصوصی در backup نیست" "ok"
fi

# 7. فایل‌های حساس gitignore شده؟
echo ""
echo "7️⃣ .gitignore محافظت می‌کند؟"
if grep -q "private.key" "$REPO_DIR/.gitignore" 2>/dev/null; then
  check ".gitignore از کلید محافظت می‌کند" "ok"
else
  check ".gitignore از کلید محافظت می‌کند" "🔴 .gitignore ناقص!"
fi

echo ""
echo "══════════════════════════"
if [ "$FAIL" = "0" ]; then
  echo -e "${GREEN}🎉 امنیت کامل است — $PASS/$PASS چک پاس شد${NC}"
  exit 0
else
  echo -e "${RED}⚠️ $FAIL مشکل امنیتی پیدا شد!${NC}"
  exit 1
fi
