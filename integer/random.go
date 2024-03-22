package integer

import "math/rand"

func RandomBetween(min, max int) int {
	return min + rand.Intn(max-min)
}
