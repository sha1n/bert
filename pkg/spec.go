package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sha1n/benchy/api"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"gopkg.in/yaml.v2"
)

// LoadSpec loads benchmark specs from the specified file.
func LoadSpec(path string) (*api.BenchmarkSpec, error) {
	var unmarshalFn func([]byte, interface{}) error

	log.Infof("Loading benchmark specs from '%s'...", path)

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		unmarshalFn = yaml.Unmarshal
	} else {
		unmarshalFn = json.Unmarshal
	}

	return load(path, unmarshalFn)
}

func SaveSpec(spec *api.BenchmarkSpec, path string) (err error) {
	if err = validate(spec); err != nil {
		return err
	}

	return save(spec, path)
}

func save(spec *api.BenchmarkSpec, path string) (err error) {
	var data []byte
	if data, err = yaml.Marshal(spec); err == nil {
		err = os.WriteFile(path, data, 0644)
	}

	return err
}

func load(path string, unmarshal func([]byte, interface{}) error) (spec *api.BenchmarkSpec, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(path); err == nil {
		err = unmarshal(bytes, &spec)
	}

	if err == nil {
		err = validate(spec)
	}

	return spec, err
}

func validate(spec *api.BenchmarkSpec) (err error) {
	v := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v, trans)

	if err = v.Struct(spec); err != nil {
		var errstrings []string
		errstrings = append(errstrings, "Invalid configuration:")
		errstrings = append(errstrings, translateError(err, trans)...)
		err = fmt.Errorf(strings.Join(errstrings, "\n\t- "))
	}

	return err
}

func translateError(err error, trans ut.Translator) (errs []string) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		log.Debug(e)
		errs = append(errs, translatedErr.Error())
	}
	return errs
}
