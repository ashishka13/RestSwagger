// Package main Movie API.
//
//     Schemes: http, https
//     Host: localhost:12345
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//     - application/xml
//     Produces:
//     - application/json
//     - application/xml
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package main

import (
	"RestSwagger/controller"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	fmt.Printf("\nStarting server localhost:12345...\n")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":12345", handlers.CORS(headers, methods, origins)(controller.UserRouter()))
}

// this program us run using: gslab@GS-4260:~/go/src/RestSwagger$ swagger generate spec -o ./swagger.json --scan-models && swagger serve -F=swagger swagger.json
