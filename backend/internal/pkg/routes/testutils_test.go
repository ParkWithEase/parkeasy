package routes

import (
	"context"
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
)

type (
	// A simple adapter for context.Value
	FakeSessionDataGetter struct{}
	FakeSessionDataKey    string
)

// Get implements SessionDataGetter.
func (FakeSessionDataGetter) Get(ctx context.Context, key string) any { //nolint: ireturn // required by interface
	return ctx.Value(FakeSessionDataKey(key))
}

func FakeUserMiddleware(ctx huma.Context, next func(huma.Context)) {
	next(ctx)
}

func JsonAnyify(v any) any { //nolint: ireturn // this is intentional
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