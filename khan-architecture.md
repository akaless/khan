# 🏠 خان — معماری فنی

> **شعار:** «خانِ شما، جای گپ‌های تیم»

# 🏗️ معماری فنی — نسخه ۱ (خان — چت سازمانی سبک فارسی) — Go Edition

**تاریخ:** ۱ آگوست ۲۰۲۶
**زبان:** Go (Golang)
**پورت پیش‌فرض:** **1727**
**خروجی:** یک فایل exe برای ویندوز (باینری تک‌فایله)

---

## 🎯 انتخاب تکنولوژی (و چرا)

| بخش | انتخاب | چرا |
|-----|--------|-----|
| **زبان** | **Go 1.22+** | باینری تک‌فایله، مصرف رم پایین (~30-60MB)، cross-compile به ویندوز |
| **وب‌سرور** | net/http (استاندارد) | داخلی Go — بدون فریم‌ورک سنگین |
| **روتر** | chi (سبک) | ساده، سازگار با استاندارد net/http |
| **ریل‌تایم** | gorilla/websocket | استاندارد WebSocket در Go |
| **دیتابیس** | SQLite + modernc.org/sqlite | **خالص Go** — بدون CGO، بدون نیاز به کامپایلر C |
| **ORM** | بدون ORM — sqlx (سبک) | کنترل کامل، کوئری‌های مستقیم |
| **رمز هش** | golang.org/x/crypto/argon2 | Argon2id استاندارد |
| **رمزنگاری پیام** | crypto/aes + crypto/cipher (AES-256-GCM) | داخلی Go — بدون کتابخانه خارجی |
| **توکن** | golang-jwt/jwt | JWT استاندارد |
| **فرانت‌اند** | Vue 3 + Vite (سبک) | باندل کوچیک، RTL فارسی عالی |
| **PWA** | manifest + service worker | قابل نصب بدون اپ‌استور |
| **فونت** | وزیرمتن (Vazirmatn) | فونت آزاد فارسی |
| **سرویس ویندوز** | kardianos/service | exe به‌عنوان سرویس ویندوز |
| **پورت** | **1727** | پورت پیش‌فرض |

### چرا Go؟
```
✅ باینری تک‌فایله: go build → khan.exe (فقط یک فایل!)
✅ مصرف رم: ~30-60MB (سبک‌ترین!)
✅ Cross-compile: GOOS=windows go build → exe از لینوکس هم!
✅ بدون CGO: modernc.org/sqlite = بدون نیاز به gcc در ویندوز
✅ gorilla/websocket: WebSocket ساده و قدرتمند
✅ goroutine: مدیریت همزمانی عالی برای هزاران اتصال
✅ استاندارد طلایی برای سرورهای چت (مثل خود SignalR ولی در Go)
```

---

## 📂 ساختار پروژه

