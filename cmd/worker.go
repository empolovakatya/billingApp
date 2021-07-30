package main

import (
	"billingApp/workers"
	_ "github.com/lib/pq"
)

func main() {
	//workers.Work()
	workers.Server()
}
