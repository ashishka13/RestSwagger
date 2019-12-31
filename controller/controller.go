package controller

import (
	"RestSwagger/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var movies []model.Movie

//InsDberror is a variable used for checking database error while inserting data
var InsDberror = false

//AllDbErr is a variable used for checking database error while fetching all records
var AllDbErr = false

//CreateMovie create new record
func CreateMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation POST /movie CreateMovie
	// Add A mov to the database
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	//   - name: movie
	//     in: body
	//     required: true
	//     description: The movie to be created.
	//     schema:
	//       $ref: '#/definitions/Movie'
	// responses:
	//   '200':
	//     description: Movie succesfully created.
	if request.Header.Get("content-type") != "application/json" {
		fmt.Printf("request type not matching")
		return
	}
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")      //create new database and collection

	var movie model.Movie
	decodeErr := json.NewDecoder(request.Body).Decode(&movie)

	if decodeErr != nil {
		fmt.Print("JSON Decode problem create")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if movie.Name == "" || movie.Budget == 0 || movie.Director == "" || request.ContentLength == 0 {
		fmt.Print("cannot send empty data")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	errIns := c.Insert(movie) //actual insert query
	if errIns != nil || InsDberror {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(movie)
}

//GetMovieByID select one movie by user id
func GetMovieByID(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /movie/{uid} GetMovieByID
	// Returns the movie from the database which user has requested
	// ---
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	// - name: uid
	//   in: path
	//   required: true
	// responses:
	//   '200':
	//     description: movie response
	//     schema:
	//       $ref: "#/definitions/Movie"
	//   default:
	//    description: unexpected error
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")
	params := mux.Vars(request)

	query := c.Find(bson.M{"uid": params["uid"]})

	result := model.Movie{}

	err := query.One(&result)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(result)

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

//GetMovies select all movies
func GetMovies(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /movies GetMovies
	// Returns all the movies from the database
	// ---
	// produces:
	// - application/json
	// - application/xml
	// responses:
	//   '200':
	//     description: movie response
	//     schema:
	//       $ref: "#/definitions/Movie"
	//   default:
	//    description: unexpected error
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")

	query := c.Find(bson.M{})
	var results []model.Movie
	err := query.All(&results)
	if err != nil || AllDbErr {
		fmt.Print("Database query problem")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	arr := []model.Movie{}
	for _, value := range results {
		arr = append(arr, value)
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(arr)
}

//UpdateMovie update one movie by name
func UpdateMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation PUT /movie/{uid} UpdateMovie
	// Updates a movie from the database
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	// - name: uid
	//   in: path
	//   required: true
	//   type: string
	//   description: The movie to be updated.
	// - name: movie
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/Movie'
	// responses:
	//   '201':
	//     description: Movie succesfully created.
	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")

	params := mux.Vars(request)
	var movie model.Movie

	decodeErr := json.NewDecoder(request.Body).Decode(&movie) //this returns error
	if decodeErr != nil {
		fmt.Print("JSON decode problem Update put")
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	what := bson.M{"uid": params["uid"]}
	change := bson.M{
		"$set": bson.M{
			"budget":   movie.Budget,
			"director": movie.Director},
	}

	err := c.Update(what, change)
	if err != nil {
		response.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(movie)
}

//UpdateMoviePatch update one movie with patch
func UpdateMoviePatch(response http.ResponseWriter, request *http.Request) {
	// swagger:operation PATCH /movie/{uid} UpdateMovie
	// Updates a movie from the database with patch
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	// - name: uid
	//   in: path
	//   required: true
	//   type: string
	//   description: The movie to be updated with patch.
	// - name: movie
	//   in: body
	//   required: true
	//   schema:
	//     $ref: '#/definitions/Movie'
	// responses:
	//   '201':
	//     description: Movie succesfully created.
	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")

	params := mux.Vars(request)
	var movie model.Movie

	decodeErr := json.NewDecoder(request.Body).Decode(&movie) //this returns error
	if decodeErr != nil {
		fmt.Print("JSON decode problem Update Patch")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	what := bson.M{"uid": params["uid"]}
	change := bson.M{
		"$set": bson.M{
			"name":     movie.Name,
			"budget":   movie.Budget,
			"director": movie.Director},
	}

	err := c.Update(what, change)
	if err != nil {
		response.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(movie)
}

//DeleteMovie delete one movie by name
func DeleteMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation DELETE /movie/{uid} DeleteMovie
	// Delete a movie from the database
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/json
	// - application/xml
	// parameters:
	// - name: uid
	//   in: path
	//   required: true
	//   type: string
	//   description: The movie to be deleted.
	//   required: true
	// responses:
	//   '200':
	//     description: Movie succesfully deleted.
	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer
	c := session.DB("ashish").C("plaza")

	params := mux.Vars(request)

	what := bson.M{"uid": params["uid"]}

	err := c.Remove(what)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		fmt.Println("delete error", err)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func generateSwagger() http.Handler {
	cmd := exec.Command("swagger", "generate", "spec", "-o", "swaggerui/swagger.json", "--scan-models")
	//swagger generate spec -o swaggerui/swagger.json --scan-models
	cmd.Run()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// }
	swgr := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	return swgr
	// return nil
}
func swaggerHome(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "http://localhost:12345/swaggerui/", http.StatusFound)
}

//UserRouter all the router related functionalities
func UserRouter() *mux.Router {
	generateSwagger()
	router := mux.NewRouter()
	router.HandleFunc("/", swaggerHome)
	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movie", CreateMovie).Methods("POST")
	router.HandleFunc("/movie/{uid}", GetMovieByID).Methods("GET")
	router.HandleFunc("/movie/{uid}", UpdateMovie).Methods("PUT")
	router.HandleFunc("/movie/{uid}", UpdateMoviePatch).Methods("PATCH")
	router.HandleFunc("/movie/{uid}", DeleteMovie).Methods("DELETE")

	router.PathPrefix("/swaggerui/").Handler(generateSwagger())

	return router
}
