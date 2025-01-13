package response

import "time"

type Restaurant struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Location  string    `bson:"location"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}


type MenuItem struct {
	ID           string    `bson:"_id,omitempty"`
	RestaurantID string    `bson:"restaurant_id"`
	Name         string    `bson:"name"`
	Description  string    `bson:"description"`
	Price        float64   `bson:"price"`
	Availability bool      `bson:"availability"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
