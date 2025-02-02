package validation

import (
	"github.com/go-playground/validator/v10"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

func InitValidator(appCtx *appctx.AppContext) *validator.Validate {
	validate := validator.New()
	appCtx.Set(appctx.ValidatorKey, validate)
	return validate
}