```
khan/
├── go.mod                         ← ماژول Go
├── go.sum
├── main.go                        ← نقطه شروع
├── config/
│   └── config.go                  ← بارگذاری تنظیمات (JSON/ENV)
├── internal/
│   ├── models/
│   │   ├── user.go                ← مدل کاربر
│   │   ├── room.go                ← مدل اتاق
│   │   ├── message.go             ← مدل پیام
│   │   ├── reaction.go            ← مدل واکنش
│   │   ├── file.go                ← مدل فایل
│   │   ├── session.go             ← مدل جلسه
│   │   └── enums.go               ← نقش‌ها و نوع اتاق‌ها
│   ├── database/
│   │   ├── db.go                  ← اتصال SQLite + مهاجرت‌ها
│   │   └── schema.sql             ← اسکیمای کامل (embeded)
│   ├── repository/
│   │   ├── user_repo.go           ← عملیات کاربران
│   │   ├── room_repo.go           ← عملیات اتاق‌ها
│   │   ├── message_repo.go        ← عملیات پیام‌ها
│   │   ├── reaction_repo.go       ← عملیات واکنش‌ها
│   │   ├── file_repo.go           ← عملیات فایل‌ها
│   │   └── session_repo.go        ← عملیات جلسات
│   ├── service/
│   │   ├── auth_service.go        ← ورود/توکن/جلسه
│   │   ├── user_service.go        ← مدیریت کاربران + نقش‌ها
│   │   ├── room_service.go        ← مدیریت اتاق‌ها + اعضا
│   │   ├── message_service.go     ← ارسال/دریافت/ویرایش پیام
│   │   ├── crypto_service.go      ← AES-256-GCM
│   │   ├── password_service.go    ← Argon2id
│   │   ├── file_service.go        ← آپلود/دریافت فایل
│   │   └── backup_service.go      ← بکاپ خودکار
│   ├── handler/
│   │   ├── auth_handler.go        ← /api/auth/*
│   │   ├── user_handler.go        ← /api/users/*
│   │   ├── room_handler.go        ← /api/rooms/*
│   │   ├── message_handler.go     ← /api/messages/*
│   │   ├── file_handler.go        ← /api/files/*
│   │   └── ws_handler.go          ← /ws (WebSocket)
│   ├── ws/
│   │   ├── hub.go                 ← مدیریت اتصالات (قلب برنامه)
│   │   ├── client.go              ← کلاینت WebSocket
│   │   ├── events.go              ← تعریف رویدادها
│   │   └── presence.go            ← حضور آنلاین
│   └── middleware/
│       ├── auth_middleware.go     ← بررسی توکن
│       ├── role_middleware.go     ← بررسی نقش
│       └── ratelimit.go           ← محدودیت درخواست
├── web/                           ← فرانت‌اند Vue (بیلدشده → embed)
│   ├── index.html
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── router.ts
│   │   ├── api/ (auth, users, rooms, messages)
│   │   ├── ws/chatClient.ts       ← اتصال WebSocket
│   │   ├── stores/ (auth, chat, users)
│   │   ├── components/ (ChatWindow, MessageBubble, RoomList, ...)
│   │   ├── views/ (Login, Chat, Admin, Profile)
│   │   └── styles/main.css        ← RTL + وزیرمتن
│   └── public/ (manifest, sw.js)
├── webui/
│   └── embed.go                   ← go:embed فرانت‌اند داخل باینری
├── cmd/
│   ├── server/
│   │   └── main.go                ← اجرای سرور
│   └── install/
│       └── main.go                ← نصب سرویس ویندوز
├── scripts/
│   ├── build-windows.sh           ← بیلد exe ویندوز
│   └── setup.iss                  ← Inno Setup نصب‌گر
└── data/                          ← دیتابیس + فایل‌ها (ساخته می‌شه)
    ├── khan.db                 ← SQLite
    ├── uploads/                   ← فایل‌های آپلودی
    └── backups/                   ← بکاپ‌ها
```

---

## 🔌 طراحی API (REST + WebSocket)

### REST API — پورت 1727

**احراز هویت:**
| متد | مسیر | توضیح | دسترسی |
|-----|------|-------|--------|
| POST | /api/auth/login | ورود → توکن | عمومی |
| POST | /api/auth/logout | خروج | همه |
| POST | /api/auth/change-password | تغییر رمز خودش | همه |
| GET | /api/auth/me | اطلاعات خودم | همه |

**کاربران:**
| متد | مسیر | توضیح | دسترسی |
|-----|------|-------|--------|
| GET | /api/users | لیست کاربران | ادمین |
| POST | /api/users | **ایجاد کاربر جدید** | ادمین+ |
| DELETE | /api/users/{id} | **حذف کلی کاربر** | ادمین+ |
| POST | /api/users/{id}/reset-password | ریست رمز | ادمین+ |
| POST | /api/users/{id}/toggle-active | فعال/غیرفعال | ادمین+ |
| PUT | /api/users/{id}/role | ارتقا/تنزل سوپروایزر | ادمین+ |
| PUT | /api/users/{id}/admin | **مدیریت نقش‌ها** | ادمین |
| PUT | /api/users/{id}/profile | ویرایش پروفایل خودش | همه |

**اتاق‌ها:**
| متد | مسیر | توضیح | دسترسی |
|-----|------|-------|--------|
| GET | /api/rooms | لیست اتاق‌های من | همه |
| POST | /api/rooms | ساخت اتاق | همه |
| POST | /api/rooms/{id}/join | عضویت | همه |
| POST | /api/rooms/{id}/members | **اضافه عضو** | سوپروایزر+ (گروهش) |
| DELETE | /api/rooms/{id}/members/{uid} | **حذف عضو** | سوپروایزر+ (گروهش) |
| PATCH | /api/rooms/{id} | تغییر نام/عکس گروه | سوپروایزر+ |

