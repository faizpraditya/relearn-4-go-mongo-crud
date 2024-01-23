package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-mongo-crud/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// using mgo.v2
// type UserController struct {
// 	session *mgo.Session
// }

// func NewUserController(s *mgo.Session) *UserController {
// 	return &UserController{s}
// }

// // struct method for UserController
// func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")

// 	if !bson.IsObjectIdHex(id) {
// 		// w=response, writeheader=tosend header (err code, status code, etc)
// 		w.WriteHeader(http.StatusNotFound)
// 	}

// 	oid := bson.ObjectIdHex(id)

// 	u := models.User{}

// 	if err := uc.session.DB("dbname").C("collectionname/user").FindId(oid).One(&u); err != nil {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	// send message to postman
// 	uj, err := json.Marshal(u)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "%s\n", uj)
// }

// // _ httprouter.Params, because we don't need to use the params (for example: id) from the postman/frontend because we created a new user
// func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	u := models.User{}

// 	json.NewDecoder(r.Body).Decode(&u)

// 	u.Id = bson.NewObjectId()

// 	uc.session.DB("go-project").C("users").Insert(u)

// 	// send message to postman/frontend
// 	uj, err := json.Marshal(u)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "%s\n", uj)
// }

// func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")

// 	if !bson.IsObjectIdHex(id) {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	oid := bson.ObjectIdHex(id)

// 	// mongo db function to delete by id
// 	if err := uc.session.DB("go-project").C("users").RemoveId(oid); err != nil {
// 		w.WriteHeader(404)
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, "Deleted user", oid, "\n")
// 	// For HTTP handlers, it's common to use fmt.Fprint with the http.ResponseWriter because it allows you to send a response to the client via the provided writer.
// 	// Using fmt.Print wouldn't be suitable in this context since it prints to the console, and it doesn't allow you to direct the output to the HTTP response.
// }

// using mongo-go-driver

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)

	u := models.User{}

	collection := uc.client.Database("go-project").Collection("users")

	if err := collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.ID = primitive.NewObjectID()

	collection := uc.client.Database("go-project").Collection("users")

	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)

	collection := uc.client.Database("go-project").Collection("users")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user", oid, "\n")
}

func (uc UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := uc.client.Database("go-project").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	uj, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}
