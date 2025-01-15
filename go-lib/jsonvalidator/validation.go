package jsonvalidator

import (
	"fmt"
	"strings"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	valid "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Validator struct {
	validate *valid.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: valid.New(),
	}
}

func (v *Validator) ValidateStruct(payload interface{}) *errutils.Error {
	err := v.validate.Struct(payload)
	if err == nil {
		return nil
	}

	var validationErrs valid.ValidationErrors
	if errors.As(err, &validationErrs) { // Safely unwrap and check for ValidationErrors
		var errMsg strings.Builder
		for _, fieldErr := range validationErrs {
			tmp := strings.Split(fieldErr.StructNamespace(), ".")
			msg := fmt.Sprintf("%s is %s", tmp[len(tmp)-1], fieldErr.Tag())
			msg = strings.ToLower(string(msg[0])) + msg[1:]
			errMsg.WriteString(msg + ", ")
		}

		// Trim trailing comma and space
		finalMsg := strings.TrimSuffix(errMsg.String(), ", ")
		return errutils.NewError(errors.New(finalMsg), errutils.BadRequest)
	}

	// Handle non-validation errors
	return errutils.NewError(err, errutils.InternalServerError)
}

// For echo validation
func (v *Validator) Validate(payload interface{}) error {
	err := v.ValidateStruct(payload)
	if err != nil {
		return err.ToEchoError()
	}

	return nil
}
