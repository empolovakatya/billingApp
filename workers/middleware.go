package workers

import "github.com/sirupsen/logrus"

func FailOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}
