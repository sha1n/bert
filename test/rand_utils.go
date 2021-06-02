package test

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// RandomStrings returns a slice of random strings
func RandomStrings() []string {
	values := []string{}
	for i := 0; i < rand.Intn(10); i++ {
		values = append(values, RandomString())
	}

	return values
}

// RandomString returns a random UUID based string...
func RandomString() string {
	uid, _ := uuid.NewRandom()
	return fmt.Sprintf("str-%s", uid.String())
}

// RandomBool ...
func RandomBool() bool {
	return time.Now().Nanosecond()%2 == 0
}

// RandomUint ...
func RandomUint() uint {
	return uint(rand.Uint32())
}
