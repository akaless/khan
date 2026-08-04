# 📡 راهنمای API خان

مستندات کامل REST API و WebSocket خان — نسخه 1.0

**Base URL:** `http://<server>:1727/api`

## 🔐 احراز هویت

اکثر اندپوینت‌ها به هدر نیاز دارند:

```
Authorization: Bearer <token>
```

توکن از `POST /auth/login` دریافت می‌شود.

---

## 🏠 نصب اولیه

### `GET /setup/needs-setup`
بررسی نیاز به نصب اولیه

**پاسخ:**
```json
{ "needs_setup": true }
```

### `POST /setup`
ساخت مدیر اصلی (فقط بار اول)

**بدنه:**
```json
{
  "admin_username": "modir",
  "admin_password": "Secure@123",
  "company_name": "شرکت آفتاب"
}
```

**پاسخ:** `{ "ok": true }`

---

## 🔑 احراز هویت

### `POST /auth/login`
ورود

**بدنه:**
```json
{ "username": "ali", "password": "Pass@123" }
```

**پاسخ:**
```json
{
  "token": "eyJ...",
  "user": {
    "id": 2,
    "username": "ali",
    "display_name": "علی رضایی",
    "role": "user",
    "active": true,
    "must_change_pwd": true
  }
}
```

### `POST /auth/logout`
خروج — توکن باطل می‌شود

### `GET /auth/me`
اطلاعات کاربر جاری

### `POST /auth/change-password`
تغییر رمز

**بدنه:**
```json
{
  "current_password": "Old@123",
  "new_password": "New@123"
}
```

### `GET /auth/sessions`
لیست جلسات فعال

### `DELETE /auth/sessions/{id}`
بستن یک جلسه

---

## 👥 کاربران (ادمین)

### `GET /users`
لیست همه کاربران

### `POST /users`
ساخت کاربر جدید

**بدنه:**
```json
{
  "username": "reza",
  "display_name": "رضا کریمی",
  "password": "Pass@123",
  "role": "user"  // user | sup
}
```


### `POST /users/{id}/role`
تغییر نقش (user ↔ sup)

**بدنه:** `{ "role": "sup" }`

### `POST /users/{id}/reset-password`
ریست رمز

**بدنه:** `{ "new_password": "New@123" }`

### `POST /users/{id}/toggle-active`
فعال/غیرفعال

### `DELETE /users/{id}`
حذف کاربر (ظرفیت آزاد می‌شود)

---

## 🏘 اتاق‌ها

### `GET /rooms`
لیست اتاق‌های من

### `POST /rooms`
ساخت گروه

**بدنه:**
```json
{
  "name": "تیم فروش",
  "type": "public"  // public | private
}
```

### `GET /rooms/{id}`
جزئیات اتاق

### `POST /rooms/{id}/join`
عضویت در گروه عمومی

### `POST /rooms/{id}/members`
افزودن عضو (سوپروایزر+)

**بدنه:** `{ "user_id": 5 }`

### `GET /rooms/{id}/members`
لیست اعضا

### `DELETE /rooms/{id}/members/{user_id}`
حذف عضو

### `POST /rooms/dm/{user_id}`
شروع گفتگوی خصوصی (DM)

---

## 💬 پیام‌ها

### `GET /messages/{room_id}`
تاریخچه پیام‌ها

**پارامترها:**
- `limit` — تعداد (پیش‌فرض ۵۰)
- `before` — پیام‌های قبل از این id (صفحه‌بندی)

**پاسخ:**
```json
[
  {
    "id": 1,
    "room_id": 1,
    "sender_id": 2,
    "sender_name": "علی رضایی",
    "text": "سلام به همه!",
    "created_at": "2026-08-02T10:30:00Z",
    "edited_at": null,
    "reactions": []
  }
]
```

### `DELETE /messages/{id}`
حذف پیام (مالک یا ادمین)

---

## 📁 فایل‌ها

### `POST /files/upload`
آپلود فایل (multipart/form-data)

### `GET /files/{id}`
دانلود فایل

---

## ⚙️ تنظیمات (ادمین)

### `GET /settings/info`
اطلاعات عمومی سرور

**پاسخ:**
```json
{
  "name": "خان",
  "version": "1.0.0",
  "port": 1727,
  "ip": "192.168.1.100",
  "dns": "",
  "address_type": "ip"
}
```

### `GET /settings/license`
وضعیت لایسنس

**پاسخ:**
```json
{
  "state": "valid",        // free | valid | tampered
  "max_users": 100,
  "licensed_to": "شرکت نمونه",
  "expiry": "2027-12-31T00:00:00Z",
  "error": ""
}
```

### `POST /settings/license`
اعمال لایسنس (multipart, فیلد `license`)

### `DELETE /settings/license`
حذف لایسنس → بازگشت به ۲۰ کاربر

### `GET /settings/network`
تنظیمات شبکه

### `POST /settings/network`
ذخیره تنظیمات شبکه

**بدنه:**
```json
{
  "address_type": "ip",   // ip | dns
  "ip": "192.168.1.100",
  "dns": "",
  "port": 1727
}
```

---

## 🔌 WebSocket

### اتصال

```
ws://<server>:1727/ws?token=<token>
```

### ارسال

```json
{ "type": "send_message", "room_id": 1, "text": "سلام!" }
```

### دریافت

```json
{ "type": "message", "room_id": 1, "payload": { ... } }
```

### رویدادها

| رویداد | جهت | توضیح |
|--------|------|-------|
| `message` | دریافت | پیام جدید |
| `message_edited` | دریافت | پیام ویرایش شد |
| `message_deleted` | دریافت | پیام حذف شد |
| `reaction` | دریافت | واکنش اضافه/حذف شد |
| `typing` | هر دو | در حال نوشتن |
| `presence` | دریافت | آنلاین/آفلاین |
| `force_logout` | دریافت | خروج اجباری |

---

## 🏷 کدهای خطا

| کد | معنی |
|----|------|
| 400 | درخواست نامعتبر |
| 401 | احراز هویت ناموفق |
| 403 | دسترسی غیرمجاز |
| 404 | پیدا نشد |
| 409 | تداخل (نام کاربری تکراری و...) |
| 429 | تلاش زیاد — قفل ۵ دقیقه |
| 500 | خطای سرور |

---
📡 **خان API** — نسخه 1.0
