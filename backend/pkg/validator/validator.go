package validator

import (
	"algoplatform/internal/domain"

	v10 "github.com/go-playground/validator/v10"
)

type V struct {
	v *v10.Validate
}

func New() *V                   { return &V{v: v10.New()} }
func (v *V) Struct(x any) error { return v.v.Struct(x) }

var _ domain.Validator = (*V)(nil)
