package validator

import (
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/ru"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func New() (*Validator, error) {
	validate := validator.New()

	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)
	trans, _ := uni.GetTranslator("ru")
	if err := translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	return &Validator{
		validate: validate,
		trans:    trans,
	}, nil
}

func (v *Validator) Validate(i interface{}) []string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	var messages []string
	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, e.Translate(v.trans))
	}
	return messages
}