**پیام‌ها:**
| متد | مسیر | توضیح | دسترسی |
|-----|------|-------|--------|
| GET | /api/rooms/{id}/messages?before=X | پیام‌ها (صفحه‌بندی) | عضو |
| PUT | /api/messages/{id} | ویرایش | فرستنده |
| DELETE | /api/messages/{id} | حذف | فرستنده/سوپروایزر گروه |
| POST | /api/messages/{id}/reactions | ریکشن | عضو |

**فایل‌ها:**
| متد | مسیر | توضیح | دسترسی |
|-----|------|-------|--------|
| POST | /api/files | آپلود → fileId | عضو |
| GET | /api/files/{id} | دانلود | عضو |

### WebSocket — /ws (رویدادهای زنده)

**کلاینت → سرور:**
```json
{"type":"send_message","room_id":1,"text":"سلام","file_id":null}
{"type":"edit_message","message_id":5,"text":"ویرایش"}
{"type":"delete_message","message_id":5}
{"type":"add_reaction","message_id":5,"emoji":"👍"}
{"type":"typing","room_id":1}
{"type":"mark_read","room_id":1,"last_message_id":42}
{"type":"join_room","room_id":1}
```

**سرور → کلاینت:**
```json
{"type":"message","room_id":1,"message":{...}}
{"type":"message_edited","message_id":5,"text":"..."}
{"type":"message_deleted","message_id":5}
{"type":"reaction","message_id":5,"emoji":"👍","user_id":3}
{"type":"typing","room_id":1,"user_id":3,"username":"ali"}
{"type":"presence","user_id":3,"online":true}
{"type":"user_added","room_id":1,"user":{...}}
{"type":"user_removed","room_id":1,"user_id":7}
{"type":"user_role_changed","user_id":4,"new_role":"supervisor"}
{"type":"user_deleted","user_id":9}
{"type":"force_logout","user_id":4}
```

---

## 🔐 طراحی امنیت

### جریان ورود:
```
۱. POST /api/auth/login {username, password}
۲. رمز با Argon2id بررسی
۳. موفق → JWT (7 روز) برمی‌گرده
۴. کلاینت توکن رو ذخیره می‌کنه
۵. هر درخواست: Authorization: Bearer <token>
```

### جریان رمزنگاری پیام (AES-256-GCM):
```
ارسال:   کلاینت → [TLS] → سرور → AES-256-GCM رمزنگاری → SQLite
دریافت:  SQLite → AES-256-GCM رمزگشایی → [TLS] → کلاینت

- کلید AES-256: در فایل config (با دسترسی محدود)
- هر پیام: IV تصادفی ۱۲ بایت + tag ۱۶ بایت
- GCM = رمزنگاری + احراز اصالت (ضد دستکاری)
```

### هش رمز (Argon2id):
```go
salt := make([]byte, 16)
rand.Read(salt)
hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32)
// ذخیره: hex(salt) + "$" + hex(hash)
// بررسی: دوباره محاسبه + مقایسه
```

### محافظت ورود:
```
- ۵ تلاش اشتباه → قفل ۵ دقیقه (IP + یوزرنیم)
- JWT: ۷ روز انقضا
- تغییر رمز: نیاز به رمز قبلی
- ریست رمز: رمز موقت + must_change_pwd=1
```

---

## 🗄️ اسکیمای SQLite (کامل)

