// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Customer struct {
	CId           int       `json:"cID"`
	Cname         string    `json:"cName"`
	Caddress      string    `json:"cAddress"`
	CTel          int       `json:"cTel"`
	CregisterDate time.Time `json:"cRegisterDate"`
}

type AllDataToSend struct {
	Size      int        `json:"size"`
	Customers []Customer `json:"customers"`
	Msg       string     `json:"msg"`
}

type OneCustomerToSend struct {
	CId           int       `json:"cID"`
	Cname         string    `json:"cName"`
	Caddress      string    `json:"cAddress"`
	CTel          int       `json:"cTel"`
	CregisterDate time.Time `json:"cRegisterDate"`
	Msg           string    `json:"msg"`
}

type Report struct {
	TotalCustomers int    `json:"totalCustomers`
	Period         int    `json:"period"`
	Msg            string `json:"msg"`
}

type Message struct {
	Msg string `json:"msg"`
}

var Customers []Customer

var indexer int = 0

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllCustomers(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Endpoint Hit: returnAllArticles")

	var strTemp string

	var keyName string

	var lengOfInput int
	r.ParseForm()
	for key, value := range r.Form {
		// fmt.Printf("%s = %s", key, value)
		strTemp = value[0]
		keyName = key
		lengOfInput = len(value)
		break
	}

	// json.NewEncoder(w).Encode(strTemp)

	if lengOfInput != 0 && keyName == "cName" {
		var flag int = 0
		for _, customer := range Customers {
			if strings.Contains(customer.Cname, strTemp) == true {
				var customerToSend OneCustomerToSend
				customerToSend.CId = customer.CId
				customerToSend.Caddress = customer.Caddress
				customerToSend.Cname = customer.Cname
				customerToSend.CTel = customer.CTel
				customerToSend.CregisterDate = customer.CregisterDate
				customerToSend.Msg = "success"
				json.NewEncoder(w).Encode(customerToSend)
				flag = 1
			}
		}

		if flag == 0 {
			var message Message
			message.Msg = "error"
			json.NewEncoder(w).Encode(message)
		}
	} else {
		var toSent AllDataToSend
		toSent.Customers = Customers
		toSent.Size = len(Customers)

		if len(Customers) != 0 {
			toSent.Msg = "success"
			json.NewEncoder(w).Encode(toSent)
		} else {
			var message Message
			message.Msg = "error"
			json.NewEncoder(w).Encode(message)
		}
	}
}

func returnSingleCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var temp int
	temp, _ = strconv.Atoi(key)

	var flag int = 0
	for _, customer := range Customers {
		if customer.CId == temp {
			var customerToSend OneCustomerToSend
			customerToSend.CId = customer.CId
			customerToSend.Caddress = customer.Caddress
			customerToSend.Cname = customer.Cname
			customerToSend.CTel = customer.CTel
			customerToSend.CregisterDate = customer.CregisterDate
			customerToSend.Msg = "success"
			json.NewEncoder(w).Encode(customerToSend)
			flag = 1
		}
	}
	if flag == 0 {
		var message Message
		message.Msg = "error"
		json.NewEncoder(w).Encode(message)
	}
}

