// Package main provides the example implemetation to connect to pubnub api.
// Runs on the console.
package main

import (
	"fmt"
	"net/http"
    "io/ioutil"
    "os"
    "log"
    "time"
    "encoding/json"
)

var historyUrl = "http://pubsub.pubnub.com/history/sub-c-32127d56-0f84-11e4-baa3-02ee2ddab7fe/redwoodcity/0/1"

var publishUrl = "http://pubsub.pubnub.com/publish/pub-c-639b1fec-45b4-4d92-b246-f553bcbeb79d/sub-c-32127d56-0f84-11e4-baa3-02ee2ddab7fe/0/redwoodcity/0/"

type User struct {
    FirstName string;
    LastName string;
    Latitude float32;
    Longitude float32;
    StackId int;
    TravelRange int;
    AggregateRating float32;
    Ratings int;
}

var users map[string]User

func userHandler(w http.ResponseWriter, r *http.Request) {
    b, err := json.Marshal(users)
	if (err != nil) {
		panic(err)
	}
    
	fmt.Fprintf(w, string(b))
}

func server() {
	http.HandleFunc("/users/", userHandler)
	//http.HandleFunc("/", rootHandler)
    
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil) 
    fmt.Println("Listening on " + os.Getenv("PORT"))
    if err != nil {
      log.Fatal(err)
    }    
}

/*
func updateMap() {
    //fmt.Println("Sending GET request to " + history + "...");
    
    response, err := http.Get(url)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        fmt.Printf("updated")
        //fmt.Printf("%s\n", string(contents))
        
        latestLoc = string(contents)
    }
}
*/

func publishPubNub(data map[string]User) {
    b, err := json.Marshal(users)
	if (err != nil) {
		panic(err)
	}
    
    fmt.Println("Sending GET request to " + publishUrl + string(b) + "...");
    
    response, err := http.Get(publishUrl + string(b))
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        fmt.Printf("%s\n", string(contents))
    }
}    

func makeUser(fName string, lName string, lat float32, lon float32, stackId int, travelRange int, aggregateRating float32, ratings int) User {
    return User{fName, lName, lat, lon, stackId, travelRange, aggregateRating, ratings}
}

func main() {
	
	go server()
	
	users = make(map[string]User)
	users["761902"] = makeUser("Anson", "Liu", 37.490315, -122.223517, 761902, 10, 12, 3);
	users["110707"] = makeUser("John", "Doe", 37.478182, -122.186032, 110707, 10, 12, 3);
	
	go publishPubNub(users)
	
	
	t := time.NewTicker(10 * time.Second)
	for now := range t.C {
		now = now
		go publishPubNub(users)

		//fmt.Println("Counter ", counter)
	}
	
	fmt.Println("exit")
}