```sql
-- کاربران
CREATE TABLE users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,          -- Argon2id: salt$hash
    display_name    TEXT DEFAULT '',
    avatar_path     TEXT DEFAULT NULL,
    role            TEXT NOT NULL DEFAULT 'user',
                    -- 'user' | 'supervisor' | 'admin' | 'super_admin'
    hidden          INTEGER NOT NULL DEFAULT 0,
    active          INTEGER NOT NULL DEFAULT 1,
    must_change_pwd INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL,
    last_seen       TEXT
);

-- اتاق‌ها
CREATE TABLE rooms (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    type        TEXT NOT NULL,              -- 'dm' | 'group' | 'public' | 'private'
    creator_id  INTEGER NOT NULL,
    created_at  TEXT NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- اعضای اتاق
CREATE TABLE room_members (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id   INTEGER NOT NULL,
    user_id   INTEGER NOT NULL,
    role      TEXT NOT NULL DEFAULT 'member',  -- 'owner' | 'member'
    joined_at TEXT NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (room_id, user_id)
);

-- پیام‌ها
CREATE TABLE messages (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id    INTEGER NOT NULL,
    sender_id  INTEGER NOT NULL,
    body       BLOB NOT NULL,               -- AES-256-GCM
    iv         BLOB NOT NULL,
    file_id    INTEGER,
    created_at TEXT NOT NULL,
    edited_at  TEXT,
    deleted_at TEXT,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);

-- واکنش‌ها
CREATE TABLE reactions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    emoji      TEXT NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (message_id, user_id, emoji)
);

-- جلسات
CREATE TABLE sessions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    token      TEXT UNIQUE NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- فایل‌ها
CREATE TABLE files (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id   INTEGER NOT NULL,
    room_id    INTEGER NOT NULL,
    file_name  TEXT NOT NULL,
    stored_as  TEXT NOT NULL,
    size       INTEGER NOT NULL,
    mime_type  TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

-- ایندکس‌ها
CREATE INDEX idx_messages_room ON messages(room_id, created_at);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_user ON room_members(user_id);
CREATE INDEX idx_sessions_user ON sessions(user_id);
```

---

## ⚙️ تنظیمات (config.json)

```json
{
  "server": {
    "port": 1727,
    "host": "0.0.0.0",
    "address_type": "ip",           // "ip" | "dns"
    "ip": "192.168.1.100",
    "dns": "",                      // مثلاً chat.mycompany.com
    "db_path": "data/khan.db",
    "data_dir": "data",
    "max_upload_mb": 50
  },
  "license": {
    "license_file": "data/license.key"
  },
  "security": {
    "encryption_key": "CHANGE_ME_RANDOM_32_BYTES",
    "jwt_secret": "CHANGE_ME_RANDOM_SECRET",
    "jwt_expire_days": 7,
    "max_login_attempts": 5,
    "lockout_minutes": 5
  },
  "backup": {
    "enabled": true,
    "interval_hours": 24,
    "keep_count": 7,
    "backup_path": "data/backups"
  },
  "admin": {
    "username": "SET_AT_INSTALL",
    "password": "SET_AT_INSTALL"
  }
}
```

### 🔧 تنظیم آدرس سرور (IP یا DNS — ۲ روش):

**روش ۱ — موقع نصب (نصب‌گر Inno Setup):**
- صفحه «تنظیمات شبکه» در نصب‌گر:
  - گزینه «نوع آدرس»: 🔘 IP 🔘 DNS
  - اگه IP: فیلد آدرس IP سرور (مثلاً 192.168.1.100)
  - اگه DNS: فیلد آدرس DNS (مثلاً chat.mycompany.com)
  - پورت (پیش‌فرض 1727 — محدوده 1700 تا 1799)
- نوشته می‌شه توی config.json
- بعد از نصب: میانبر دسکتاپ → آدرس تنظیم‌شده

**روش ۲ — از داخل برنامه (پنل ادمین):**
- بخش «تنظیمات سرور» در پنل ادمین:
  - گزینه «نوع آدرس»: 🔘 IP 🔘 DNS
  - فیلد IP یا DNS
  - فیلد پورت (فقط دو رقم آخر: ۱۷[۰۰-۹۹] → 1700 تا 1799)
  - دکمه ذخیره
- تغییر → ذخیره در config.json + نمایش «بعد از ری‌استارت اعمال می‌شه»
- دکمه «ری‌استارت سرور» (توسط ادمین)

### 🔢 محدوده پورت (دو رقم بعد از ۱۷):

**قانون:** پورت پیش‌فرض `1727` — اما ادمین می‌تونه **دو رقم آخر** رو تغییر بده:

| فیلد | محدوده | مثال |
|------|--------|------|
| پورت | 1700 – 1799 | 1727، 1745، 1790، 1710 |

**پیاده‌سازی UI:**
```
پورت:  17[ ▢▢ ]        ← فقط ۲ رقم قابل تایپ (00-99)
        ↑↑
        ثابت  قابل تغییر
```

