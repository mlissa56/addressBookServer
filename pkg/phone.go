package pkg

import (
	"errors"
	"strings"
)

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
    var sb strings.Builder
	for i := range phone {
		n := phone[i]
		if n >= '0' && n <= '9' {
			sb.WriteByte(n)
		}
	}
	normalizedPhone = sb.String()

	if sb.Len() == 0 {
		return "", errors.New("Invalid phone number")
	}
	if normalizedPhone[0] == '8' {
		normalizedPhone = "7" + normalizedPhone[1:]
	}
	if normalizedPhone[0] != '7' {
		return "", errors.New("Invalid phone number")
	}
	if len(normalizedPhone) != 11 {
		return "", errors.New("Invalid phone number")
	}

	return normalizedPhone, nil
}
