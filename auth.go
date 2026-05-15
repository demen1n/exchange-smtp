package exchangesmtp

import (
	"errors"
	"net/smtp"
)

// loginAuth is analog for plain auth for Exchange server.
// This is smtp.Auth interface implementation.
type loginAuth struct {
	username, password string
}

// Start begins an authentication with a server.
func (a *loginAuth) Start(info *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

// Next continues the authentication.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server: " + string(fromServer))
		}
	}
	return nil, nil
}

// LoginAuth uses for authorization in Exchange.
// Plain auth for Exchange server doesn't work since 2017.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