**اعتبارسنجی:**
- ورودی فقط ۲ رقم (00-99)
- ذخیره: `1700 + value` → پورت نهایی
- مثال: کاربر می‌زنه `45` → پورت `1745`
- اگه ۱۷۲۷ گرفته بود، ادمین می‌تونه `45` بزنه و مشکل حل بشه

**توی URL و مستندات:**
```
http://192.168.1.100:1745
http://chat.mycompany.com:1745
```

### 🎯 قفل ۲۰ کاربر — نسخه حرفه‌ای (ضد دستکاری!)

**⚠️ مهم:** محدودیت کاربر **عدد ثابت توی کد نیست!** — یک **لایسنس رمزنگاری‌شده با امضای دیجیتال** هست که:
- با کلید خصوصی (پیش ما) امضا می‌شه
- کلید عمومی داخل باینری تعبیه می‌شه
- اگه کسی لایسنس رو دستکاری کنه → امضا نامعتبر → برنامه کار نمی‌کنه

#### 🔐 فایل لایسنس (license.key):

```json
{
  "version": 1,
  "licensed_to": "My Company",
  "max_users": 20,
  "issued_at": "2026-08-01T00:00:00Z",
  "expires_at": "2027-08-01T00:00:00Z",
  "features": ["chat", "files", "reactions"],
  "signature": "Ed25519_SIGNATURE_HERE"
}
```

**چرخه تولید لایسنس (پیش ما):**
```
۱. کلید خصوصی Ed25519 رو نگه می‌داریم (هیچ‌جا انتشار نمی‌دیم)
۲. برای هر شرکت: {licensed_to, max_users, expires_at} → امضا با کلید خصوصی
۳. خروجی: license.key → به شرکت داده می‌شه
۴. شرکت فایل رو می‌ذاره توی data/license.key
```

**بررسی در زمان اجرا (سرور):**
```
۱. برنامه استارت می‌شه
۲. license.key رو می‌خونه
۳. امضا رو با کلید عمومی (داخل باینری) بررسی می‌کنه
   ├── امضا درست → max_users از لایسنس خونده می‌شه (۵۰/۱۰۰/۵۰۰/نامحدود)
   ├── امضا غلط / دستکاری‌شده → **محدود به ۵ کاربر** (جریمه!)
   └── فایل نیست → **۲۰ کاربر رایگان** (پیش‌فرض برای شرکت‌های کوچیک)
۴. max_users هر بار از لایسنس خونده می‌شه — نه از کد!
```

### ⚠️ هشدار موقع وارد کردن لایسنس (UI ادمین):

**پنل ادمین → بخش لایسنس:**

```
🛡️ مدیریت لایسنس

┌─────────────────────────────────────────────┐
│ ⚠️ توجه مهم:                                │
│                                            │
│ وارد کردن لایسنس اشتباه یا دستکاری‌شده      │
│ محدودیت کاربران را به ۵ نفر کاهش می‌دهد!    │
│                                            │
│ • لایسنس معتبر → حداکثر کاربران افزایش می‌یابد │
│ • بدون لایسنس → ۲۰ کاربر رایگان             │
│ • لایسنس اشتباه → فقط ۵ کاربر (جریمه)       │
│                                            │
│ [📁 انتخاب فایل لایسنس]  [✅ اعمال]         │
└─────────────────────────────────────────────┘
```

**جریان:**
```
۱. ادمین روی «مدیریت لایسنس» کلیک می‌کنه
۲. هشدار قرمز بالا نمایش داده می‌شه (قبل از انتخاب فایل!)
۳. ادمین فایل license.key رو انتخاب می‌کنه
۴. دکمه «اعمال» → سرور امضا رو بررسی می‌کنه
۵. نتیجه:
   ✅ معتبر → «لایسنس اعمال شد: ۵۰ کاربر»
   ❌ نامعتبر → «⚠️ لایسنس نامعتبر! محدودیت به ۵ کاربر کاهش یافت»
6. بعد از خطا: دکمه «حذف لایسنس» → برگشت به ۲۰ کاربر رایگان
```

