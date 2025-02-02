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

func (b *Base) GetID() string {
	return b.ID.Hex()
}

func (b *Base) SetID(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	b.ID = objectID
	return nil
}
