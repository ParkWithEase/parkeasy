package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/stephenafamo/bob"

	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"

	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"

	carRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"

	parkingSpotRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/parkingspot"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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
func (c *Config) RegisterRoutes(api huma.API, sessionManager *scs.SessionManager) {
	db := bob.NewDB(stdlib.OpenDBFromPool(c.DBPool))

	passwordRepository := resettoken.NewMemoryRepository()
	authRepository := authRepo.NewPostgres(db)
	authService := auth.NewService(authRepository, passwordRepository)
	authRoute := routes.NewAuthRoute(authService, sessionManager)

	userRepository := userRepo.NewPostgres(db)
	userService := user.NewService(authService, userRepository)
	userRoute := routes.NewUserRoute(userService, sessionManager)

	parkingSpotRepository := parkingSpotRepo.NewPostgres(db)
	parkingSpotService := parkingspot.New(parkingSpotRepository)
	parkingSpotRoute := routes.NewParkingSpotRoute(parkingSpotService, sessionManager)

	carRepository := carRepo.NewPostgres(db)
	carService := car.New(carRepository)
	carRoute := routes.NewCarRoute(carService, sessionManager)

	routes.UseHumaMiddlewares(api, sessionManager, userService)
	huma.AutoRegister(api, authRoute)
	huma.AutoRegister(api, userRoute)
	huma.AutoRegister(api, parkingSpotRoute)
	huma.AutoRegister(api, carRoute)
}

// Creates a new Huma API instance with routes configured
func (c *Config) NewHumaAPI() huma.API {
	router := http.NewServeMux()
	config := routes.NewHumaConfig()
	api := humago.New(router, config)
	sessionManager := routes.NewSessionManager(pgxstore.New(c.DBPool))
	sessionManager.Cookie.Secure = !c.Insecure

	c.RegisterRoutes(api, sessionManager)

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
