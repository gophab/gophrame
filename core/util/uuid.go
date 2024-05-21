package util

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	result := uuid.NewString()
	return strings.ReplaceAll(result, "-", "")
}