func createNewCustomer(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var customer Customer

	json.Unmarshal(reqBody, &customer)

	indexer += 1
	customer.CId = indexer

	customer.CregisterDate = time.Now()

	Customers = append(Customers, customer)

	var customerToSend OneCustomerToSend
	customerToSend.CId = customer.CId
	customerToSend.Caddress = customer.Caddress
	customerToSend.Cname = customer.Cname
	customerToSend.CTel = customer.CTel
	customerToSend.CregisterDate = customer.CregisterDate
	customerToSend.Msg = "success"

	json.NewEncoder(w).Encode(customerToSend)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var temp int
	temp, _ = strconv.Atoi(id)

	var flag int = 0

	for index, customer := range Customers {
		if customer.CId == temp {
			Customers = append(Customers[:index], Customers[index+1:]...)

			var message Message
			message.Msg = "success"
			json.NewEncoder(w).Encode(message)
			flag = 1
		}
	}
	if flag == 0 {
		var message Message
		message.Msg = "error"
		json.NewEncoder(w).Encode(message)
	}

}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var temp int
	temp, _ = strconv.Atoi(id)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var customerTemp Customer

	json.Unmarshal(reqBody, &customerTemp)

	var flag int = 0

	for index, customer := range Customers {
		if customer.CId == temp {
			Customers[index].Cname = customerTemp.Cname
			Customers[index].Caddress = customerTemp.Caddress
			Customers[index].CTel = customerTemp.CTel

			var customerToSend OneCustomerToSend
			customerToSend.CId = Customers[index].CId
			customerToSend.Caddress = Customers[index].Caddress
			customerToSend.Cname = Customers[index].Cname
			customerToSend.CTel = Customers[index].CTel
			customerToSend.CregisterDate = Customers[index].CregisterDate
			customerToSend.Msg = "success"
			json.NewEncoder(w).Encode(customerToSend)
			flag = 1
		}
	}
	if flag == 0 {
		var message Message
		message.Msg = "error"
		json.NewEncoder(w).Encode(message)
	}

}

func returnReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	month := vars["month"]

	var temp int
	temp, _ = strconv.Atoi(month)

	var tempMonth string

	switch temp {
	case 0:
		tempMonth = "January"
	case 1:
		tempMonth = "February"
	case 2:
		tempMonth = "March"
	case 3:
		tempMonth = "April"
	case 4:
		tempMonth = "May"
	case 5:
		tempMonth = "June"
	case 6:
		tempMonth = "July"
	case 7:
		tempMonth = "August"
	case 8:
		tempMonth = "September"
	case 9:
		tempMonth = "October"
	case 10:
		tempMonth = "November"
	case 11:
		tempMonth = "December"
	}

	var numOfCustomer int = 0
	for _, customer := range Customers {
		// json.NewEncoder(w).Encode(customer.CregisterDate.Month().String())
		if customer.CregisterDate.Month().String() == tempMonth {
			// json.NewEncoder(w).Encode(tempMonth)

			numOfCustomer += 1
		}
	}

	var reportToSend Report
	reportToSend.TotalCustomers = numOfCustomer
	reportToSend.Period = temp
	reportToSend.Msg = "success"
	json.NewEncoder(w).Encode(reportToSend)

}

// func returnCustomerByName(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.
// 	id := vars["id"]

// 	vars := mux.Vars(r)
// 	key := vars["cName"]

// 	var flag int = 0
// 	for _, customer := range Customers {
// 		if strings.Contains(customer.Cname, key) == true {
// 			var customerToSend OneCustomerToSend
// 			customerToSend.CId = customer.CId
// 			customerToSend.Caddress = customer.Caddress
// 			customerToSend.Cname = customer.Cname
// 			customerToSend.Cphone = customer.Cphone
// 			customerToSend.CregisterDate = customer.CregisterDate
// 			customerToSend.Msg = "success"
// 			json.NewEncoder(w).Encode(customerToSend)
// 			flag = 1
// 		}
// 	}
// 	if flag == 0 {
// 		json.NewEncoder(w).Encode("error")
// 	}
// }

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/customers", returnAllCustomers).Methods("GET")
	myRouter.HandleFunc("/customers/{id}", returnSingleCustomer).Methods("GET")
	myRouter.HandleFunc("/report/{month}", returnReport).Methods("GET")
	// myRouter.HandleFunc("/customers/{cName}", returnCustomerByName).Methods("GET")

	myRouter.HandleFunc("/customers", createNewCustomer).Methods("POST")
	myRouter.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	myRouter.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	// Customers = []Customer{
	// 	Customer{CId: 1, Cname: "Hell23o", Caddress: "Article Description", Cphone: 8765465, CregisterDate: time.Now()},
	// 	Customer{CId: 2, Cname: "Hello 2", Caddress: "Article Description", Cphone: 67498494, CregisterDate: time.Now()},
	// }
	handleRequests()
}
