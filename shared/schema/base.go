package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Base struct {
	ID        bson.ObjectID `bson:"_id"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	DeletedAt *time.Time    `bson:"deletedAt,omitempty"`
}
