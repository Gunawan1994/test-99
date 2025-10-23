package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func Join(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func ToObjectID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return objID
}
