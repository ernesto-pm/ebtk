package files

import (
	"errors"
	"github.com/ernesto-pm/ebtk/pkg/reflect"
	"strings"
)

func TransformFile[T any](f EpiFile) (T, error) {
	var result T

	switch any(result).(type) {
	case JSONEpiFile:
		if strings.HasSuffix(f.Extension, ".json") {
			newInstance := *new(T)
			err := reflect.CopyFields(f, newInstance)
			if err != nil {
				return *new(T), errors.New("couldn't copy fields over to the new instance")
			}
			return newInstance, nil
		}
	default:
		return *new(T), errors.New("unsupported file transformation")
	}

	return *new(T), errors.New("file cannot be transformed")
}
