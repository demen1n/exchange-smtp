# exchangesmtp

Simple Go library for sending emails through Microsoft Exchange servers using LOGIN authentication.

## Why?

Exchange dropped support for PLAIN auth back in 2017. This library implements LOGIN auth mechanism that still works.

## Install

```bash
go get github.com/yourusername/exchangesmtp
```

## Usage

### Quick send

```go
qs := exchangesmtp.NewQuickSender(
    "user@company.com",           // username
    "password",                    // password
    "smtp.company.com:587",        // server
    "sender@company.com",          // from
    []string{"recipient@company.com"}, // to
)

err := qs.Send("Subject", "Body text")
```

### With attachments

```go
ms := exchangesmtp.NewMailSender("user@company.com", "password", "smtp.company.com:587")

mail := exchangesmtp.Mail{
    MT:      exchangesmtp.HTML,
    From:    "sender@company.com",
    To:      []string{"recipient@company.com"},
    Subject: "Check this out",
    Body:    "<h1>Hello!</h1>",
    Attachment: []exchangesmtp.AttachmentFile{
        {
            Name: "report.pdf",
            Body: fileBytes,
        },
    },
}

err := ms.Send(mail)
```

## Features

- LOGIN authentication for Exchange
- Plain text and HTML emails
- File attachments
- UTF-8 support
- Quoted-printable encoding

## License

MIT