**جزئیات پیام خطا (نمایش به ادمین):**
```
❌ لایسنس نامعتبر است!
   امضای دیجیتال تأیید نشد — فایل ممکن است دستکاری شده باشد.
   
   ⚠️ محدودیت کاربران شما به ۵ نفر کاهش یافت.
   
   [🗑️ حذف لایسنس و بازگشت به ۲۰ کاربر]
```

**نکته:** اگه لایسنس اشتباه وارد شد، ادمین می‌تونه حذفش کنه → برگرده به ۲۰ کاربر. فقط تا وقتی فایل غلط هست، ۵ کاربره.

### 🆓 سیاست کاربران (خلاصه):

| وضعیت | حداکثر کاربر | برای کی |
|-------|:-----------:|---------|
| **بدون لایسنس** | **۲۰ کاربر** | همه — شرکت‌های کوچیک راحت استفاده کنن |
| لایسنس معتبر | طبق لایسنس (۵۰/۱۰۰/۵۰۰/∞) | شرکت‌هایی که بزرگ‌تر شدن |
| **لایسنس دستکاری‌شده** | **۵ کاربر** (جریمه) | کسی که خواسته تقلب کنه |

**منطق:**
```
- هدف: شرکت‌های کوچیک بدون هیچ دردسری از ۲۰ کاربر استفاده کنن
- لایسنس فقط برای بزرگ‌تر شدن لازمه (نه برای شروع!)
- دستکاری لایسنس = نه فقط بی‌اثر، بلکه جریمه (از ۲۰ می‌افته به ۵)
- بنابراین تقلب هیچ سودی نداره — ضرر داره!
```

#### 🛡️ چرا دستکاریش نمی‌شه؟

| حمله | چرا شکست می‌خوره |
|------|-----------------|
| عوض کردن عدد 20 توی license.key | امضا نامعتبر می‌شه — برنامه رد می‌کنه |
| جستجوی `max_users` توی کد | اسم فیلد توی کد **مبهم‌سازی‌شده** — پیدا نمی‌شه |
| جستجوی عدد 20 توی کد | عدد 20 **هیچ‌جا توی کد نیست** — فقط داخل لایسنس امضاشده |
| کامپایل مجدد با عدد جدید | کلید عمومی امضا رو چک می‌کنه — لایسنس جعلی رد می‌شه |
| تغییر باینری | امضای باینری خودش هم بررسی می‌شه |
| فیک کردن کلید | کلید خصوصی پیش ماست — غیرممکنه |

#### 🧩 مبهم‌سازی در کد (Obfuscation):

```go
// ❌ نه اینطوری (قابل جستجو):
// maxUsers := 20

// ✅ اینطوری — عدد 20 رمزنگاری‌شده و توی کد نیست:
// مقدار اصلی با XOR مخفی شده — حتی grep هم پیداش نمی‌کنه
const _licKey int32 = 0x46  // 20 XOR 0x5A
// در runtime:
func licenseMaxUsers() int {
    if v, ok := verifyLicense(_licensePath); ok {
        return v.MaxUsers          // از فایل امضاشده
    }
    return int(_licKey ^ 0x5A)     // 20 — حالت آزمایشی
}
```

**نام‌گذاری مبهم:**
- متغیرهای لایسنس: اسم‌های بی‌ربط (`_x7kQ`, `seat`, `cfgHash`)
- فیلدهای JSON لایسنس: `max_users` → `m_u` (مختصر)
- توابع بررسی: `verifyLicense` → `_chk()` 
- رشته‌های خطا: رمزنگاری‌شده، فقط در runtime دیکد می‌شن

#### 📊 نمایش در پنل ادمین:

```
📊 ظرفیت: ۱۲ / ۲۰ کاربر       ← ۲۰ از لایسنس خونده می‌شه
███████████░░░░░░░░░
[+ کاربر جدید] ← اگه پر باشه غیرفعال می‌شه

🛡️ لایسنس: My Company
   معتبر تا: ۲۰۲۷/۰۸/۰۱
```

#### 💼 سیستم لایسنس برای نسخه‌های بعد:

| نسخه | max_users | توضیح |
|------|-----------|-------|
| نسخه ۱ (رایگان/شروع) | ۲۰ | لایسنس پایه |
| نسخه ۱.۵ | ۵۰ | لایسنس کوچک |
| نسخه ۲ | ۱۰۰ | لایسنس متوسط |
| نسخه ۲.۵ | ۵۰۰ | لایسنس بزرگ |
| نسخه سازمانی | نامحدود | لایسنس سازمانی |

