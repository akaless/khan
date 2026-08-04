<div align="center">

# 🏠 خان (Khan)

### چت سازمانی سبک فارسی — «خانِ شما، جای گپ‌های تیم»
### Lightweight Persian LAN Chat — "Your Khan, Where Teams Chat"

![Version](https://img.shields.io/badge/version-1.0.1-blue)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)
![License](https://img.shields.io/badge/license-MIT-green)
![Size](https://img.shields.io/badge/size-11MB-orange)

**خان** یک چت سازمانی سبک، مدرن و امن برای شبکه‌های محلی است — یک فایل، بدون نصب، بدون اینترنت، بدون سرور جدا.

**Khan** is a lightweight, modern, and secure team chat for local networks — one file, no install, no internet, no separate server.

<img src="branding/khan-logo.jpg" alt="Khan Logo" width="200">

</div>

---

## ✨ چرا خان؟ / Why Khan?

| نیاز / Need | راه‌حل خان / Khan's Solution |
|-----------|------|
| 🚫 بدون اینترنت / Offline | کاملاً آفلاین روی LAN / Fully offline on LAN |
| 📦 نصب سخت / Complex setup | **یک فایل exe** — اجرا و تمام! / **One exe** — run and done! |
| 💰 سرور گران / Expensive server | بدون نیاز به سرور جدا / No separate server needed |
| 🔒 حریم خصوصی / Privacy | داده‌ها روی سیستم خودتان / Data stays on your machine |
| 🌐 فیلترینگ / Censorship | هیچ وابستگی به CDN/اینترنت / No CDN/internet dependency |

## 🚀 شروع سریع / Quick Start

```bash
# ۱. دانلود باینری / Download binary (from Releases or dist/)
# ۲. اجرا / Run:
./khan-linux        # لینوکس/مک / Linux/macOS
khan.exe            # ویندوز / Windows

# ۳. مرورگر / Browser:
#    http://localhost:1727
```

> 💡 کاربران LAN با `http://IP-سرور:1727` وصل می‌شوند
> 💡 LAN users connect via `http://SERVER-IP:1727`

## 📦 باینری‌ها و نصب‌گرها / Binaries & Installers

| فایل / File | پلتفرم / Platform | نوع / Type |
|------|--------|-----|
| `dist/khan.exe` | ویندوز 10/11 (x64) | پورتابل / Portable |
| `dist/khan-linux` | لینوکس (x64) | پورتابل / Portable |
| `dist/khan-macos-intel` | مک اینتل / macOS Intel | پورتابل / Portable |
| `dist/khan-macos-arm64` | مک سیلیکون / macOS Silicon | پورتابل / Portable |
| `dist/installers/KhanSetup-1.0.0.exe` | ویندوز / Windows | **نصب‌گر / Installer** |
| `dist/installers/khan_1.0.0_amd64.deb` | اوبونتو/دبیان / Ubuntu/Debian | **بسته نصب / Package** |
| `installers/linux/install.sh` | هر لینوکس / Any Linux | **اسکریپت نصب / Script** |
| `installers/macos/build-dmg.sh` | مک / macOS | ساخت اپ / App builder |

## 🎯 ویژگی‌ها / Features

### 💬 چت / Chat
- پیام‌رسانی زنده (WebSocket) / Real-time messaging (WebSocket)
- ویرایش / حذف / واکنش (👍❤️😂👏) / Edit / Delete / Reactions
- ایموجی‌پیکر + تایپینگ زنده / Emoji picker + live typing
- گروه‌بندی پیام‌ها بر اساس روز (شمسی) / Day-grouped messages (Persian calendar)
- تاریخچه + بارگذاری نامحدود / History + infinite scroll

### 👥 نقش‌ها / Roles (3 levels)
| نقش / Role | امکانات / Permissions |
|-----|---------|
| 👤 کاربر / User | چت + تغییر رمز خودش / Chat + change own password |
| 🛡 سوپروایزر / Supervisor | مدیریت گروه‌ها / Manage groups |
| ⚙️ ادمین / Admin | مدیریت کاربران + تنظیمات / Manage users + settings |

### 🎫 لایسنس / License (Ed25519)
| وضعیت / Status | کاربران / Users |
|--------|---------|
| 🆓 بدون لایسنس / No license | ۲۰ / 20 |
| ✅ معتبر / Valid | ۵۰ / ۱۰۰ / ۵۰۰ / ∞ |
| ⚠️ دستکاری / Tampered | ۵ / 5 (penalty) |

### 🔐 امنیت / Security
- Argon2id (هش رمز / password hashing)
- AES-256-GCM (رمزنگاری پیام / message encryption)
- قفل ۵ دقیقه بعد از ۵ تلاش ناموفق / 5-min lockout after 5 failed attempts
- توکن‌های امن + مدیریت جلسه / Secure tokens + session management

### 🎨 UI
- الهام از تلگرام — تیره و مدرن / Telegram-inspired — dark & modern
- فارسی RTL + انگلیسی / Persian RTL + English
- PWA (قابل نصب روی موبایل / installable on mobile)
- واکنش‌گرا — موبایل/دسکتاپ / Responsive — mobile/desktop

## 🛠 توسعه / Development

```bash
git clone https://codeberg.org/adiib/khan.git
cd khan/code

# اجرا / Run
go run ./cmd/server

# تست / Test
go test ./...

# بیلد همه پلتفرم‌ها / Build all platforms
cd .. && make build-all
```

## 📚 مستندات / Documentation

- 📖 [راهنمای نصب / Installation Guide](docs/installation.md)
- 🎫 [سیستم لایسنس / Licensing](docs/licensing.md)
- 📡 [راهنمای API / API Reference](docs/api.md)
- 🏛 [معماری / Architecture](khan-architecture.md)
- 📝 [اسپک نسخه ۱ / v1 Spec](khan-spec.md)
- 📜 [تغییرات / Changelog](CHANGELOG.md)

## 🤝 مشارکت / Contributing

راهنمای مشارکت: [CONTRIBUTING.md](CONTRIBUTING.md) / See [CONTRIBUTING.md](CONTRIBUTING.md)

- گزارش باگ / Report bugs → [Issues](https://codeberg.org/adiib/khan/issues)
- پیشنهاد ویژگی / Feature request → Issue با برچسب / with label `enhancement`
- کد / Code → Fork + Pull Request

## 🔐 امنیت / Security

آسیب‌پذیری پیدا کردید؟ / Found a vulnerability? → [SECURITY.md](SECURITY.md)

**هرگز کلید خصوصی لایسنس یا توکن‌ها را در ریپو نگذارید.**
**Never commit private license keys or tokens.**

## 🐳 داکر / Docker

```bash
docker build -t khan .
docker run -d -p 1727:1727 -v $(pwd)/data:/app/data khan
```

## 📄 لایسنس / License

[MIT](LICENSE) © 2026 Adib Mombini (aDiB)

---

<div align="center">

**ساخته شده با ❤️ و Go — خان، چت سازمانی سبک فارسی 🏠🧔**
**ساخته شده توسط aDiB — Built with ❤️ and Go, by aDiB**

</div>
