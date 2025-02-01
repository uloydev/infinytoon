package appctx

import (
	"context"
)

type AppContext struct {
	ctx context.Context
}

func NewAppContext() *AppContext {
	return &AppContext{
		ctx: context.Background(),
	}
}

func (a *AppContext) Context() context.Context {
	return a.ctx
}

func (a *AppContext) SetContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *AppContext) SpanChild() *AppContext {
	return &AppContext{
		ctx: context.WithoutCancel(a.ctx),
	}
}

func (a *AppContext) Set(key AppContextKey, value interface{}) {
	a.ctx = context.WithValue(a.ctx, key, value)
}

func (a *AppContext) Get(key AppContextKey) interface{} {
	return a.ctx.Value(key)
}

func (a *AppContext) GetStr(key AppContextKey) string {
	return a.ctx.Value(key).(string)
}

func (a *AppContext) GetInt(key AppContextKey) int {
	return a.ctx.Value(key).(int)
}

func (a *AppContext) GetBool(key AppContextKey) bool {
	return a.ctx.Value(key).(bool)
}

func (a *AppContext) GetFloat64(key AppContextKey) float64 {
	return a.ctx.Value(key).(float64)
}
