package validation

import (
	"errors"
	"fmt"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/constants"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrEmailBadFormat   = errors.New("invalid email format")
	ErrPhoneBadFormat   = errors.New("invalid phone format")
	ErrConfirmFormat    = errors.New("password not match with confirm password")
	ErrPassLengh        = errors.New("invalid password length")
	ErrUnresolvableHost = errors.New("unresolvable host")

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])+.(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*.[a-z]{2,4}$")
	phoneRegexp = regexp.MustCompile("^(62)+[0-9]+$")
)

func ValidateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrEmailBadFormat
	}
	return nil
}

func ValidatePhoneFormat(phone string) error {
	if !phoneRegexp.MatchString(phone) {
		return ErrPhoneBadFormat
	}
	return nil
}

func CheckUserManagement(req *model.CreateUserReq) error {
	if len(req.Email) > 0 {
		if err := ValidateFormat(req.Email); err != nil {
			return err
		}
	}

	if len(req.Phone) > 0 {
		if err := ValidatePhoneFormat(req.Phone); err != nil {
			return err
		}
	}

	if len(req.Password) > 0 {
		if req.Password != req.PasswordConfirm {
			return ErrConfirmFormat
		}

		if err := VerifyPassword(req.Password); err != nil {
			return err
		}
	}

	return nil
}

func VerifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(constants.MinPassLength <= passLen && passLen <= constants.MaxPassLength) {
		appendError(fmt.Sprintf("password length must be %d characters long or more", constants.MinPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}

