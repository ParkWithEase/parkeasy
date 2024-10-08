package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

type Config struct {
	// Database pool for Postgres connection
	DBPool *pgxpool.Pool
	// The address to run the server on
	Addr string
	// Whether to run server in insecure mode. This allows cookies to be transferred over plain HTTP.
	Insecure bool
}

// Register all routes
func RegisterRoutes(api huma.API, sessionManager *scs.SessionManager) {
	authMiddleware := routes.NewSessionMiddleware(api, sessionManager)
	api.UseMiddleware(authMiddleware)

	passwordRepository := resettoken.NewMemoryRepository()
	authRepository := authRepo.NewMemoryRepository()
	authService := auth.NewService(authRepository, passwordRepository)
	authRoute := routes.NewAuthRoute(authService, sessionManager)

	userRepository := userRepo.NewMemoryRepository()
	userService := user.NewService(authService, userRepository)
	userRoute := routes.NewUserRoute(userService, sessionManager)
	huma.AutoRegister(api, authRoute)
	huma.AutoRegister(api, userRoute)
}

// Creates a new Huma API instance with routes configured
func (c *Config) NewHumaAPI() huma.API {
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
	api := c.NewHumaAPI()
	huma.NewError = routes.NewErrorFiltered

	srv := http.Server{
		Addr:              c.Addr,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           api.Adapter(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	if c.Insecure {
		corsMiddleware := cors.New(cors.Options{
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPatch,
			},
			// NOTE: This allow all credentials to be passed across CORS
			//
			// It is very insecure, and as such should only be enabled for development
			AllowOriginFunc: func(_ string) bool {
				return true
			},
			AllowCredentials: true,
		})
		srv.Handler = corsMiddleware.Handler(srv.Handler)
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
