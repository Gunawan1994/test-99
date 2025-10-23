package entity

type User struct {
	ID        int64  `json:"id" bson:"user_id"`
	Name      string `json:"name" bson:"name"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}
