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

func TestCreateMovie(t *testing.T) {

	var jsonStr = []byte(` {"name": "ashishashish" , "budget":456789 , "director":"karmarkar"}`)

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

	if beforeCount == afterCount {
		t.Errorf("handler returned %v, %v", beforeCount, afterCount)
	}
}

func TestCreateMovieWithErr(t *testing.T) {
	var jsonStr = []byte(` {"name": "type mismatch" , "budget":456789 , "director":"karmarkar"}`)

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	request.Header.Set("content-type", "application/javascript")

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

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
}

func TestCreateMovieWithDataErr(t *testing.T) {

	var jsonStr = []byte(`{"name":""}`)

	session, _ := mgo.Dial("localhost:27017") //establish connection
	c := session.DB("ashish").C("plaza")      //create new database and collection
	beforeCount, _ := c.Find(bson.M{}).Count()

	request, err := http.NewRequest("POST", "/movie", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	UserRouter().ServeHTTP(response, request)
	afterCount, _ := c.Find(bson.M{}).Count()

	fmt.Println(beforeCount, afterCount)

	if beforeCount != afterCount {
		t.Errorf("expecting error not getting it %v, %v", beforeCount, afterCount)
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

func TestGetMovie(t *testing.T) {
	request, err := http.NewRequest("GET", "/movie/ssss", nil)
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
	request, err := http.NewRequest("DELETE", "/movie/ashishashish", nil)
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

	request, err := http.NewRequest("PUT", "/movie/Titanic", bytes.NewBuffer(updateStr))
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

	request, err := http.NewRequest("PUT", "/movie/Tfitanic", bytes.NewBuffer(updateStr))
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

	var updateStr = []byte(` {"budget":456789 , "director":"Bruice wills"}`)

	request, err := http.NewRequest("PATCH", "/movie/Titanic", bytes.NewBuffer(updateStr))
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

	request, err := http.NewRequest("PATCH", "/movie/Tfitanic", bytes.NewBuffer(updateStr))
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
