package main

import (
	"net/http"

	"./controllers"
)

func main() {
	controllers.RegisterControllers()

	http.ListenAndServe(":3000", nil)
}

// func startWebServer(port int) (int, error) {
// 	fmt.Println("Starting server...")

// 	fmt.Println("Server started on port", port)
// 	return port, nil
// }
