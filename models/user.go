package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// using mgo.v2
// type User struct {
// 	Id     bson.ObjectId `json:"id" bson:"_id"`
// 	Name   string        `json:"name" bson:"name"`
// 	Gender string        `json:"gender" bson:"gender"`
// 	Age    int           `json:"age" bson:"age"`
// }

// using mongo-go-driver
type User struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Gender string             `json:"gender" bson:"gender"`
	Age    int                `json:"age" bson:"age"`
}
