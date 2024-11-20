package routes

// import (
// 	"context"
// 	"errors"
// 	"net/http"

// 	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
// 	"github.com/danielgtaylor/huma/v2"
// )

// // Service provider for `DemoRoute`
// type DemoServicer interface {
// 	// Get the data stored.
// 	Get(ctx context.Context) (string, error)
// }

// type DemoRoute struct {
// 	service       DemoServicer
// 	sessionGetter SessionDataGetter
// }

// var DemoTag = huma.Tag{
// 	Name:        "Demo",
// 	Description: "Operations for demo.",
// }

// type DemoOutput struct {
// 	Body models.Demo
// }

// // Returns a new `ParkingSpotRoute`
// func NewDemoRoute(
// 	service DemoServicer,
// 	sessionGetter SessionDataGetter,
// ) *DemoRoute {
// 	return &DemoRoute{
// 		service:       service,
// 		sessionGetter: sessionGetter,
// 	}
// }

// func (r *DemoRoute) RegisterDemoTag(api huma.API) {
// 	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &DemoTag)
// }

// func (r *DemoRoute) RegisterDemoRoutes(api huma.API) {
// 	huma.Register(api, *withUserID(&huma.Operation{
// 		OperationID: "get-demo-data",
// 		Method:      http.MethodGet,
// 		Path:        "/demo",
// 		Summary:     "Get demo data",
// 		Tags:        []string{DemoTag.Name},
// 	}), func(ctx context.Context, input *struct{}) (*DemoOutput, error) {
// 		// _ = r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)

// 		resultString, err := r.service.Get(ctx)
// 		if err != nil {
// 			if errors.Is(err, models.ErrNoData) {
// 				err = NewHumaError(ctx, http.StatusNoContent, err)
// 			} else {
// 				err = NewHumaError(ctx, http.StatusInternalServerError, err)
// 			}

// 			return nil, err
// 		}

// 		result := DemoOutput{Body: models.Demo(resultString)}

// 		return &result, nil
// 	})
// }
