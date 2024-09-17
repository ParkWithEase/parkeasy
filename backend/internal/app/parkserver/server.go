package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// Register all routes
func RegisterRoutes(api huma.API) {
	huma.AutoRegister(api, routes.NewGreetingRoute(&services.SimpleGreeting{}))
}

// Creates a new Huma API instance with routes configured
func NewHumaApi() huma.API {
	router := http.NewServeMux()
	config := huma.DefaultConfig("Greeting API", "0.0.0")
	api := humago.New(router, config)
	RegisterRoutes(api)

	return api
}

// Listen and serve at `addr`.
//
// If `ctx` is cancelled, the server will shutdown gracefully and no error will be returned.
func ListenAndServe(ctx context.Context, addr string) error {
	api := NewHumaApi()

	srv := http.Server{
		Addr:        addr,
		BaseContext: func(net.Listener) context.Context { return ctx },
		Handler:     api.Adapter(),
	}

	go func() {
		<-ctx.Done()
		// Ignore shutdown errors, not that we can do anything about them
		_ = srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	// ServerClosed just meant that the server has shutdown
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
