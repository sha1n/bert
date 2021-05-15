package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"gopkg.in/yaml.v2"
)

// CommandSpec benchmark command execution specs
type CommandSpec struct {
	WorkingDirectory string   `json:"workingDir" yaml:"workingDir"`
	Cmd              []string `json:"cmd" yaml:"cmd" validate:"required"`
}

// ScenarioSpec benchmark scenario specs
type ScenarioSpec struct {
	Name             string `json:"name" yaml:"name" validate:"required"`
	WorkingDirectory string `json:"workingDir" yaml:"workingDir"`
	Env              map[string]string
	BeforeAll        *CommandSpec `json:"beforeAll" yaml:"beforeAll"`
	AfterAll         *CommandSpec `json:"afterAll" yaml:"afterAll"`
	BeforeEach       *CommandSpec `json:"beforeEach" yaml:"beforeEach"`
	AfterEach        *CommandSpec `json:"afterEach" yaml:"afterEach"`
	Command          *CommandSpec `validate:"required,dive"`
}

// BenchmarkSpec benchmark specs top level structure
type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required,dive"`
	Executions int             `validate:"required,gte=1"`
	Alternate  bool
}

// ID returns a unique identifier
func (s *ScenarioSpec) ID() string {
	return s.Name
}

// Load loads benchmark specs from the specified file.
func Load(path string) (*BenchmarkSpec, error) {
	var unmarshalFn func([]byte, interface{}) error

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		unmarshalFn = yaml.Unmarshal
	} else {
		unmarshalFn = json.Unmarshal
	}

	return load(path, unmarshalFn)
}

func load(path string, unmarshal func([]byte, interface{}) error) (spec *BenchmarkSpec, err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(path); err == nil {
		defer jsonFile.Close()

		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			err = unmarshal(bytes, &spec)
		}

		if err == nil {
			err = validate(spec)
		}
	}

	return spec, err
}

func validate(spec *BenchmarkSpec) (err error) {
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
