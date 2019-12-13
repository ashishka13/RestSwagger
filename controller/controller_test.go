package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func changeVal(variable bool) {
	variable = false
}

//check flags of test command
func TestCreateMovie(t *testing.T) {

	var jsonStr = []byte(` {"uid":"3", "name": "ssss" , "budget":456789 , "director":"drax"}`)

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("content-type", "application/json")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)

	if beforeCount == afterCount {
		t.Errorf("handler returned %v, %v", beforeCount, afterCount)
	}
}

func TestCreateMovieDbErr(t *testing.T) {

	InsDberror = true		//make it false on purpose, this is referred in controller again
	var jsonStr = []byte(` {"uid":"t3", "name": "ssss" , "budget":456789 , "director":"drax"}`)

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("content-type", "application/json")
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)
	cod := response.Code

	if cod != 500 {
		t.Errorf("Expecting database error but not getting error \n entries before%v, entries after %v", beforeCount, afterCount)
	}
	defer changeVal(AllDbErr)
}

func TestCreateMovieTypeErr(t *testing.T) {
	var jsonStr = []byte(` {"uid":"r4","name": "type mismatch" , "budget":456789 , "director":"karmarkar"}`)

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	request.Header.Set("content-type", "application/javascript")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, err := c.Find(bson.M{}).Count()

	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)

	if beforeCount != afterCount {
		t.Errorf("expecting error not getting error  %v, %v", beforeCount, afterCount)
	}
	// defer Dberror = false
}

func TestCreateMovieWithJsonErr(t *testing.T) {
	var jsonStr = []byte(` {q2w34er5t6y7u89i90o0-p}`)

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	request.Header.Set("content-type", "application/json")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, err := c.Find(bson.M{}).Count()

	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)
	cod := response.Code
	if cod != 500 {
		t.Errorf("expecting error not getting error  %v, %v", beforeCount, afterCount)
	}

}

func TestCreateMovieWithNoData1(t *testing.T) {

	var jsonStr = []byte(`{}`)

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	request.Header.Set("content-type", "application/json")

	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)

	cod := response.Code
	if cod != 400 {
		t.Errorf("expecting error not getting error  %v, %v", beforeCount, afterCount)
	}
}

func TestCreateMovieWithNoData2(t *testing.T) {
	var jsonStr = []byte(` {"id":"", "name": "" , "budget":0 , "director":""}`)

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	request.Header.Set("content-type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)

	cod := response.Code

	if cod != 400 {
		t.Errorf("expecting error not getting error  %v, %v", beforeCount, afterCount)
	}
}

func TestGetMovies(t *testing.T) {
	request, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	assert.NotEqual(t, 404, response.Code)
}

func TestGetMoviesWithErr(t *testing.T) {
	AllDbErr = true
	request, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	assert.NotEqual(t, 404, response.Code)

	defer changeVal(AllDbErr)
}

func TestGetMovie(t *testing.T) {
	request, err := http.NewRequest("GET", "/movie/4", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code == 404 {
		t.Error("no record found")
	}
}

func TestGetMovieWithErr(t *testing.T) {
	request, err := http.NewRequest("GET", "/movie/sssshgf", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 404 {
		t.Error("expecting error not getting error")
	}
}

func TestDeleteMovie(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/movie/3", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 200 {
		t.Error("operation failed")
	}
}

func TestDeleteMovieWithErr(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/movie/ashifshashish", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 404 {
		t.Error("Expecting error but not getting")
	}
}

func TestUpdateMovie(t *testing.T) {

	var updateStr = []byte(` {"budget":456789 , "director":"Bruice wills"}`)

	request, err := http.NewRequest("PUT", "/movie/t1", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 200 {
		t.Error("operation failed")
	}
}

func TestUpdateMovieWithError(t *testing.T) {

	var updateStr = []byte(` {"budget":456789 , "director":"Bruice wills"}`)

	request, err := http.NewRequest("PUT", "/movie/w2w2w", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 406 {
		t.Error("Expecting error but not getting")
	}
}

func TestUpdateMovieWithJsonError(t *testing.T) {

	var updateStr = []byte(` {"q2w3 # @ !4ey7u  ^ * $ )89i0"}`)

	request, err := http.NewRequest("PUT", "/movie/r5r5r", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 406 {
		t.Error("Expecting error but not getting")
	}
}

func TestUpdatePatch(t *testing.T) {

	var updateStr = []byte(` {"name":"same", "budget":456789 , "director":"Bruice wills"}`)

	request, err := http.NewRequest("PATCH", "/movie/t3", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 200 {
		t.Error("operation failed")
	}
}

func TestUpdatePatchWithError(t *testing.T) {

	var updateStr = []byte(` {"budget":456789 , "director":"Bruice wills"}`)

	request, err := http.NewRequest("PATCH", "/movie/e54e54e", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 406 {
		t.Error("Expecting error but not getting")
	}
}

func TestUpdatePatchWithJsonError(t *testing.T) {

	var updateStr = []byte(` {"q2w3 # @ !4ey7u  ^ * $ )89i0"}`)

	request, err := http.NewRequest("PATCH", "/movie/t2", bytes.NewBuffer(updateStr))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	fmt.Println(response.Code)

	if response.Code != 500 {
		t.Error("Expecting error but not getting")
	}
}
