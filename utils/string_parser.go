package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const UsernameTooShortErr = "Le nom d'utilisateur doit être supérieur ou égal à 4"
const UsernameTooLongErr = "Le nom d'utilisateur doit être inférieur ou égal à 20 caractères"
const DoubleSpaceForbiddenErr = "L'utilisation d'un espace à la suite d'un autre n'est pas autorisé"
const InvalidEmailErr = "Cet email est invalide"

func UsernameValidity(str string) error {
	if len(str) < 3 {
		return errors.New(UsernameTooShortErr)
	}

	if len(str) > 20 {
		return errors.New(UsernameTooLongErr)
	}

	for _, c := range str {
		if !(c == 32 ||
			(c >= 48 && c <= 57) ||
			(c >= 65 && c <= 90) ||
			(c >= 97 && c <= 122)) {
			return fmt.Errorf("L'utilisation du caractère “%c“ n'est pas autorisé", c)
		}
	}

	for i := 0; i < len(str); i++ {
		if str[i] == 32 && str[i+1] == 32 {
			return errors.New(DoubleSpaceForbiddenErr)
		}
	}

	return nil
}

func EmailValidity(email string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if emailRegex.MatchString(email) {
		return nil
	}
	return errors.New(InvalidEmailErr)
}

func TokenGenerator(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	const alphaNumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	sb := strings.Builder{}
	sb.Grow(length)

	for i := 0; i < sb.Cap(); i++ {
		sb.WriteByte(alphaNumeric[rand.Intn(len(alphaNumeric)-1)])
	}

	return sb.String()
}
