package translator

import (
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type (
	Translator interface {
		// Translate struct fields validation error
		TranslateValidationError(ctx *fiber.Ctx, fields validator.FieldError, plurals interface{})

		// Translate Arbitrary message
		TranslateMessage(ctx *fiber.Ctx, key string, param TranslateParam, plurals interface{}) string

		// Validate struct fileds validation error helper
		ValidateRequest(ctx *fiber.Ctx, validator *validator.Validate, validateStruct interface{}) []string
	}

	// Translate parameter
	TranslateParam map[string]interface{}
)

type UTtrans struct {
	fs.FileInfo
	bundle *i18n.Bundle
}

func NewUtTrans(fs embed.FS, transFolderName string) (*UTtrans, error) {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	dirs, err := fs.ReadDir(transFolderName)
	if err != nil {

		return nil, err
	}

	for _, file := range dirs {
		file_path := filepath.Join(transFolderName, file.Name())

		_, err = bundle.LoadMessageFileFS(fs, file_path)
		if err != nil {
			return nil, err

		}
	}

	return &UTtrans{bundle: bundle}, err
}

func (u *UTtrans) TranslateValidationError(
	locale string,
	fe validator.FieldError,
	plurals interface{},
) string {

	localizer := i18n.NewLocalizer(u.bundle, locale)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: fe.Tag(),
		TemplateData: map[string]interface{}{
			"Field": fe.Field(),
			"Param": fe.Param(),
		},
		PluralCount: plurals,
	})
	if err != nil {
		return "Error Translate Message"
	}

	return message
}

// Translate message with associated key , plurals define the plurals variable that determine more than one form
func (u *UTtrans) TranslateMessage(locale string,
	key string,
	para TranslateParam,
	plurals interface{},
) string {
	localizer := i18n.NewLocalizer(u.bundle, locale)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: para,
		PluralCount:  plurals,
	})
	if err != nil {
		return "Error TranslateMessage"
	}
	return message
}

// Translate Validate Error when parsing request body to a more user friendly form or other lang depend on locale store on url query
func (u *UTtrans) ValidateRequest(
	locale string,
	v *validator.Validate,
	validateStruct interface{},
) []string {
	var return_err []string

	validate_errs := v.Struct(validateStruct)

	if validate_errs != nil {
		for _, err := range validate_errs.(validator.ValidationErrors) {
			err_msg := u.TranslateValidationError(locale, err, nil)
			return_err = append(return_err, err_msg)

		}
		return return_err
	}
	return nil
}
