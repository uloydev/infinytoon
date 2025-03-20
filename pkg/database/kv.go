package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/logger"
)

type KVClient struct {
	appCtx *appctx.AppContext
	client *redis.Client
}

func NewKVClient(appCtx *appctx.AppContext) *KVClient {
	log := appCtx.Get(appctx.LoggerKey).(*logger.Logger)
	log.Info().Msg("initializing kv client")
	kv := &KVClient{
		appCtx: appCtx,
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}),
	}

	// ping the client to check if it's connected
	_, err := kv.client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to kv client")
	}

	appCtx.Set(appctx.KVClientKey, kv)
	log.Info().Msg("kv client initialized")
	return kv
}

func GetKVClientFromCtx(appCtx *appctx.AppContext) *KVClient {
	db := appCtx.Get(appctx.KVClientKey).(*KVClient)
	if db == nil {
		panic("kv client is not initialized")
	}
	return db
}

func (k *KVClient) Client() *redis.Client {
	return k.client
}
