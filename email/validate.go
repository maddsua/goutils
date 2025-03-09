package email

import (
	"errors"
	"net"
	"net/mail"
	"strings"
)

func ValidateAddress(email string) error {

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	if addr.Name != "" {
		return errors.New("input is not a single email address")
	}

	_, host, ok := strings.Cut(addr.Address, "@")
	if !ok {
		return errors.New("email host empty")
	}

	if _, err := net.LookupMX(host); err != nil {
		return errors.New("email host unreachable")
	}

	return nil
}
