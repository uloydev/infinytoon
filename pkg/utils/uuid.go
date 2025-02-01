package utils

import "github.com/google/uuid"

func GetUUID() string {
	return uuid.Must(uuid.NewV7()).String()
}

func ValidateUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
