package main

import (
	"./model"
	"./webserver"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var databaseHost string

func main() {
	// TODO:
	// - Run in kubernetes
	// - MongoDB runs with StatefulSet
	// - Build with Jenkins (Scripted, docker, podTemplate)

	fmt.Println("> Preparing initial data")
	databaseHost = os.Getenv("DB")
	if len(databaseHost) == 0 {
		fmt.Println("> Database host not set, falling back to 'localhost'")
		databaseHost = "localhost"
	}
	trips := prepareInitialData()

	c := make(chan bool)
	webserverData := &webserver.WebserverData{Trips: trips}
	port := "8888"
	go webserver.StartServer(port, webserverData, c)

	fmt.Printf("> Started the web server on port %s\n", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for i := 1; ; i++ { // this is still infinite
		t := time.NewTicker(time.Second * 30)
		select {
		case <-stop:
			fmt.Println("> Shutting down polling")
			break
		case <-t.C:
			continue
		}
		break // only reached if the quitCh case happens
	}
	fmt.Println("> Shutting down webserver")
	c <- true
	if b := <-c; b {
		fmt.Println("> Webserver shut down")
	}
	fmt.Println("> Shut down app")

}

func prepareInitialData() []model.Trip {
	db := databaseHost + ":27017"
	fmt.Printf("> Connecting to MongoDB @: %s\n", db)
	session, err := mgo.Dial(db)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("go-travel").C("trips")

	netherlandsTrip := createTheNetherlandsTrip()
	tripToFind := "Netherlands 2018"

	count, err := c.Find(bson.M{"name": tripToFind}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		fmt.Println("> Did not find the Netherlands trip, creating")
		err = c.Insert(&netherlandsTrip)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("> Trip already exists, not recreating")
	}

	//result := model.Trip{}
	//err = c.Find(bson.M{"name": "Netherlands 2018"}).One(&result)
	//if err != nil {
	//    log.Fatal(err)
	//}
	//
	//fmt.Println("Trip:", result)
	var response []model.Trip
	err = c.Find(bson.M{}).All(&response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func createTheNetherlandsTrip() model.Trip {

	activityKroket := model.Activity{
		Name:         "Broodje Kroket",
		ActivityType: "Eating",
		Description:  "Eat a Broodje Kroket",
		Location:     "Amsterdam",
		Order:        1,
	}

	visitRijksmuseum := model.Activity{
		Name:         "Rijksmuseum",
		ActivityType: "Museum",
		Description:  "Visit the Rijksmuseum",
		Location:     "Amsterdam",
		Order:        2,
	}
	activities := make([]model.Activity, 2)
	activities[0] = activityKroket
	activities[1] = visitRijksmuseum

	nemo := model.PointOfInterest{
		Order:       3,
		Location:    "Amsterdam",
		Description: "Science Museum",
		Name:        "Nemo",
	}
	pointsAmsterdam := make([]model.PointOfInterest, 1)
	pointsAmsterdam[0] = nemo

	amsterdamTrip := model.Trip{
		Name:             "Amsterdam",
		Location:         "Amsterdam",
		Activities:       activities,
		LocationType:     "City",
		PointsOfInterest: pointsAmsterdam,
	}

	utrecht := model.PointOfInterest{
		Order:       4,
		Location:    "Utrecht",
		Description: "Old town",
		Name:        "Utrecht",
	}
	points := make([]model.PointOfInterest, 1)
	points[0] = utrecht

	amsterdamToUtrecht := model.Travel{
		Order: 0,
		From:  "Amsterdam",
		Too:   "Utrecht",
		Mode:  "Train",
	}
	travels := make([]model.Travel, 1)
	travels[0] = amsterdamToUtrecht

	itineraries := make([]model.Trip, 1)
	itineraries[0] = amsterdamTrip

	netherlandsTrip := model.Trip{
		PointsOfInterest: points,
		LocationType:     "Country",
		Location:         "The Netherlands",
		Name:             "Netherlands 2018",
		TripType:         "Country",
		Itineraries:      itineraries,
		Travels:          travels,
	}
	return netherlandsTrip

}
