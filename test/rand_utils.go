package test

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func RandomLabels() []string {
	return []string{RandomString(), RandomString()}
}

func RandomString() string {
	uid, _ := uuid.NewRandom()
	return fmt.Sprintf("label-%s", uid.String())
}

func RandomBool() bool {
	return time.Now().Nanosecond()%2 == 0
}

func RandomUint() uint {
	return uint(rand.Uint32())
}
