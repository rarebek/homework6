package codegen

import (
	"fmt"
	"math/rand"
)

func GenerateCode() string {
	code := rand.Intn(1000000)

	temp := fmt.Sprintf("%06d", code)
	return temp
}
