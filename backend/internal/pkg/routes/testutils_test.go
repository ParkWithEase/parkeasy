package routes

import (
	"context"
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
)

type (
	// A simple adapter for context.Value
	fakeSessionDataGetter struct{}
	fakeSessionDataKey    string
)

// Get implements SessionDataGetter
func (fakeSessionDataGetter) Get(ctx context.Context, key string) any { //nolint: ireturn // this is intentional
	return ctx.Value(fakeSessionDataKey(key))
}

func fakeUserMiddleware(ctx huma.Context, next func(huma.Context)) {
	next(ctx)
}

func jsonAnyify(v any) any { //nolint: ireturn // this is intentional
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var result any
	err = json.Unmarshal(j, &result)
	if err != nil {
		panic(err)
	}

	return result
}
