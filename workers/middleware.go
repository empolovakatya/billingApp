package workers

import (
	"github.com/sirupsen/logrus"
	"math"
)

const presicion = 1e6

//FailOnError returns fatal errors
func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}

// FloatToInt converts float to int representation
func FloatToInt(x float64) float64 {
	return math.Round(x * presicion)
}
