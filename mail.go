package exchangesmtp

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"os"
	"strings"
)

type MailType int

const (
	PlainText MailType = iota
	HTML
)

var mailTypeNames = [...]string{"text/plain", "text/html"}

func (mt MailType) String() string {
	return mailTypeNames[mt]
}

const charset = "UTF-8"

// Mail is a struct for two types of email: plain text and html like.
type Mail struct {
	MT MailType

	From        string
	To          []string
	Subject     string
	Body        string
	contentType string

	Attachment []AttachmentFile
}

type AttachmentFile struct {
	Name        string
	ContentType string
	Body        []byte
}

// generateBoundary creates a random MIME boundary
func generateBoundary() string {
	var buf [16]byte
	rand.Read(buf[:])
	return fmt.Sprintf("boundary-%x", buf[:])
}

func (m *Mail) ToBytes() ([]byte, error) {
	msg := bytes.NewBuffer(nil)

	if len(m.To) == 0 {
		return nil, errors.New("recipient list is empty")
	}

	if len(m.Body) == 0 {
		return nil, errors.New("email body is empty")
	}

	// write headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ", ")))
	sbj := mime.QEncoding.Encode("utf-8", m.Subject)
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", sbj))
	msg.WriteString("MIME-Version: 1.0\r\n")

	boundary := generateBoundary()
	if len(m.Attachment) > 0 {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}

	// write body
	msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=%s\r\n", m.MT.String(), charset))
	msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qp := quotedprintable.NewWriter(msg)
	_, err := qp.Write([]byte(m.Body))
	if err != nil {
		return nil, err
	}
	err = qp.Close()
	if err != nil {
		return nil, err
	}

	// add attachments
	if len(m.Attachment) > 0 {
		for _, file := range m.Attachment {
			msg.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))

			contentType := file.ContentType
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", contentType, file.Name))
			msg.WriteString("Content-Transfer-Encoding: base64\r\n")
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", file.Name))

			if len(file.Body) > 0 {
				if err := m.writeBytes(msg, file.Body); err != nil {
					return nil, err
				}
			} else {
				if err := m.writeFile(msg, file.Name); err != nil {
					return nil, err
				}
			}
		}
		msg.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	}

	return msg.Bytes(), nil
}

func (m *Mail) writeBytes(msg *bytes.Buffer, file []byte) error {
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	msg.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		msg.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			msg.WriteString("\r\n")
		}
	}

	return nil
}

func (m *Mail) writeFile(msg *bytes.Buffer, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err = m.writeBytes(msg, file); err != nil {
		return err
	}

	return nil
}
