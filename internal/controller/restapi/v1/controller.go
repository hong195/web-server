package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/hong195/web-server/internal/usecase"
	"github.com/hong195/web-server/pkg/logger"
)

type V1 struct {
	l     logger.Interface
	v     *validator.Validate
	user  usecase.User
	items usecase.Items
}
