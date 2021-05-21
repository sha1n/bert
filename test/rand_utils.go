package test

import (
	"fmt"
	"time"
)

func RandomLabels() []string {
	return []string{RandomString(), RandomString()}
}

func RandomString() string {
	return fmt.Sprintf("label-%d", time.Now().Nanosecond())
}
