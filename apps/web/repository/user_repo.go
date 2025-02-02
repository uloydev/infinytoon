package repository

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/database"
	"infinitoon.dev/infinitoon/shared/schema"
)

type UserRepo interface {
	FindByID(id bson.ObjectID) (user *schema.User, err error)
	FindByEmail(email string) (user *schema.User, err error)
}

type userRepo struct {
	appCtx *appctx.AppContext
	db     *database.MongoDB
}

func NewUserRepo(appCtx *appctx.AppContext) UserRepo {
	return &userRepo{
		appCtx: appCtx,
		db:     database.GetMongoFromCtx(appCtx),
	}
}

func (r *userRepo) getDb() *mongo.Database {
	return r.db.Client().Database("infinitoon")
}

func (r *userRepo) FindByID(id bson.ObjectID) (user *schema.User, err error) {
	user = &schema.User{}
	res := r.getDb().
		Collection(user.CollectionName()).
		FindOne(r.appCtx.Context(), bson.M{"_id": id})

	if err = res.Err(); err != nil {
		return
	}

	err = res.Decode(user)
	return
}

func (r *userRepo) FindByEmail(email string) (user *schema.User, err error) {
	user = &schema.User{}
	res := r.getDb().
		Collection(user.CollectionName()).
		FindOne(r.appCtx.Context(), bson.M{"email": email})

	if err = res.Err(); err != nil {
		return
	}

	err = res.Decode(user)
	return
}
