package routes

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// OpenAPI cookie security scheme for the API
var CookieSecurityScheme = huma.SecurityScheme{
	Type: "apiKey",
	In:   "cookie",
	Name: "session",
}

const (
	CookieSecuritySchemeName = "cookieAuth"
	SessionKeyAuthID         = "authid"
	SessionKeyUserID         = "userid"
	SessionKeyPersist        = "persist"
	DefaultSessionLifetime   = 30 * 24 * time.Hour
)

type SessionDataGetter interface {
	Get(ctx context.Context, key string) any
}

type SessionDataPutter interface {
	Put(ctx context.Context, key string, data any)
}

type SessionDataGetterPutter interface {
	SessionDataGetter
	SessionDataPutter
}

// Headers to commit into the result
type SessionHeaderOutput struct {
	SetCookie    []string `header:"Set-Cookie"`
	CacheControl []string `header:"Cache-Control"`
}

// Returns a session manager configured for use with this package.
//
// `store` can be nil, in which case the default memory store is used.
func NewSessionManager(store scs.Store) *scs.SessionManager {
	result := scs.New()
	if store != nil {
		result.Store = store
	}
	gob.Register(uuid.Nil)
	result.Lifetime = DefaultSessionLifetime
	result.Cookie.Secure = true
	result.Cookie.HttpOnly = true
	return result
}

// Returns a middleware to be attached to `api`.
func NewSessionMiddleware(api huma.API, manager *scs.SessionManager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		var token string
		cookie, err := huma.ReadCookie(ctx, manager.Cookie.Name)
		if err == nil {
			token = cookie.Value
		}

		newCtx, err := manager.Load(ctx.Context(), token)
		if err != nil {
			log.Err(err).Msg("internal error loading from session store")
			_ = huma.WriteErr(api, ctx, 500, "")
			return
		}

		ctx = huma.WithContext(ctx, newCtx)

		if isCookieAuthorizationRequired(ctx.Operation()) {
			if !manager.Exists(ctx.Context(), SessionKeyAuthID) {
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "")
				return
			}
		}

		next(ctx)
	}
}

// Commit changes to the session and output new cookies
func CommitSession(ctx context.Context, manager *scs.SessionManager) (SessionHeaderOutput, error) {
	var result SessionHeaderOutput
	switch manager.Status(ctx) {
	case scs.Modified:
		token, expiry, err := manager.Commit(ctx)
		if err != nil {
			return result, err
		}
		result.SetCookie = append(result.SetCookie, newSessionCookie(ctx, manager, token, expiry).String())
	case scs.Destroyed:
		result.SetCookie = append(result.SetCookie, newSessionCookie(ctx, manager, "", time.Time{}).String())
	case scs.Unmodified:
		// nothing to do
	}

	if len(result.SetCookie) > 0 {
		result.CacheControl = append(result.CacheControl, `no-cache="Set-Cookie"`)
	}
	return result, nil
}

// Creates a new session cookie
func newSessionCookie(ctx context.Context, manager *scs.SessionManager, token string, expiry time.Time) *http.Cookie {
	result := http.Cookie{
		Name:     manager.Cookie.Name,
		Value:    token,
		Path:     manager.Cookie.Path,
		Domain:   manager.Cookie.Domain,
		Secure:   manager.Cookie.Secure,
		HttpOnly: manager.Cookie.HttpOnly,
		SameSite: manager.Cookie.SameSite,
	}

	if expiry.IsZero() {
		result.Expires = time.Unix(1, 0)
		result.MaxAge = -1
	} else if manager.Cookie.Persist || manager.GetBool(ctx, SessionKeyPersist) {
		result.Expires = time.Unix(expiry.Unix()+1, 0)
		result.MaxAge = int(time.Until(expiry).Seconds() + 1)
	}

	return &result
}

func isCookieAuthorizationRequired(op *huma.Operation) bool {
	if op == nil {
		return false
	}

	for _, scheme := range op.Security {
		if _, ok := scheme[CookieSecuritySchemeName]; ok {
			return true
		}
	}
	return false
}
