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
    "strconv"
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

func addUserHandler(w http.ResponseWriter, r *http.Request) {

    	r.ParseForm()
    	fmt.Println(r.Form)
    
    	//bypass same origin policy
    	w.Header().Set("Access-Control-Allow-Origin", "*")
    
    	fName := r.Form["fName"][0]
    	lName := r.Form["lName"][0]
    	lat, _ := strconv.ParseFloat(r.Form["lat"][0], 32)
    	lon, _ := strconv.ParseFloat(r.Form["lon"][0], 32)
    	stackId, _ := strconv.ParseInt(r.Form["stackId"][0], 0, 64)
    	travelRange, _ := strconv.ParseInt(r.Form["travelRange"][0], 0, 64)
    	
    	users[r.Form["stackId"][0]] = makeUser(fName, lName, float32(lat), float32(lon), int(stackId), int(travelRange), 0, 0);
    	
    	fmt.Println(r.Form);
    	
    	fmt.Fprintf(w, "User added")
}

func server() {
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/adduser", addUserHandler)
    
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
    
    //fmt.Println("Sending GET request to " + publishUrl + string(b) + "...");
    
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
	users["761902"] = makeUser("Anson", "Liu", 37.75315, -122.423517, 761902, 10, 12, 3);
	users["3858448"] = makeUser("Emily", "Wang", 37.572939, -122.6, 3858448, 20, 12, 3);
	users["110707"] = makeUser("John", "Doe", 37.478182, -122.186032, 110707, 10, 12, 3);
	users["2970947"] = makeUser("Elliot", "Frisch", 37.478182, -122.186032, 2970947, 10, 12, 3);
	users["157247"] = makeUser("Thomas", "Crowder", 37.574363, -122.332224, 157247, 10, 12, 3);
	users["2767207"] = makeUser("Elliot", "Frisch", 37.3, -121.7, 2767207, 10, 12, 3);
	
	go publishPubNub(users)
	
	
	t := time.NewTicker(10 * time.Second)
	for now := range t.C {
		now = now
		go publishPubNub(users)

		//fmt.Println("Counter ", counter)
	}
	
	fmt.Println("exit")
}