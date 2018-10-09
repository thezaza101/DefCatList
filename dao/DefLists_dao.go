package dao

import (
	"log"

	. "../models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DefListsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "deflists"
)

// Establish a connection to database
func (m *DefListsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of deflist
func (m *DefListsDAO) FindAll() ([]DefList, error) {
	var deflists []DefList
	err := db.C(COLLECTION).Find(bson.M{}).All(&deflists)
	return deflists, err
}

// Find a deflist by its id
func (m *DefListsDAO) FindById(id string) (DefList, error) {
	var deflist DefList
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&deflist)
	return deflist, err
}

// Insert a deflist into database
func (m *DefListsDAO) Insert(deflist DefList) error {
	err := db.C(COLLECTION).Insert(&deflist)
	return err
}

// Delete an existing deflist
func (m *DefListsDAO) Delete(deflist DefList) error {
	err := db.C(COLLECTION).Remove(&deflist)
	return err
}

// Update an existing deflist
func (m *DefListsDAO) Update(deflist DefList) error {
	err := db.C(COLLECTION).UpdateId(deflist.ID, &deflist)
	return err
}
