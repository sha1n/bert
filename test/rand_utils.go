package test

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomLabels() []string {
	return []string{RandomString(), RandomString()}
}

func RandomString() string {
	return fmt.Sprintf("label-%d", time.Now().Nanosecond())
}

func RandomBool() bool {
	return time.Now().Nanosecond()%2 == 0
}

func RandomUint() uint {
	return uint(rand.Uint32())
}
