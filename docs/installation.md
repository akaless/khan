# 📖 راهنمای نصب خان

این راهنما نصب و راه‌اندازی خان را در سیستم‌عامل‌های مختلف توضیح می‌دهد.

## 📋 پیش‌نیازها

- **ویندوز:** 10 یا 11 (64-bit)
- **لینوکس:** هر توزیع با هسته 4.18+ (x86_64)
- **مک:** Monterey+ (Intel یا Apple Silicon)
- **رم:** حداقل ۶۴MB
- **فضای دیسک:** ۵۰MB
- بدون نیاز به نصب هیچ وابستگی — باینری استاتیک

## ⚡ نصب سریع

### ویندوز

1. از [Releases](https://codeberg.org/adiib/khan/releases) فایل `khan.exe` را دانلود کنید
2. آن را در پوشه دلخواه (مثلاً `C:\Khan\`) قرار دهید
3. دوبار کلیک کنید یا از CMD:

```cmd
C:\Khan\khan.exe
```

### لینوکس

```bash
# دانلود
curl -LO https://codeberg.org/adiib/khan/releases/download/v1.0.0/khan-linux
chmod +x khan-linux

# اجرا
./khan-linux
```

### مک

```bash
# اپل سیلیکون (M1-M4)
curl -LO https://codeberg.org/adiib/khan/releases/download/v1.0.0/khan-macos-arm64
chmod +x khan-macos-arm64 && ./khan-macos-arm64

# اینتل
curl -LO https://codeberg.org/adiib/khan/releases/download/v1.0.0/khan-macos-intel
chmod +x khan-macos-intel && ./khan-macos-intel
```

> 💡 اگر Gatekeeper اجازه نداد: **کلیک راست → Open**

## 🚀 راه‌اندازی اول

1. مرورگر را باز کنید: `http://localhost:1727`
2. **ویزارد نصب اولیه** ظاهر می‌شود (فقط بار اول)
3. نام کاربری مدیر، نام شرکت و رمز عبور بسازید
4. وارد شوید و شروع کنید!

## 🌐 دسترسی از شبکه محلی (LAN)

### ویندوز — فایروال

```powershell
# PowerShell (مدیریت)
New-NetFirewallRule -DisplayName "Khan Chat" -Direction Inbound -Port 1727 -Protocol TCP -Action Allow
```

### لینوکس — فایروال

```bash
# firewalld
sudo firewall-cmd --permanent --add-port=1727/tcp
sudo firewall-cmd --reload

# یا UFW (اوبونتو/دبیان)
sudo ufw allow 1727/tcp
```

### مک — فایروال

سیستم‌ها → حریم خصوصی و امنیت → فایروال → اجازه ورودی

### پیدا کردن IP سرور

```bash
# ویندوز
ipconfig | findstr IPv4

# لینوکس
ip addr show | grep inet

# مک
ifconfig | grep inet
```

سپس بقیه کاربران با `http://IP-سرور:1727` وصل می‌شوند (مثلاً `http://192.168.1.20:1727`).

## ⚙️ اجرا به‌صورت سرویس

### لینوکس (systemd)

```bash
sudo tee /etc/systemd/system/khan.service <<'EOF'
[Unit]
Description=Khan Chat Server
After=network.target

[Service]
ExecStart=/opt/khan/khan-linux
WorkingDirectory=/opt/khan
Restart=always
RestartSec=3
User=khan

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now khan
sudo systemctl status khan
```

### ویندوز (Task Scheduler)

1. Task Scheduler را باز کنید
2. Create Task → Triggers → At startup
3. Action → Start program → `C:\Khan\khan.exe`
4. Settings → Restart on failure ✅

## 🛠 گزینه‌های خط فرمان

```bash
khan-linux [flags]

Flags:
  -config string   مسیر فایل تنظیمات (پیش‌فرض: config.json)
  -data string     مسیر پوشه داده (پیش‌فرض: data/)
  -port int        پورت (پیش‌فرض: 1727)
  -host string     آدرس گوش‌دادن (پیش‌فرض: 0.0.0.0)
  -version         نمایش نسخه
```

## 🔄 به‌روزرسانی

1. باینری جدید را دانلود کنید
2. باینری قدیمی را جایگزین کنید
3. داده‌ها در پوشه `data/` حفظ می‌شوند
4. سرویس را ری‌استارت کنید

## 🗄 بکاپ

داده‌ها در فایل‌های JSON در پوشه `data/` ذخیره می‌شوند:

```bash
# بکاپ ساده — فقط کپی فایل
cp -r /opt/khan/data /backup/khan-data-$(date +%Y%m%d)

# بازیابی
cp -r /backup/khan-data-20260801/* /opt/khan/data/
```

## ❓ مشکلات رایج

| مشکل | راه‌حل |
|------|--------|
| پورت ۱۷۲۷ اشغال است | از `-port 1728` یا از پنل مدیریت تغییر دهید |
| دیگران وصل نمی‌شوند | فایروال را بررسی کنید (پورت TCP 1727) |
| صفحه خالی بعد از ورود | نسخه جدید را دانلود کنید (باگ CDN رفع شده) |
| فراموشی رمز ادمین | پوشه data را بکاپ کنید، دوباره نصب و بازیابی |

---
🏠 **خان** — سوالی دارید؟ [Issue باز کنید](https://codeberg.org/adiib/khan/issues)
