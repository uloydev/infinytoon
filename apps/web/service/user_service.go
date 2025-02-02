package service

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"infinitoon.dev/infinitoon/apps/web/repository"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/shared/schema"
)

type UserService interface {
	GetByID(id string) (user *schema.User, err error)
	GetByEmail(email string) (user *schema.User, err error)
}

type userService struct {
	appCtx   *appctx.AppContext
	userRepo repository.UserRepo
}

func NewUserService(appCtx *appctx.AppContext, userRepo repository.UserRepo) UserService {
	return &userService{
		appCtx:   appCtx,
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(id string) (user *schema.User, err error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	return s.userRepo.FindByID(objId)
}

func (s *userService) GetByEmail(email string) (user *schema.User, err error) {
	return s.userRepo.FindByEmail(email)
}
