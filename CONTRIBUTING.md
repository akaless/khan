# Contributing to خان (Khan)

خوش آمدید! 🙏 ممنون که می‌خواهید به خان کمک کنید.

## 🚀 شروع

1. ریپو را fork کنید
2. کلون کنید: `git clone https://codeberg.org/adiib/khan.git`
3. شاخه بسازید: `git checkout -b feature/your-feature`
4. تغییرات را اعمال و commit کنید
5. Push کنید و Pull Request بسازید

## 🐛 گزارش باگ

از [Issues](https://codeberg.org/adiib/khan/issues) استفاده کنید و شامل این موارد باشد:

- نسخه خان (خروجی `khan --version` یا از فایل `code/go.mod`)
- سیستم عامل (Windows/Linux/macOS) و معماری
- قدم‌های بازتولید باگ
- خروجی/خطای مورد انتظار و واقعی
- اسکرین‌شات اگر ممکن است

## ✨ درخواست ویژگی

ویژگی‌های جدید را به‌صورت Issue با برچسب `enhancement` ثبت کنید:

- مشکل را توصیف کنید (چه کاری می‌خواهید انجام دهید)
- راه‌حل پیشنهادی
- اگر ایده‌ای برای پیاده‌سازی دارید، بنویسید

## 🛠 استانداردهای کد

- زبان: **Go 1.22+**
- قالب‌بندی: `gofmt` (اجباری)
- تست: هر تابع عمومی باید تست داشته باشد
- کامنت: توابع عمومی با کامنت انگلیسی توضیح داده شوند
- نام‌گذاری: Go استاندارد (CamelCase)
- از `go vet` قبل از commit استفاده کنید

```bash
# بررسی کیفیت قبل از ارسال
cd code
gofmt -l .          # قالب‌بندی
go vet ./...        # تحلیل استاتیک
go test ./...       # تست‌ها
```

## 📁 ساختار

```
code/
├── cmd/server/          ← نقطه ورود سرور
├── config/              ← تنظیمات
├── internal/
│   ├── models/          ← مدل‌های داده
│   ├── database/        ← دیتاستور
│   ├── repository/      ← لایه دسترسی داده
│   ├── service/         ← منطق کسب‌وکار
│   ├── handler/         ← هندلرهای HTTP
│   └── ws/              ← WebSocket
└── scripts/             ← ابزارها (gen-license)
```

## 📝 Conventional Commits

از [Conventional Commits](https://www.conventionalcommits.org/) استفاده کنید:

```
feat: add new feature
fix: fix a bug
docs: update documentation
style: formatting
refactor: code refactoring
test: add tests
chore: maintenance
```

## 🏷 برچسب‌ها

- `bug` — باگ
- `enhancement` — ویژگی جدید
- `documentation` — مستندات
- `good-first-issue` — مناسب شروع
- `help-wanted` — نیاز به کمک

## ✅ قبل از Pull Request

- [ ] کد `gofmt` شده
- [ ] `go vet` بدون خطا
- [ ] تست‌ها پاس می‌شوند
- [ ] تست جدید برای تغییرات نوشته شده
- [ ] CHANGELOG به‌روز شده

## 🎯 کد نویسی

- کیفیت > سرعت
- سادگی > پیچیدگی
- امنیت اولویت اول
- فارسی/انگلیسی پشتیبانی شود

---
سپاس از مشارکت شما! 🏠🧔
