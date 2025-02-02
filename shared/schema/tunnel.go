package schema

import "go.mongodb.org/mongo-driver/v2/bson"

type TunnelStatus string

const (
	TunnelStatusActive   TunnelStatus = "active"
	TunnelStatusInactive TunnelStatus = "inactive"
	TunnelStatusError    TunnelStatus = "error"
)

type Tunnel struct {
	Base
	User      bson.ObjectID `bson:"user"`
	Name      string        `bson:"name"`
	Domain    string        `bson:"domain"`
	LocalIP   string        `bson:"localIP"`
	LocalPort string        `bson:"localPort"`
	Status    TunnelStatus  `bson:"status"`
}

func (t *Tunnel) CollectionName() string {
	return "tunnels"
}
