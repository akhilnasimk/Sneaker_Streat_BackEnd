package helpers

import (
	"log"

	"github.com/google/uuid"
)

func StringToUUID(val string) uuid.UUID {
	id, err := uuid.Parse(val)
	if err != nil {
		log.Printf("Invalid UUID string: %v", err)
		return uuid.Nil // returns a zero UUID
	}
	return id
}