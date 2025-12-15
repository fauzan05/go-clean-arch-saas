# Email Service Documentation

## Overview

The email service uses SMTP to send verification emails to users. It supports two modes:
- **Development Mode**: Emails are logged to console (not sent)
- **Production Mode**: Emails are sent via SMTP

## Architecture

### Async Pattern: Goroutine (Non-blocking)

Emails are sent **asynchronously using goroutines**, not channels. Reasons:

1. **Fire-and-forget pattern**: We don't need to wait for email sending results
2. **Non-blocking**: Registration/resend doesn't fail if email errors occur
3. **Simple**: No need for worker pool or channel complexity
4. **Graceful degradation**: Errors are only logged, not thrown as exceptions

### Implementation

```go
// Di auth_usecase.go - Register method
go func() {
    if err := u.EmailService.SendVerificationEmail(user.Email, user.Name, token, u.BaseURL); err != nil {
        u.Log.Warnf("Failed to send verification email to %s: %+v", user.Email, err)
    } else {
        u.Log.Infof("Verification email sent to %s", user.Email)
    }
}()
```

**Why goroutine instead of channels?**

| Approach | Use Case | Pros | Cons |
|----------|----------|------|------|
| **Goroutine** (Current) | Fire-and-forget background task | Simple, no coordination needed | No result tracking, no rate limiting |
| **Channel + Worker Pool** | High-volume email queue | Rate limiting, retry logic, ordered processing | Complex, overkill for verification emails |
| **Synchronous** | Critical emails that must succeed | Guaranteed delivery attempt | Blocks registration, bad UX |

**Decision**: Goroutine is the best choice for verification emails because:
- Email verification is not a critical path (users can resend)
- No coordination needed between emails
- Simple & maintainable

### Future Improvements (If Needed)

If advanced features are required later:

```go
// Option 1: Channel-based worker pool (for high-volume)
type EmailJob struct {
    To      string
    Subject string
    Body    string
}

emailQueue := make(chan EmailJob, 1000)

// Worker pool
for i := 0; i < 10; i++ {
    go emailWorker(emailQueue)
}

// Send to queue
emailQueue <- EmailJob{...}

// Option 2: Background job with Redis/RabbitMQ
// - For mission-critical emails
// - With retry mechanism
// - For distributed systems
```

## Template System

### HTML Templates

Templates are stored in `pkg/email/templates/*.html` and embedded into the binary using Go 1.16+ `embed.FS`.

**Benefits of embed.FS:**
- Templates included in binary (no external files needed)
- Zero dependencies at deployment
- Fast template loading (from memory)

### Template Structure

```
pkg/email/
├── email.go              # Email service
└── templates/
    └── verify_email.html # Email verification template
```

### Adding New Templates

1. Create `.html` file in `templates/`
2. Use Go template syntax `{{.Variable}}`
3. Load with `template.ParseFS(templateFS, "templates/your_template.html")`

Example for welcome email:

```html
<!-- templates/welcome.html -->
<!DOCTYPE html>
<html>
<body>
    <h1>Welcome {{.UserName}}!</h1>
    <p>Your account has been verified.</p>
</body>
</html>
```

```go
// Di email.go
func (s *EmailService) SendWelcomeEmail(toEmail, userName string) error {
    tmpl, _ := template.ParseFS(templateFS, "templates/welcome.html")
    
    data := struct {
        UserName string
    }{UserName: userName}
    
    var body bytes.Buffer
    tmpl.Execute(&body, data)
    
    return s.send(toEmail, "Welcome!", body.String())
}
```

## SMTP Configuration

### Gmail Setup

1. Enable 2FA: https://myaccount.google.com/security
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Configure `.env`:
```bash
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-16-char-app-password
EMAIL_FROM=noreply@yourdomain.com
```

### SendGrid

```bash
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=your-sendgrid-api-key
EMAIL_FROM=noreply@yourdomain.com
```

### AWS SES

```bash
EMAIL_HOST=email-smtp.us-east-1.amazonaws.com
EMAIL_PORT=587
EMAIL_USERNAME=your-ses-smtp-username
EMAIL_PASSWORD=your-ses-smtp-password
EMAIL_FROM=verified@yourdomain.com
```

## Development Mode

Leave email config empty for development:

```bash
EMAIL_HOST=
EMAIL_USERNAME=
```

Emails will be logged to console:

```
WARN Email service not configured. Would send to test@example.com: Verify Your Email Address
INFO Email body:
<!DOCTYPE html>...
```

Verification token is available in logs, can be copied for testing.

## Error Handling

Email service uses **graceful degradation**:

1. **SMTP connection error**: Log error, return error (but registration succeeds)
2. **Template error**: Log error, return error
3. **Invalid recipient**: SMTP server returns error

Registration **doesn't fail** if email errors occur because:
- Email is not a critical path
- Users can resend verification
- Better UX (fast registration)

## Testing

Email service is tested with 7 test cases:

1. ✅ Register with email_verified=false
2. ✅ Verify with valid token
3. ✅ Verify with invalid token
4. ✅ Verify already verified email
5. ✅ Resend verification
6. ✅ Resend to non-existent email
7. ✅ Resend to already verified email

Run tests:
```bash
make test
# or
go test ./test -v
```

## Security

- **Token**: 64-character random hex (crypto/rand)
- **Token expiry**: No expiry (token cleared after use)
- **Rate limiting**: None (consider adding if abuse detected)
- **Token reuse**: Token is cleared after verification (cannot be reused)

## Performance

- **Goroutine overhead**: ~2KB per goroutine (negligible)
- **Template parsing**: Cached in memory (fast)
- **SMTP connection**: ~100-300ms (doesn't block registration)
- **Email delivery**: Asynchronous (users don't wait)

## Monitoring

Log events:
- `INFO Verification email sent to {email}` - Success
- `WARN Failed to send verification email to {email}` - Error
- `WARN Email service not configured` - Development mode

Production monitoring:
- Track email send success/failure rate
- Alert if failure rate > 5%
- Monitor SMTP connection errors