**نکته:** همه با همون کلید امضا می‌شن — فقط عدد max_users در لایسنس فرق می‌کنه. سیستم از قبل آماده‌ست!


---

## 🖥️ بیلد و استقرار

### بیلد exe ویندوز (از هر سیستمی!):
```bash
# از لینوکس/مک/ویندوز:
GOOS=windows GOARCH=amd64 go build -o khan.exe ./cmd/server

# با embed فرانت‌اند (یک فایل کامل):
cd web && npm install && npm run build
cd .. && GOOS=windows GOARCH=amd64 go build -o khan.exe ./cmd/server

# خروجی: khan.exe (~40-60MB) — همه چیز داخلشه!
```

### اجرا:
```bash
# مستقیم:
khan.exe
# → سرور روی http://0.0.0.0:1727

# به‌عنوان سرویس ویندوز:
khan.exe install
khan.exe start
```

### نصب‌گر (Inno Setup):
```
- KhanSetup.exe
- نصب به C:\Program Files\Khan\
- ساخت سرویس: khan.exe install
- شروع خودکار با ویندوز
- میانبر → http://localhost:1727
```

### استفاده نهایی:
```
۱. نصب KhanSetup.exe
   → موقع نصب: آدرس IP سرور رو وارد می‌کنی (مثلاً 192.168.1.100)
۲. سرویس شروع می‌شه → http://localhost:1727 (خود سرور)
۳. بقیه کامپیوترها: http://192.168.1.100:1727 (IP تنظیم‌شده)
۴. ورود با یوزنیم + رمز (ساخته‌شده توسط ادمین)
۵. ادمین می‌تونه IP رو بعداً از پنل مدیریت عوض کنه
```

### 🔄 تغییر IP بعد از نصب (از پنل ادمین):
```
پنل ادمین → تنظیمات سرور → تغییر IP/پورت → ذخیره → ری‌استارت
```

---

## 📦 وابستگی‌های Go (go.mod)

```
module khan

go 1.22

require (
    github.com/go-chi/chi/v5            // روتر
    github.com/gorilla/websocket        // WebSocket
    modernc.org/sqlite                  // SQLite خالص Go (بدون CGO)
    github.com/jmoiron/sqlx             // کوئری‌های راحت‌تر
    golang.org/x/crypto/argon2          // هش رمز
    github.com/golang-jwt/jwt/v5        // توکن JWT
    github.com/kardianos/service        // سرویس ویندوز
)
```

---

## 📅 نقشه راه توسعه

### هفته ۱: زیرساخت
- [ ] go mod + ساختار پوشه‌ها + config
- [ ] SQLite + schema.sql + مهاجرت
- [ ] AuthService + ورود + JWT
- [ ] UI پایه (ورود + لیست اتاق‌ها)

### هفته ۲: چت
- [ ] WebSocket Hub + ارسال/دریافت
- [ ] DM + گروه + عمومی/خصوصی
- [ ] پیام آفلاین + تایپینگ + خوانده/نخوانده
- [ ] رمزنگاری AES-256-GCM

### هفته ۳: مدیریت
- [ ] نقش‌ها و مجوزها (۴ سطح)
- [ ] CRUD کاربران (ادمین)
- [ ] ریست رمز + اجبار تغییر

### هفته ۴: تکمیل
- [ ] فایل‌شیرینگ + ایموجی + ریکشن + ذکر
- [ ] RTL + وزیرمتن + تقویم شمسی
- [ ] PWA + حالت تاریک
- [ ] بکاپ + بیلد exe + نصب‌گر

---

## 🎯 خلاصه

```
✅ Go — سبک‌ترین و ساده‌ترین برای سرور چت
✅ پورت 1727
✅ باینری تک‌فایله exe (~40-60MB)
✅ مصرف رم ~30-60MB
✅ SQLite خالص Go (بدون CGO)
✅ Argon2id + AES-256-GCM
✅ WebSocket (gorilla)
✅ Vue 3 + PWA فارسی
✅ Cross-compile: از هر سیستم‌عصلی exe بساز
✅ embed: فرانت‌اند داخل باینری — یک فایل همه چیز
```
