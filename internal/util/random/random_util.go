package random_util

import "github.com/google/uuid"

func RandomCorrelationId() string {
	return uuid.NewString()
}
