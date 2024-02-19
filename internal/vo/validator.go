package vo

import (
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

var Validate *validator.Validate
var trans ut.Translator

type ValidateError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}
type ValidateErrors struct {
	Validates []*ValidateError `json:"validates"`
}

func init() {

	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ = uni.GetTranslator("en")

	Validate = validator.New()

	_ = en_translations.RegisterDefaultTranslations(Validate, trans)

	// register function to get tag name from json tags.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterValidation("hexAddress", hexAddressValidate)
	Validate.RegisterValidation("txHash", txHashValidate)

}

func ValidateStruct(s interface{}) error {
	return _newValidateErrors(Validate.Struct(s))
}

func _newValidateErrors(err error) error {
	if err == nil {
		return nil
	}
	var es = &ValidateErrors{Validates: []*ValidateError{}}
	for _, err := range err.(validator.ValidationErrors) {
		var element ValidateError
		element.Field = err.Field()
		element.Tag = err.Tag()
		element.Value = err.Param()
		element.Message = err.Translate(trans)
		es.Validates = append(es.Validates, &element)
	}
	return es
}

func (v *ValidateErrors) Error() string {
	ms := make([]string, len(v.Validates))

	for _, validate := range v.Validates {
		ms = append(ms, validate.Message)
	}
	return strings.Join(ms, "\n")
}

func hexAddressValidate(fl validator.FieldLevel) bool {
	return utils.IsHexAddress(fl.Field().String())
}

func txHashValidate(fl validator.FieldLevel) bool {
	return utils.IsHashHex(fl.Field().String())
}
