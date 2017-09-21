package main

import (
	"net/http"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"github.com/julienschmidt/httprouter"
)

type Room struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name:"name"}
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//transfer creation request information to type value
	r := new(Room)
	errs := binding.Bind(req, r)
	if errs != nil {
		return
	}

	//create mongodb session
	session := mongoSession.Copy()
	defer session.Close()

	//create mongodb id
	r.ID = bson.NewObjectId()
	//create mongodb collection instance
	c := session.DB("test").C("rooms")

	//store room information
	if err := c.Insert(r); err != nil {
		//occur error
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusCreated, r)
}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//create mongodb session
	session := mongoSession.Copy()
	defer session.Close()

	var rooms []Room
	//retrieve all rooms
	err := session.DB("test").C("rooms").Find(nil).All(&rooms)
	if err != nil {
		//occur error
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, rooms)
}