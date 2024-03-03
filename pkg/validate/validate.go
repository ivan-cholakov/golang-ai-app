package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*]?$`)
)

type RuleFunc func() RuleSet

type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}

type Fields map[string][]RuleSet

type Messages map[string]string

func Equal(n string) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "equal",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return str == n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s sould be equal", set.FieldName)
			},
		}
	}
}

func Password() RuleSet {
	return RuleSet{
		Name: "password",
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s should be valid", set.FieldName)
		},
		ValidateFunc: func(rule RuleSet) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			_, ok = ValidatePassword(str)
			return ok
		},
	}
}

func Required() RuleSet {
	return RuleSet{
		Name: "required",
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s is a required field", set.FieldName)
		},
		ValidateFunc: func(rule RuleSet) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) > 0
		},
	}
}

func Message(msg string) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "message",
			RuleValue: msg}
	}
}

func Url() RuleSet {
	return RuleSet{
		Name: "url",
		MessageFunc: func(set RuleSet) string {
			return "not a valid url"
		},
		ValidateFunc: func(set RuleSet) bool {
			u, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return urlRegex.MatchString(u)
		},
	}
}

func Email() RuleSet {
	return RuleSet{
		Name: "email",
		MessageFunc: func(set RuleSet) string {
			return "email address is invalid"
		},
		ValidateFunc: func(set RuleSet) bool {
			email, ok := set.FieldValue.(string)
			if !ok {
				return false
			}

			return emailRegex.MatchString(email)
		},
	}
}

func Max(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "min",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) <= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s sould be maxiumum %d characters long", set.FieldName, n)
			},
		}
	}
}

func Min(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{Name: "min", RuleValue: n, ValidateFunc: func(set RuleSet) bool {
			str, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) >= n
		},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s sould be at least %d characters long", set.FieldName, n)
			},
		}
	}
}

func Rules(rules ...RuleFunc) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]()
	}
	return ruleSets
}

type Validator struct {
	data   any
	fields Fields
}

func New(data any, fields Fields) *Validator {
	return &Validator{
		fields: fields,
		data:   data,
	}
}
func Validate(in any, out any, fields Fields) bool {
	return true
}

func (v *Validator) Validate(target any) bool {
	ok := true
	for fieldName, ruleSets := range v.fields {
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v.data, fieldName)
		for _, set := range ruleSets {
			set.FieldValue = fieldValue
			set.FieldName = fieldName
			if set.Name == "message" {
				setErrorMessage(target, fieldName, set.RuleValue.(string))
				continue
			}
			if !set.ValidateFunc(set) {
				msg := set.MessageFunc(set)
				setErrorMessage(target, fieldName, msg)
				ok = false
			}
		}
	}
	return ok
}

func setErrorMessage(v any, fieldName string, msg string) {
	if v == nil {
		return
	}
	switch t := v.(type) {
	case map[string]string:
		t[fieldName] = msg
	default:
		structVal := reflect.ValueOf(v)
		if structVal.Kind() != reflect.Ptr || structVal.IsNil() {
			return
		}
		structVal = structVal.Elem()
		field := structVal.FieldByName(fieldName)
		field.Set(reflect.ValueOf(msg))
	}
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil

	}

	return fieldVal.Interface()
}

func ValidatePassword(password string) (string, bool) {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) < 8 {
		return "Password must contain at least 8 characters", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain an upper case character", false
	}
	if !hasLower {
		return "Password must contain a lower case character", false
	}
	if !hasNumber {
		return "Password must contain a number", false
	}
	if !hasSpecial {
		return "Password must contain a special character", false
	}
	return "", true

}
