package utils

import "github.com/pborman/uuid"

// GenerateUUID generates a new uuid for a minion
func GenerateUUID(name string) uuid.UUID {
	u := uuid.NewSHA1(uuid.NameSpace_DNS, []byte(name))

	return u
}
