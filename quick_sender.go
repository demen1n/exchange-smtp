package exchangesmtp

type QuickSender struct {
	ms   *MailSender
	from string
	to   []string
}

func NewQuickSender(user, password, server, from string, to []string) *QuickSender {
	return &QuickSender{
		ms:   NewMailSender(user, password, server),
		from: from,
		to:   to,
	}
}
func (qs *QuickSender) Send(subject, body string) error {
	m := Mail{
		MT:      PlainText,
		From:    qs.from,
		To:      qs.to,
		Subject: subject,
		Body:    body,
	}
	return qs.ms.Send(m)
}
