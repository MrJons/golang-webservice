package main

import (
	"fmt"
	"github.com/mrjons/webservice/controllers"
	"net/http"
)

const port = 3000

// Start server. Currently runs on static port
func main() {
	fmt.Println("Starting server on port", port)

	controllers.RegisterControllers()
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
