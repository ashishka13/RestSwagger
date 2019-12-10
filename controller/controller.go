package controller

import (
	"RestSwagger/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var movies []model.Movie

//CreateMovie create new record
func CreateMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation POST /movie CreateMovie
	// Adds new movie to the database
	// could be any movie
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/xml
	// - application/json
	// parameters:
	//   - name: movie
	//     in: body
	//     required: true
	//     description: The movie to create.
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
	_ = json.NewDecoder(request.Body).Decode(&movie)

	if movie.Name == "" || request.ContentLength == 0 {
		fmt.Print("cannot send empty data")
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	errIns := c.Insert(movie) //actual insert query
	if errIns != nil {
		response.WriteHeader(http.StatusNotAcceptable)
	} else {
		response.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(response).Encode(movie)
}

//GetMovie select one movie
func GetMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /movie/{name} GetMovie
	// Returns the movie from the database which user has requested
	// Could be any movie
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// parameters:
	// - name: name
	//   in: path
	//   required: true
	//   type: string
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
	result := model.Movie{}

	query := c.Find(bson.M{"name": params["name"]})
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
	// Returns the movies from the database
	// List of all movies
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// responses:
	//   '200':
	//     description: movie response
	//     schema:
	//       $ref: "#/definitions/Movie"
	//   default:
	//    description: unexpected error
	fmt.Printf("%v", request.Header.Get("content-type"))
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	session, _ := mgo.Dial("localhost:27017") //establish connection
	defer session.Close()                     //close it in defer

	c := session.DB("ashish").C("plaza")
	var results []model.Movie

	query := c.Find(bson.M{})
	query.All(&results)

	arr := []model.Movie{}
	for _, value := range results {
		// fmt.Println(value,err)
		arr = append(arr, value)
	}
	json.NewEncoder(response).Encode(arr)
}

//UpdateMovie update one movie
func UpdateMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation PUT /movie/{name} UpdateMovie
	// Updates a movie from the database
	// could be any movie
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/xml
	// - application/json
	// parameters:
	// - name: name
	//   in: path
	//   required: true
	//   type: string
	//   description: The movie to update.
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
	var movie model.Movie
	params := mux.Vars(request)

	_ = json.NewDecoder(request.Body).Decode(&movie) //this returns error

	what := bson.M{"name": params["name"]}
	change := bson.M{
		"$set": bson.M{
			"budget":   movie.Budget,
			"director": movie.Director},
	}

	err := c.Update(what, change)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else if err != nil {
		response.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err)
		return
	}
}

//UpdateMoviePatch update one movie only with patch
func UpdateMoviePatch(response http.ResponseWriter, request *http.Request) {
	// swagger:operation PATCH /movie/{name} UpdateMovie
	// Updates a movie from the database with patch
	// could be any movie
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/xml
	// - application/json
	// parameters:
	// - name: name
	//   in: path
	//   required: true
	//   type: string
	//   description: The movie to update patch.
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
	var movie model.Movie
	params := mux.Vars(request)

	_ = json.NewDecoder(request.Body).Decode(&movie) //this returns error

	what := bson.M{"name": params["name"]}
	change := bson.M{
		"$set": bson.M{
			"budget":   movie.Budget,
			"director": movie.Director},
	}

	err := c.Update(what, change)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else if err != nil {
		response.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err)
		return
	}
}

//DeleteMovie delete one movie
func DeleteMovie(response http.ResponseWriter, request *http.Request) {
	// swagger:operation DELETE /movie/{name} DeleteMovie
	// Delete a movie from the database
	// could be any movie
	// ---
	// consumes:
	// - application/json
	// - application/xml
	// produces:
	// - application/xml
	// - application/json
	// parameters:
	// - name: name
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

	what := bson.M{"name": params["name"]}
	err := c.Remove(what)

	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else if err != nil {
		response.WriteHeader(http.StatusNotFound)
		fmt.Println("ashish", err)
		return
	}
}

//UserRouter all the router related functionalities
func UserRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movie", CreateMovie).Methods("POST")
	router.HandleFunc("/movie/{name}", GetMovie).Methods("GET")
	router.HandleFunc("/movie/{name}", UpdateMovie).Methods("PUT")
	router.HandleFunc("/movie/{name}", UpdateMoviePatch).Methods("PATCH")
	router.HandleFunc("/movie/{name}", DeleteMovie).Methods("DELETE")
	return router
}
