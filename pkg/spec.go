package pkg

import (
	"encoding/json"
	"fmt"
	"io"
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
func LoadSpec(path string) (api.BenchmarkSpec, error) {
	var unmarshalFn func([]byte, interface{}) error

	log.Infof("Loading benchmark specs from '%s'...", path)

	if strings.HasSuffix(path, ".json") {
		unmarshalFn = json.Unmarshal
	} else {
		unmarshalFn = yaml.Unmarshal
	}

	return load(path, unmarshalFn)
}

// LoadSpecFromYamlData loads a spec from the specified slice of bytes, assuming YAML data.
func LoadSpecFromYamlData(data []byte) (spec api.BenchmarkSpec, err error) {
	err = yaml.Unmarshal(data, &spec)

	if err == nil {
		err = validate(spec)
	}

	return spec, err
}

// SaveSpec saves the specified spec to the provided writer in YAML format and closes.
func SaveSpec(spec api.BenchmarkSpec, wc io.WriteCloser) (err error) {
	if err = validate(spec); err != nil {
		return err
	}

	return save(spec, wc)
}

func save(spec api.BenchmarkSpec, wc io.WriteCloser) (err error) {
	var data []byte
	if data, err = yaml.Marshal(spec); err == nil {
		_, err = wc.Write(data)
	}

	if closeErr := wc.Close(); closeErr != nil {
		log.Error(closeErr)
	}

	return err
}

func load(path string, unmarshal func([]byte, interface{}) error) (spec api.BenchmarkSpec, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(path); err == nil {
		err = unmarshal(bytes, &spec)
	}

	if err == nil {
		err = validate(spec)
	}

	return spec, err
}

func validate(spec api.BenchmarkSpec) (err error) {
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
