package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Config struct {
	// The address to run the server on
	Addr string
	// Whether to run server in insecure mode. This allows cookies to be transferred over plain HTTP.
	Insecure bool
}

// Register all routes
func RegisterRoutes(api huma.API, sessionManager *scs.SessionManager) {
	authMiddleware := routes.NewSessionMiddleware(api, sessionManager)
	api.UseMiddleware(authMiddleware)

	authRepo := authRepo.NewMemoryRepository()
	authService := auth.NewService(authRepo)
	authRoute := routes.NewAuthRoute(authService, sessionManager)

	userRepo := userRepo.NewMemoryRepository()
	userService := user.NewService(*authService, userRepo)
	userRoute := routes.NewUserRoute(userService, sessionManager)
	huma.AutoRegister(api, authRoute)
	huma.AutoRegister(api, userRoute)
}

// Creates a new Huma API instance with routes configured
func (c *Config) NewHumaApi() huma.API {
	router := http.NewServeMux()
	config := huma.DefaultConfig("ParkEasy API", "0.0.0")
	api := humago.New(router, config)
	api.OpenAPI().Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	api.OpenAPI().Components.SecuritySchemes[routes.CookieSecuritySchemeName] = &routes.CookieSecurityScheme
	sessionManager := routes.NewSessionManager(nil)
	sessionManager.Cookie.Secure = !c.Insecure

	RegisterRoutes(api, sessionManager)

	return api
}

// Listen and serve at `addr`.
//
// If `ctx` is cancelled, the server will shutdown gracefully and no error will be returned.
func (c *Config) ListenAndServe(ctx context.Context) error {
	api := c.NewHumaApi()
	huma.NewError = routes.NewErrorFiltered

	srv := http.Server{
		Addr:        c.Addr,
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
