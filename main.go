package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Member struct {
	Id        string `json:"Id,omitempty"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PlanType  string `json:"PlanType"`
	Active    bool   `json:"Active"`
}

var counter int = 3

// declare a global Members array
// that we can then populate in our main function
// to simulate a database
var members []Member

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func createNewMember(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var member Member
	json.Unmarshal(reqBody, &member)
	member.Id = strconv.Itoa(counter)
	// update our global members array to include
	// our new Member
	members = append(members, member)
	counter++
	fmt.Println("Endpoint Hit: createNewMember")
	json.NewEncoder(w).Encode(member)
}

func deleteMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// extract the `id` of the member we wish to delete
	id := vars["id"]

	// loop through all our articles
	for index, member := range members {
		// if our id path parameter matches one of our members
		if member.Id == id {
			// updates our Articles array to remove the member
			members = append(members[:index], members[index+1:]...)
		}
	}
	fmt.Println("Endpoint Hit: deleteMember")
}

func returnAllMembers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllMembers")
	json.NewEncoder(w).Encode(members)
}

func returnSingleMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our members
	// if the member.Id equals the key we pass in
	// return the member encoded as JSON
	for _, member := range members {
		if member.Id == key {
			fmt.Println("Endpoint Hit: returnSingleMember")
			json.NewEncoder(w).Encode(member)
		}
	}
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var mbr Member
	var members2 []Member

	json.Unmarshal(reqBody, &mbr)

	// loop through all our articles
	for _, member := range members {
		// if our id path parameter matches one of our members
		if member.Id == mbr.Id {
			// updates our members array to update the member
			member = mbr
		}
		members2 = append(members2, member)
		members = members2
	}
	fmt.Println("Endpoint Hit: updateMember")
	json.NewEncoder(w).Encode(members)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/members", returnAllMembers)
	myRouter.HandleFunc("/member", createNewMember).Methods("POST")
	myRouter.HandleFunc("/member/{id}", updateMember).Methods("PUT")
	myRouter.HandleFunc("/member/{id}", deleteMember).Methods("DELETE")
	myRouter.HandleFunc("/member/{id}", returnSingleMember)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	members = []Member{
		Member{Id: "1", FirstName: "Tony", LastName: "Stark", PlanType: "Medicare", Active: true},
		Member{Id: "2", FirstName: "Peter", LastName: "Parker", PlanType: "Medicaid", Active: true},
	}
	handleRequests()
}
