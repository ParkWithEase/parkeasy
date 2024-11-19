package parkserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/geocoding"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/resettoken"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/routes"
	"github.com/sourcegraph/conc"
	"github.com/stephenafamo/bob"

	authRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/auth"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/health"

	userRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/user"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/user"

	carRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/car"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/car"

	parkingSpotRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/parkingspot"

	demoRepo "github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/demo"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/demo"

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
	// The prefix to the API endpoint
	APIPrefix string
	// Geocodio API key
	GeocodioAPIKey string
	// The address to run the server on
	Addr string
	// The origin to allow cross-origin request from.
	CorsOrigin string
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

	geocodioRepository := geocoding.NewGeocodio(http.DefaultClient, c.GeocodioAPIKey)

	parkingSpotRepository := parkingSpotRepo.NewPostgres(db)
	parkingSpotService := parkingspot.New(parkingSpotRepository, geocodioRepository)
	parkingSpotRoute := routes.NewParkingSpotRoute(parkingSpotService, sessionManager)

	carRepository := carRepo.NewPostgres(db)
	carService := car.New(carRepository)
	carRoute := routes.NewCarRoute(carService, sessionManager)

	demoRepository := demoRepo.NewPostgres(db)
	demoServicer := demo.New(demoRepository)
	demoRoute := routes.NewDemoRoute(demoServicer, sessionManager)

	healthService := health.New(c.DBPool)
	healthRoute := routes.NewHealthRoute(healthService)

	routes.UseHumaMiddlewares(api, sessionManager, userService)
	huma.AutoRegister(api, authRoute)
	huma.AutoRegister(api, userRoute)
	huma.AutoRegister(api, parkingSpotRoute)
	huma.AutoRegister(api, carRoute)
	huma.AutoRegister(api, demoRoute)
	huma.AutoRegister(api, healthRoute)
}

// Creates a new Huma API instance with routes configured
func (c *Config) NewHumaAPI() huma.API {
	router := http.NewServeMux()
	config := routes.NewHumaConfig()

	// Swap /docs handler for Scalar
	config.DocsPath = ""
	router.HandleFunc("/docs", handleDocs)

	api := humago.New(router, config)
	sessionManager := routes.NewSessionManager(pgxstore.New(c.DBPool))
	sessionManager.Cookie.Secure = !c.Insecure
	if c.CorsOrigin != "" {
		sessionManager.Cookie.SameSite = http.SameSiteNoneMode
	}

	if c.APIPrefix != "" {
		api.OpenAPI().Servers = append(api.OpenAPI().Servers, &huma.Server{
			URL:         c.APIPrefix,
			Description: "API server endpoint",
		})
	}

	c.RegisterRoutes(api, sessionManager)

	return api
}

// Listen and serve at `addr`.
//
// If `ctx` is cancelled, the server will shutdown gracefully and no error will be returned.
func (c *Config) ListenAndServe(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg conc.WaitGroup
	defer wg.Wait()

	api := c.NewHumaAPI()

	srv := http.Server{
		Addr:              c.Addr,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           api.Adapter(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	srv.Handler = LogMiddleware(srv.Handler)

	if c.Insecure || c.CorsOrigin != "" {
		corsOpts := cors.Options{
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPatch,
			},
			AllowCredentials: true,
			ExposedHeaders:   []string{"Link"},
		}

		if c.Insecure {
			// NOTE: This allow all credentials to be passed across CORS
			//
			// It is very insecure, and as such should only be enabled for development
			corsOpts.AllowOriginFunc = func(_ string) bool {
				return true
			}
		} else {
			corsOpts.AllowedOrigins = append(corsOpts.AllowedOrigins, c.CorsOrigin)
		}
		corsMiddleware := cors.New(corsOpts)
		srv.Handler = corsMiddleware.Handler(srv.Handler)
	}

	wg.Go(func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// Ignore shutdown errors, not that we can do anything about them
		_ = srv.Shutdown(shutdownCtx)
	})

	err := srv.ListenAndServe()
	// ServerClosed just meant that the server has shutdown

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
