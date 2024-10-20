package listing

import (
	"context"
	"errors"

	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/models"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/listing"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/parkingspot"
	"github.com/ParkWithEase/parkeasy/backend/internal/pkg/repositories/timeunit"
)

type Service struct {
	spotRepo     parkingspot.Repository
	listingRepo  listing.Repository
	timeunitRepo timeunit.Repository
}

func New(spotRepo parkingspot.Repository, listingRepo listing.Repository, timeunitRepo timeunit.Repository) *Service {
	return &Service{
		spotRepo:     spotRepo,
		listingRepo:  listingRepo,
		timeunitRepo: timeunitRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, list *models.ListingCreationInput) (int64, models.Listing, error) {
	//----Check if the parking spot belongs to the user----
	// Get the parking spot ID
	result, err := s.spotRepo.GetByUUID(ctx, list.ID)
	if err != nil {
		if errors.Is(err, parkingspot.ErrNotFound) {
			err = models.ErrParkingSpotNotFound
		}
		return -1, models.Listing{}, err
	}

	//Check for create privilages
	if result.OwnerID != userID {
		return -1, models.Listing{}, models.ErrParkingSpotNotFound
	}
	//-------------------------------------------------------

	//Create listing
	ID, res, err := s.listingRepo.Create(ctx, result.InternalID, list)
	if err != nil {
		if errors.Is(err, listing.ErrDuplicatedListing) {
			err = models.ErrListingDuplicate
		}
		return -1, models.Listing{}, err
	}

	//Add corresponding timeslots
	timeunits, err := s.timeunitRepo.Create(ctx, list.Availability)
	if err != nil {
		if errors.Is(err, timeunit.ErrDuplicatedTimeUnit) {
			err = models.ErrListingDuplicate
		}
		return -1, models.Listing{}, err
	}

	entry := models.Listing{
		Spot:         result.ParkingSpot,
		Availability: timeunits.TimeSlots,
		ID:           res.ID,
		PricePerHour: res.PricePerHour,
	}

	return ID, entry, nil
}
