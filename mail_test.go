package exchangesmtp

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestMail_ToBytes_PlainText(t *testing.T) {
	mail := Mail{
		MT:      PlainText,
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "This is a plain text body.",
	}

	msg, err := mail.ToBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(msg, []byte("Content-Type: text/plain")) {
		t.Errorf("expected Content-Type to be text/plain, got: %s", msg)
	}
}

func TestMail_ToBytes_HTML(t *testing.T) {
	mail := Mail{
		MT:      HTML,
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test HTML Subject",
		Body:    "<h1>This is HTML content.</h1>",
	}

	msg, err := mail.ToBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(msg, []byte("Content-Type: text/html")) {
		t.Errorf("expected Content-Type to be text/html, got: %s", msg)
	}
}

func TestMail_ToBytes_Attachment(t *testing.T) {
	attachmentContent := []byte("This is the content of the file.")

	mail := Mail{
		MT:      PlainText,
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Attachment",
		Body:    "Please see the attached file.",
		Attachment: []AttachmentFile{ // ИСПРАВЛЕНИЕ: правильный тип
			{
				Name: "testfile.txt",
				Body: attachmentContent,
			},
		},
	}

	msg, err := mail.ToBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(msg, []byte("Content-Disposition: attachment; filename=\"testfile.txt\"")) {
		t.Errorf("expected attachment filename in email, got: %s", msg)
	}

	if !bytes.Contains(msg, []byte(base64.StdEncoding.EncodeToString(attachmentContent))) {
		t.Errorf("expected base64 encoded attachment content in email, got: %s", msg)
	}
}

func TestMail_ToBytes_EmptyTo(t *testing.T) {
	mail := Mail{
		MT:      PlainText,
		From:    "sender@example.com",
		Subject: "Test Empty To",
		Body:    "This email has no recipients.",
	}

	_, err := mail.ToBytes()
	if err == nil || err.Error() != "recipient list is empty" {
		t.Errorf("expected error 'recipient list is empty', got: %v", err)
	}
}

func TestMail_ToBytes_EmptyBody(t *testing.T) {
	mail := Mail{
		MT:      PlainText,
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Empty Body",
	}

	_, err := mail.ToBytes()
	if err == nil || err.Error() != "email body is empty" {
		t.Errorf("expected error 'email body is empty', got: %v", err)
	}
}

func TestMail_ToBytes_MultipleRecipients(t *testing.T) {
	mail := Mail{
		MT:      PlainText,
		From:    "sender@example.com",
		To:      []string{"recipient1@example.com", "recipient2@example.com"},
		Subject: "Test Multiple Recipients",
		Body:    "Test body",
	}

	msg, err := mail.ToBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(msg, []byte("To: recipient1@example.com, recipient2@example.com")) {
		t.Errorf("expected comma-separated recipients, got: %s", msg)
	}
}
