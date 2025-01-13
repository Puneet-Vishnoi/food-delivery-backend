package deliverypersonal

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/adapter/repository"
	constant "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/ports/service"
)

type DeliveryPersonalService struct {
	locationService      service.GeolocationWorkflowService
	deliveryPersonalRepo repository.DeliveryPersonalRepo
}

func NewDeliveryPersonalServiceProvider(deliveryPersonalRepo repository.DeliveryPersonalRepo) *DeliveryPersonalService {
	return &DeliveryPersonalService{
		deliveryPersonalRepo: deliveryPersonalRepo,
	}
}

// Function to assign order to nearest personnel
func (serv *DeliveryPersonalService) AssignOrder(orderID string, pickupLocation request.Location, deliveryLocation request.Location) (*request.DeliveryPersonal, error) {
	// Fetch all available personnel
	personnelList, err := serv.deliveryPersonalRepo.GetAvailablePersonnel(constant.PersonnelCollection)
	if err != nil {
		log.Println("Error fetching personnel:", err)
		return nil, err
	}

	var nearestPersonnel *request.DeliveryPersonal
	minDistance := float32(math.MaxFloat32) // Initialize to the largest possible float32 value

	for _, personnel := range personnelList {
		// Calculate distance between pickup location and personnel's location
		distanceStr, err := serv.locationService.GetDistanceFromPickupToPersonnel(
			pickupLocation.Lat, pickupLocation.Lng,
			personnel.CurrentLocation.Lat, personnel.CurrentLocation.Lng,
		)
		if err != nil {
			log.Printf("Error calculating distance for personnel ID %s: %v", personnel.Id, err)
			continue
		}

		// Convert distance to float32
		distanceFloat64, err := strconv.ParseFloat(distanceStr, 32)
		if err != nil {
			log.Printf("Error converting distance to float32 for personnel ID %s: %v", personnel.Id, err)
			continue
		}
		distance := float32(distanceFloat64)

		// Check if this personnel is the closest one
		if nearestPersonnel == nil || distance < minDistance {
			nearestPersonnel = personnel
			minDistance = distance
		}
	}

	// Assign the order to the nearest personnel
	if nearestPersonnel != nil {
		err = serv.deliveryPersonalRepo.AssignOrderToPersonnel(orderID, nearestPersonnel.Id, &pickupLocation, &deliveryLocation, constant.AssignOrderCollection)
		if err != nil {
			log.Println("Error assigning order:", err)
			return nil, err
		}
		return nearestPersonnel, nil
	}

	return nil, fmt.Errorf("no available personnel found")
}

func (serv *DeliveryPersonalService) UpdatePersonnelLocation(personnelID string, currentLocation request.Location) error {
	err := serv.deliveryPersonalRepo.UpdatePersonnelLocation(personnelID, &currentLocation, constant.PersonnelCollection)
	if err != nil {
		return fmt.Errorf("failed to update personnel status to in_zone: %w", err)
	}
	return nil
}

func (serv *DeliveryPersonalService) UpdateDeliveryStatus(personnelID string, currentLocation request.Location) error {
	// Define geo-fence bounds
	geoFenceBounds := []request.Location{
		{Lat: 12.91, Lng: 77.60}, // example lower bound
		{Lat: 13.01, Lng: 77.70}, // example upper bound
	}

	// Check if personnel is within geo-fence
	if !serv.locationService.IsWithinGeoFence(currentLocation, geoFenceBounds) {
		// Mark as out of zone
		err := serv.deliveryPersonalRepo.UpdatePersonnelStatus(personnelID, "out_of_zone", constant.PersonnelCollection)
		if err != nil {
			return fmt.Errorf("failed to update personnel status to out_of_zone: %w", err)
		}
		return fmt.Errorf("personnel is out of geo-fence bounds")
	}

	// Mark the personnel as in-zone
	err := serv.deliveryPersonalRepo.UpdatePersonnelStatus(personnelID, "in_zone", constant.PersonnelCollection)
	if err != nil {
		return fmt.Errorf("failed to update personnel status to in_zone: %w", err)
	}

	// Fetch assigned order for the personnel
	order, err := serv.deliveryPersonalRepo.GetAssignedOrderByPersonnelId(personnelID, constant.PersonnelCollection)
	if err != nil {
		return fmt.Errorf("failed to fetch assigned order for personnel %s: %w", personnelID, err)
	}

	if order.Id == "" {
		return fmt.Errorf("no assigned order found for personnel %s", personnelID)
	}

	deliveryDistance, err := serv.locationService.CalculateDistance(
		currentLocation.Lat, currentLocation.Lng,
		order.DeliveryLocation.Lat, order.DeliveryLocation.Lng,
	)
	if err != nil {
		return fmt.Errorf("failed to calculate delivery distance: %w", err)
	}
	distanceFloat64, err := strconv.ParseFloat(deliveryDistance, 32)
	if err != nil {
		log.Printf("Error converting distance to float32 for personnel ID %s: %v", personnelID, err)
		return fmt.Errorf("failed to calculate delivery distance: %w", err)
	}
	distance := float32(distanceFloat64)

	const deliveryRadius = 0.5 // 0.5 km or 500 meters

	if distance <= deliveryRadius {
		// Mark the order as delivered
		err := serv.deliveryPersonalRepo.UpdatePersonnelStatus(personnelID, "delivered", constant.PersonnelCollection)
		if err != nil {
			return fmt.Errorf("failed to update order status to delivered: %w", err)
		}

		// Mark personnel as available
		// err = serv.deliveryPersonalRepo.UpdatePersonnelAvailability(personnelID, true, constant.PersonnelCollection)
		// if err != nil {
		// 		return fmt.Errorf("failed to mark personnel as available: %w", err)
		// }
	} else {
		// Update order status as in-progress
		err := serv.deliveryPersonalRepo.UpdatePersonnelStatus(personnelID, "in_progress", constant.PersonnelCollection)
		if err != nil {
			return fmt.Errorf("failed to update order status to in_progress: %w", err)
		}
	}



	return nil
}

func (serv *DeliveryPersonalService) GetDeliveryStatus(orderId string) (*request.DeliveryPersonal, error) {
	personnelID, err := serv.deliveryPersonalRepo.GetAssignedPersonnelOrderById(orderId, constant.AssignOrderCollection)
	if err != nil {
		return nil, err
	}

	deliverypersonal, err := serv.deliveryPersonalRepo.GetDeliveryStatus(personnelID, constant.PersonnelCollection)
	if err != nil {
		return nil, err
	}

	return deliverypersonal, nil
}
