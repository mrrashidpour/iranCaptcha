# 🇮🇷 Iran Captcha

کپچا با اعداد **فارسی** (و پشتیبانی از اعداد عربی/انگلیسی در نسخه‌های بعدی)  
مناسب برای فرم‌های ثبت‌نام، ورود، نظرسنجی و هر جایی که نیاز به تأیید امنیتی با اعداد فارسی دارید.

---

## ✨ ویژگی‌ها

- نمایش اعداد به صورت **فارسی** (۱۲۳)
- قابلیت ذخیره‌سازی در **حافظه موقت** یا پیاده‌سازی `Store` دلخواه
- پشتیبانی از **فونت زیبا** (B Titr)
- قابلیت تنظیم **تعداد نویز**، طول کپچا، رنگ و اندازه
- خروجی **Base64** برای نمایش مستقیم در HTML
- تابع آماده `Verify` برای بررسی پاسخ کاربر

---

## 📦 نصب

```bash
go get github.com/mrrashidpour/iranCaptcha
```


## 🏗️ ساختار ذخیره‌سازی (Store)
برای استفاده از کپچا، باید یک فضای ذخیره‌سازی (مثل حافظه یا Redis) پیاده‌سازی کنید:

```go
type Store interface {
// Set ارقام مربوط به شناسه کپچا را تنظیم می‌کند.
Set(id string, value string)

// Get ارقام ذخیره شده برای شناسه کپچا را برمی‌گرداند. Clear نشان می‌دهد
// آیا کپچا باید از حافظه حذف شود یا خیر.
Get(id string, clear bool) string

//Verify بررسی درست بودن کپچا  
Verify(id, answer string, clear bool) bool
}
```

### درایورهای آماده:

1. [عداد فارسی](driver_digit.go)  


## 🧠 راه‌اندازی حافظه (Memory Store)
```go
import (
    "time"
    "github.com/mrrashidpour/iranCaptcha"
)

var Store = iranCaptcha.NewMemoryStore(10240, 10 * time.Minute)
// ظرفیت: ۱۰۲۴۰ آیتم
// زمان انقضا: ۱۰ دقیقه
```

## 🎨 ساخت کپچا (Generate)

ابتدا یک درایور بسازید:

```go
driver := iranCaptcha.NewDriverDigit(
    80,                   // ارتفاع
    240,                  // عرض
    20,                   // تعداد نویز
    6,                    // تعداد کاراکترها
    nil,                  // رنگ پس‌زمینه (nil یعنی پیش‌فرض)
    fonts.BTitrBd,        // اسم فونت (BTitrBd.ttf)
)

```

سپس کپچا را تولید کنید:


```go
id, b64s, answer, err := iranCaptcha.Generate(driver, Store)
if err != nil {
    // خطا در تولید
}

// id    : شناسه یکتای کپچا
// b64s  : تصویر کپچا به صورت Base64 (قابل استفاده در img src)
// answer: پاسخ صحیح (برای ذخیره در session یا دیتابیس)
```

## ✅ تأیید کپچا (Verify)

```go
if Store.Verify(captchaId, userInput, true) {
    // کاربر کپچا را درست وارد کرده
} else {
    // کپچا اشتباه است
}
```
پارامتر true یعنی بعد از بررسی، مقدار از حافظه پاک شود.



## 🔧 نیازمندی فونت

```bash
fonts.BTitrBd = "BTitrBd.ttf"
```

## 📄 لایسنس

این پروژه تحت شرایط **لایسنس MIT** منتشر شده است. استفاده، کپی، تغییر، ادغام، انتشار، توزیع، فروش و استفاده از کد آن با ذکر نام نویسنده اصلی، آزاد و بدون مانع است.

```text
MIT License

Copyright (c) 2024 mrrashidpour

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 📞 ارتباط با توسعه‌دهنده

اگر سوالی دارید، نیاز به راهنمایی دارید، یا پیشنهادی برای بهبود پروژه دارید، خوشحال می‌شوم بشنوم.

- 📧 **ایمیل:** [mr.rashidpour@gmail.com](mailto:mr.rashidpour@gmail.com)
- 🐙 **گیت‌هاب:** [github.com/mrrashidpour](https://github.com/mrrashidpour)

می‌توانید از طریق **Issue** در گیت‌هاب هم سوال خود را مطرح کنید تا دیگران نیز از پاسخ آن بهره‌مند شوند.
