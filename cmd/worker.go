package main

import (
	"billingApp/workers"
	_ "github.com/lib/pq"
)

//Run worker
func main() {
	workers.Server()
}
