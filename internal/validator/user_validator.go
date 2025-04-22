package validator

import (
	"test-project/internal/domain"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/ru"
)

type UserValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewUserValidator() (*UserValidator, error) {
	validate := validator.New()

	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)
	trans, _ := uni.GetTranslator("ru")
	err := translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, err
	}

	return &UserValidator{
		validate: validate,
		trans:    trans,
	}, nil
}

func (v *UserValidator) ValidateUser(user domain.User) []string {
	err := v.validate.Struct(user)
	if err == nil {
		return nil
	}

	var messages []string
	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, e.Translate(v.trans))
	}
	return messages
}
