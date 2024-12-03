package routes

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Service provider for `BookingRoute`
type BookingServicer interface {
	// Creates a new booking attached to `userID`.
	//
	// Returns the booking internal ID and the model.
	Create(ctx context.Context, userID int64, bookingDetails *models.BookingCreationInput) (int64, models.BookingWithTimes, error)
	// Get at most `count` bookings associated with the given `userID` that satisfies the given filter conditions.
	//
	// If there are more entries following the result, a non-empty cursor will be returned
	// which can be passed to the next invocation to get the next entries.
	GetManyForSeller(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error)
	// Get at most `count` bookings associated with the given `userID` that satisfies the given filter conditions.
	//
	// If there are more entries following the result, a non-empty cursor will be returned
	// which can be passed to the next invocation to get the next entries.
	GetManyForBuyer(ctx context.Context, userID int64, count int, after models.Cursor, filter models.BookingFilter) ([]models.Booking, models.Cursor, error)
	// Get the booking with `bookingID` if `userID` has enough permission to view the resource.
	GetByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) (models.BookingWithTimes, error)
	// Get booked times with `bookingID if `userID` has enough permission to view the resource.
	GetBookedTimesByUUID(ctx context.Context, userID int64, bookingID uuid.UUID) ([]models.TimeUnit, error)
}

// BookingRoute represents booking-related API routes
type BookingRoute struct {
	service       BookingServicer
	sessionGetter SessionDataGetter
}

type bookingListOutput struct {
	Link []string         `header:"Link" doc:"Contains details on getting the next page of resources" example:"<https://example.com/bookings?after=gQL>; rel=\"next\""`
	Body []models.Booking `nullable:"false"`
}

type bookingWithTimesOutput struct {
	Body models.BookingWithTimes
}

type bookedTimes struct {
	Body []models.TimeUnit
}

var BookingTag = huma.Tag{
	Name:        "Booking",
	Description: "Operations for handling bookings.",
}

// Returns a new `BookingRoute`
func NewBookingRoute(
	service BookingServicer,
	sessionGetter SessionDataGetter,
) *BookingRoute {
	return &BookingRoute{
		service:       service,
		sessionGetter: sessionGetter,
	}
}

func (r *BookingRoute) RegisterBookingTag(api huma.API) {
	api.OpenAPI().Tags = append(api.OpenAPI().Tags, &BookingTag)
}

