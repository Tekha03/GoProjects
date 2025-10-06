package main

import "bank_server/router"

func main() {
	r := router.SetupRouter()
	r.Run(":8080")
}
