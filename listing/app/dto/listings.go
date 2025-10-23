package dto

type Listing struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id" validate:"required"`
	ListingType string `json:"listing_type" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
}
