package httpx

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	en_lang "github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/suyuan32/simple-admin-core/common/logmessage"
)

const xForwardedFor = "X-Forwarded-For"

// GetFormValues returns the form values.
func GetFormValues(r *http.Request) (map[string]interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return params, nil
}

// GetRemoteAddr returns the peer address, supports X-Forward-For.
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardedFor)
	if len(v) > 0 {
		return v
	}

	return r.RemoteAddr
}

type Validator struct {
	Validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     map[string]ut.Translator
}

func NewValidator() *Validator {
	v := Validator{}
	en := en_lang.New()
	zh := zh_Hans.New()
	v.Uni = ut.New(zh, en, zh)
	v.Validator = validator.New()
	enTrans, _ := v.Uni.GetTranslator("en")
	zhTrans, _ := v.Uni.GetTranslator("zh")
	v.Trans = make(map[string]ut.Translator)
	v.Trans["en"] = enTrans
	v.Trans["zh"] = zhTrans

	err := en_translations.RegisterDefaultTranslations(v.Validator, enTrans)
	if err != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", err.Error()))
		return nil
	}
	err = zh_translations.RegisterDefaultTranslations(v.Validator, zhTrans)
	if err != nil {
		logx.Errorw(logmessage.DatabaseError, logx.Field("Detail", err.Error()))
		return nil
	}

	return &v
}

func (v *Validator) Validate(data interface{}, lang string) string {
	err := v.Validator.Struct(data)
	if err == nil {
		return ""
	}

	errs, ok := err.(validator.ValidationErrors)
	if ok {
		transData := errs.Translate(v.Trans[lang])
		s := strings.Builder{}
		for _, v := range transData {
			s.WriteString(v)
			s.WriteString(" ")
		}
		return s.String()
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return invalid.Error()
	}

	return ""
}
