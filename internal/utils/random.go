package utils

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

var onlyOnce sync.Once

func GetRandomInt(max int) int {
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	return rand.Intn(max)
}

func GetRandomIntDist(max, pow int) int {
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	return int(float64(max) * math.Pow(rand.Float64(), float64(pow)))
}
