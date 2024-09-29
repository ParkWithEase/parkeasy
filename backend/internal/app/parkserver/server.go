package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	// The address to run the server on
	Addr string
	// Whether to run server in insecure mode. This allows cookies to be transferred over plain HTTP.
	Insecure bool
}

// Register all routes
func RegisterRoutes(api huma.API, sessionManager *scs.SessionManager, dbPool *pgxpool.Pool) {
	authMiddleware := routes.NewSessionMiddleware(api, sessionManager)
	api.UseMiddleware(authMiddleware)

	authRepository := authRepo.NewMemoryRepository()
	authService := auth.NewService(authRepository)
	authRoute := routes.NewAuthRoute(authService, sessionManager)

	userRepository := userRepo.NewMemoryRepository()
	userService := user.NewService(authService, userRepository)
	userRoute := routes.NewUserRoute(userService, sessionManager)
	huma.AutoRegister(api, authRoute)
	huma.AutoRegister(api, userRoute)

	carService := car.NewService(dbPool)
	carRoute := routes.NewCarRoute(carService)
	huma.AutoRegister(api, carRoute)

}

// Creates a new Huma API instance with routes configured
func (c *Config) NewHumaAPI(dbPool *pgxpool.Pool) huma.API {
	router := http.NewServeMux()
	config := huma.DefaultConfig("ParkEasy API", "0.0.0")
	api := humago.New(router, config)
	api.OpenAPI().Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	api.OpenAPI().Components.SecuritySchemes[routes.CookieSecuritySchemeName] = &routes.CookieSecurityScheme
	sessionManager := routes.NewSessionManager(nil)
	sessionManager.Cookie.Secure = !c.Insecure

	RegisterRoutes(api, sessionManager, dbPool) // Pass dbPool to register car routes

	return api
}

// Listen and serve at `addr`.
//
// If `ctx` is cancelled, the server will shutdown gracefully and no error will be returned.
func (c *Config) ListenAndServe(ctx context.Context, dbPool *pgxpool.Pool) error {
	api := c.NewHumaAPI(dbPool)
	huma.NewError = routes.NewErrorFiltered

	srv := http.Server{
		Addr:              c.Addr,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           api.Adapter(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