// Registers `/booking` routes
func (r *BookingRoute) RegisterBookingRoutes(api huma.API) {
	apiPrefix := getAPIPrefix(api.OpenAPI())

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID:   "create-booking",
		Method:        http.MethodPost,
		Path:          "/book",
		Summary:       "Create a new booking",
		Tags:          []string{BookingTag.Name},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity},
	}), func(ctx context.Context, input *struct {
		Body models.BookingCreationInput
	},
	) (*bookingWithTimesOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		_, result, err := r.service.Create(ctx, userID, &input.Body)
		if err != nil {
			var detail error
			switch {
			case errors.Is(err, models.ErrEmptyBookingTimes):
				detail = &huma.ErrorDetail{
					Location: "body.book_times",
					Value:    input.Body.BookedTimes,
				}
			case errors.Is(err, models.ErrParkingSpotNotFound):
				detail = &huma.ErrorDetail{
					Location: "body.parkingspot_id",
					Value:    input.Body.ParkingSpotID,
				}
			case errors.Is(err, models.ErrCarNotFound):
				detail = &huma.ErrorDetail{
					Location: "body.car_id",
					Value:    input.Body.CarID,
				}
			case errors.Is(err, models.ErrCarNotOwned):
				detail = &huma.ErrorDetail{
					Location: "body.car_id",
					Value:    input.Body.CarID,
				}
			case errors.Is(err, models.ErrDuplicateBooking):
				detail = &huma.ErrorDetail{
					Location: "body.book_times",
					Value:    input.Body.BookedTimes,
				}
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err, detail)
		}
		return &bookingWithTimesOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "list-bookings",
		Method:      http.MethodGet,
		Path:        "/users/bookings",
		Summary:     "Get bookings associated to the current user (buyer)",
		Tags:        []string{BookingTag.Name},
	}), func(ctx context.Context, input *struct {
		Filter models.BookingFilter
		After  models.Cursor `query:"after" doc:"Token used for requesting the next page of resources"`
		Count  int           `query:"count" minimum:"1" default:"50" doc:"The maximum number of bookings that appear per page."`
	},
	) (*bookingListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		bookings, nextCursor, err := r.service.GetManyForBuyer(ctx, userID, input.Count, input.After, input.Filter)
		if err != nil {
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := bookingListOutput{Body: bookings}
		if nextCursor != "" {
			nextURL := apiPrefix.JoinPath("/users/bookings")
			nextURL.RawQuery = url.Values{
				"count": []string{strconv.Itoa(input.Count)},
				"after": []string{string(nextCursor)},
			}.Encode()
			result.Link = append(result.Link, "<"+nextURL.String()+`>; rel="next"`)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "list-leasings",
		Method:      http.MethodGet,
		Path:        "/users/leasings",
		Summary:     "Get leasings associated to the current user (seller)",
		Tags:        []string{BookingTag.Name},
	}), func(ctx context.Context, input *struct {
		Filter models.BookingFilter
		After  models.Cursor `query:"after" doc:"Token used for requesting the next page of resources"`
		Count  int           `query:"count" minimum:"1" default:"50" doc:"The maximum number of bookings that appear per page."`
	},
	) (*bookingListOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		bookings, nextCursor, err := r.service.GetManyForSeller(ctx, userID, input.Count, input.After, input.Filter)
		if err != nil {
			if errors.Is(err, models.ErrSpotNotOwned) {
				detail := &huma.ErrorDetail{
					Location: "path.parkingspot_id",
					Value:    input.Filter.ParkingSpotID,
				}
				return nil, NewHumaError(ctx, http.StatusForbidden, err, detail)
			}
			return nil, NewHumaError(ctx, http.StatusUnprocessableEntity, err)
		}

		result := bookingListOutput{Body: bookings}
		if nextCursor != "" {
			nextURL := apiPrefix.JoinPath("/users/leasings")
			nextURL.RawQuery = url.Values{
				"count": []string{strconv.Itoa(input.Count)},
				"after": []string{string(nextCursor)},
			}.Encode()
			result.Link = append(result.Link, "<"+nextURL.String()+`>; rel="next"`)
		}
		return &result, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-booking",
		Method:      http.MethodGet,
		Path:        "/bookings/{id}",
		Summary:     "Get information about a booking",
		Tags:        []string{BookingTag.Name},
		Errors:      []int{http.StatusNotFound},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*bookingWithTimesOutput, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		result, err := r.service.GetByUUID(ctx, userID, input.ID)
		if err != nil {
			var detail error
			status := http.StatusUnprocessableEntity

			switch {
			case errors.Is(err, models.ErrBookingNotFound):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			case errors.Is(err, models.ErrInvalidRequest):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				status = http.StatusForbidden
			}
			return nil, NewHumaError(ctx, status, err, detail)
		}
		return &bookingWithTimesOutput{Body: result}, nil
	})

	huma.Register(api, *withUserID(&huma.Operation{
		OperationID: "get-booked-time-slots-of-a-booking",
		Method:      http.MethodGet,
		Path:        "/bookings/{id}/availability",
		Summary:     "Get booked time slots for the booking",
		Tags:        []string{BookingTag.Name},
	}), func(ctx context.Context, input *struct {
		ID uuid.UUID `path:"id"`
	},
	) (*bookedTimes, error) {
		userID := r.sessionGetter.Get(ctx, SessionKeyUserID).(int64)
		spots, err := r.service.GetBookedTimesByUUID(ctx, userID, input.ID)
		if err != nil {
			var detail error
			status := http.StatusUnprocessableEntity

			switch {
			case errors.Is(err, models.ErrBookingNotFound):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
			case errors.Is(err, models.ErrInvalidRequest):
				detail = &huma.ErrorDetail{
					Location: "path.id",
					Value:    input.ID,
				}
				status = http.StatusForbidden
			}
			return nil, NewHumaError(ctx, status, err, detail)
		}

		result := bookedTimes{Body: spots}

		return &result, nil
	})
}
