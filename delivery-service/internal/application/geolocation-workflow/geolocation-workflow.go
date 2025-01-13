package geolocationworkflow

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/adapter/repository"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"
)

type GeolocationWorkflowService struct{
}

func NewGeolocationWorkflowServiceProvider(geolocationWorkflowRepo repository.GeolocationWorkflowRepo)(*GeolocationWorkflowService){
	return &GeolocationWorkflowService{
		// geolocationWorkflowRepo : geolocationWorkflowRepo,
	}
}

//  a geo-fenced area 
func (serv *GeolocationWorkflowService) IsWithinGeoFence(currentLocation request.Location, geoFenceBounds []request.Location) bool {
	//use a rectangular bounding box for geo-fencing.
	latMin, latMax := geoFenceBounds[0].Lat, geoFenceBounds[1].Lat
	lngMin, lngMax := geoFenceBounds[0].Lng, geoFenceBounds[1].Lng

	// Check if current location is within bounding box
	if currentLocation.Lat >= latMin && currentLocation.Lat <= latMax &&
		currentLocation.Lng >= lngMin && currentLocation.Lng <= lngMax {
		return true
	}

	return false
}

// Define structure for Google Maps Distance API response
type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
		} `json:"elements"`
	} `json:"rows"`
}

const googleAPIKey = "YOUR_GOOGLE_MAPS_API_KEY" // Replace with your API key

// Function to calculate distance from pickup to personnel location
func (serv *GeolocationWorkflowService)GetDistanceFromPickupToPersonnel(pickupLat, pickupLng, personnelLat, personnelLng float32) (string, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%f,%f&destinations=%f,%f&key=%s",
		pickupLat, pickupLng, personnelLat, personnelLng, googleAPIKey)

	// Make the request to Google Maps API
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	var result DistanceMatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	// Get the distance between pickup and personnel
	if len(result.Rows) > 0 && len(result.Rows[0].Elements) > 0 {
		distance := result.Rows[0].Elements[0].Distance.Text
		return distance, nil
	}

	return "", fmt.Errorf("no route found")
}

func (serv *GeolocationWorkflowService)CalculateDistance(pickupLat, pickupLng, personnelLat, personnelLng float32) (string, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%f,%f&destinations=%f,%f&key=%s",
		pickupLat, pickupLng, personnelLat, personnelLng, googleAPIKey)

	// Make the request to Google Maps API
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	var result DistanceMatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	// Get the distance between pickup and personnel
	if len(result.Rows) > 0 && len(result.Rows[0].Elements) > 0 {
		distance := result.Rows[0].Elements[0].Distance.Text
		return distance, nil
	}

	return "", fmt.Errorf("no route found")
}