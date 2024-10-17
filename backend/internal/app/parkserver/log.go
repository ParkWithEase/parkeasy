package parkserver

import (
	"net/http"

	"github.com/rs/zerolog/hlog"
)

// Middleware to handle zerolog
func LogMiddleware(next http.Handler) http.HandlerFunc {
	next = hlog.RequestHandler("request")(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Duplicate the base logger
		log := hlog.FromRequest(r).
			With().
			Logger()
		r = r.WithContext(log.WithContext(r.Context()))

		next.ServeHTTP(w, r)
	})
}
