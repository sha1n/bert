package specs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/sha1n/bert/api"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"gopkg.in/yaml.v2"
)

// LoadSpec loads benchmark specs from the specified file.
func LoadSpec(path string) (api.BenchmarkSpec, error) {
	var unmarshalFn func([]byte, interface{}) error

	slog.Info(fmt.Sprintf("Loading benchmark specs from '%s'...", path))

	if strings.HasSuffix(path, ".json") {
		unmarshalFn = json.Unmarshal
	} else {
		unmarshalFn = yaml.Unmarshal
	}

	return load(path, unmarshalFn)
}

// CreateSpecFrom creates a spec from the specified parameters.
//
// Returns an error if the specified 'executions' is non-positive.
func CreateSpecFrom(executions int, alternate bool, failFast bool, commands ...api.CommandSpec) (spec api.BenchmarkSpec, err error) {
	if executions < 1 {
		return spec, errors.New("executions must be positive")
	}

	spec = api.BenchmarkSpec{
		Executions: executions,
		Alternate:  alternate,
		FailFast:   failFast,
		Scenarios:  []api.ScenarioSpec{},
	}

	for i := range commands {
		command := commands[i]
		scenario := api.ScenarioSpec{
			Name:    fmt.Sprintf("[%s]", strings.Join(command.Cmd, " ")),
			Command: &command,
		}
		spec.Scenarios = append(spec.Scenarios, scenario)
	}

	return
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
		slog.Error(closeErr.Error())
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
		err = errors.New(strings.Join(errstrings, "\n\t- "))
	}

	return err
}

func translateError(err error, trans ut.Translator) (errs []string) {
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := errors.New(e.Translate(trans))
		slog.Debug(e.Error())
		errs = append(errs, translatedErr.Error())
	}
	return errs
}
