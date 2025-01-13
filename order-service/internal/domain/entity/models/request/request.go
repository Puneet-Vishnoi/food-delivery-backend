package request

type OrderItem struct {
	MenuItemID string  `bson:"menu_item_id"`
	Name       string  `bson:"name"`
	Quantity   int     `bson:"quantity"`
	Price      float64 `bson:"price"`
}

type Order struct {
	ID           string      `bson:"_id,omitempty"`
	UserID       string      `bson:"user_id"`
	RestaurantID string      `bson:"restaurant_id"`
	Items        []OrderItem `bson:"items"`
	TotalAmount  float64     `bson:"total_amount"`
	Status       string      `bson:"status"`
	CreatedAt    int64       `bson:"created_at"`
	UpdatedAt    int64       `bson:"updated_at"`
}
