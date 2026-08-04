# 📦 نصب‌گرهای خان (Installers)

برای هر سه پلتفرم نصب‌گر آماده موجود است.

## 🪟 ویندوز

### `KhanSetup-1.0.0.exe` (NSIS Installer)

- **نصب:** دوبار کلیک → Next → Finish
- **مسیر:** `C:\Program Files\Khan\`
- **میان‌بر:** دسکتاپ + منوی استارت
- **حذف:** از Control Panel → Programs → خان

**ساخت از سورس:**
```bash
# پیش‌نیاز: NSIS (apt install nsis)
cd installers/windows
cp ../../dist/khan.exe .
makensis installer.nsi
# خروجی: KhanSetup-1.0.0.exe
```

---

## 🐧 لینوکس

### گزینه ۱: بسته `.deb` (دبیان/اوبونتو)

```bash
sudo dpkg -i khan_1.0.0_amd64.deb
```

- نصب: `/usr/bin/khan`
- حذف: `sudo apt remove khan`

### گزینه ۲: اسکریپت نصب (هر توزیع)

```bash
chmod +x install.sh
sudo ./install.sh               # فقط باینری
sudo ./install.sh --service     # + سرویس systemd
sudo ./install.sh --port 1800   # + پورت دلخواه
```

- نصب: `/opt/khan/` + سرویس systemd
- حذف: `sudo ./uninstall.sh`

### گزینه ۳: داکر

```bash
docker run -d -p 1727:1727 -v $(pwd)/data:/app/data khan
```

---

## 🍎 مک

### گزینه ۱: باینری مستقیم

```bash
chmod +x khan-macos-arm64   # یا khan-macos-intel
./khan-macos-arm64
```

### گزینه ۲: اپلیکیشن + DMG

```bash
# روی خود مک (نیاز به Xcode CLT)
chmod +x build-dmg.sh
./build-dmg.sh
# خروجی: Khan-1.0.0-<arch>.dmg
```

- DMG را باز کنید → خان را به Applications بکشید
- برای باز کردن از Gatekeeper: **کلیک راست → Open**

---

## 🧪 تست نصب‌گرها

### ویندوز (VM یا کامپیوتر واقعی)
```powershell
# نصب
.\KhanSetup-1.0.0.exe /S
# اجرا
& "C:\Program Files\Khan\khan.exe"
# حذف
& "C:\Program Files\Khan\Uninstall.exe" /S
```

### لینوکس
```bash
# نصب .deb
sudo dpkg -i khan_1.0.0_amd64.deb
khan --version
systemctl status khan   # اگر سرویس

# حذف
sudo apt remove khan
```

---

## 📁 ساختار

```
installers/
├── windows/
│   └── installer.nsi          ← NSIS script
├── linux/
│   ├── deb/                   ← ساختار .deb
│   ├── install.sh             ← نصب عمومی
│   ├── uninstall.sh           ← حذف
│   └── khan_1.0.0_amd64.deb   ← بسته آماده
└── macos/
    └── build-dmg.sh           ← ساخت DMG روی مک
```
