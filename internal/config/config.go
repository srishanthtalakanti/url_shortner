package config

import (
	"fmt"
	"os"
)

var required = []string{
	"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	"JWT_SECRET_KEY",
}

func Validate() error {
	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", missing)
	}
	return nil
}
