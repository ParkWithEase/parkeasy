package routes

import (
	"context"
	"net/http"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/services/health"
	"github.com/danielgtaylor/huma/v2"
)

type HealthRoute struct {
	srv *health.Service
}

var HealthTag = huma.Tag{
	Name:        "API Health",
	Description: "Queries about API status.",
}

func NewHealthRoute(srv *health.Service) *HealthRoute {
	return &HealthRoute{
		srv: srv,
	}
}

func (*HealthRoute) RegisterHealthTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &HealthTag)
}

func (r *HealthRoute) RegisterHealthRoute(api huma.API) {
	huma.Register(api, *skipSession(&huma.Operation{
		OperationID: "check-health",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Return whether API server is ready to serve",
		Tags:        []string{HealthTag.Name},
		Errors:      []int{http.StatusInternalServerError},
	}), func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		err := r.srv.CheckHealth(ctx)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusInternalServerError, err)
		}
		return nil, nil
	})
}
