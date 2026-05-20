package config

import "fmt"

var (
	ErrConfigNotFoundByKey = func(key string) error {
		return fmt.Errorf("config not found by key = %q", key)
	}
)
