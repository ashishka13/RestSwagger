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
	"github.com/gorilla/handlers"
	"net/http"
)

func main() {

	fmt.Printf("\nStarting server...\n")
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	controller.UserRouter().PathPrefix("/swaggerui/").Handler(sh)

	http.ListenAndServe(":12345", handlers.CORS(headers, methods, origins)(controller.UserRouter()))
}

// this program us run using: gslab@GS-4260:~/go/src/RestSwagger$ swagger generate spec -o ./swagger.json --scan-models && swagger serve -F=swagger swagger.json
