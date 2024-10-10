package routes

import (
	"context"
	"encoding/json"
)

type (
	// A simple adapter for context.Value
	fakeSessionDataGetter struct{}
	fakeSessionDataKey    string
)

// Get implements SessionDataGetter
func (fakeSessionDataGetter) Get(ctx context.Context, key string) any {
	return ctx.Value(fakeSessionDataKey(key))
}

func jsonAnyify(v any) any {
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
