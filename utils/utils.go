package utils

import (
	"math/rand"
	"time"
)

func RandomRange(min, max float32) float32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + (max-min)*r.Float32()
}
