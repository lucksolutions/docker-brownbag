package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"net/http"

	"time"

	"io/ioutil"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type config struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
}

/*
Guest of the docker brownbag
*/
type Guest struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Date      string `json:"time"`
}

/*
AllGuestResponse to request for all guests
*/
type AllGuestResponse struct {
	Guest []Guest `json:"guest"`
}

var version = "1.0.1"
var dbConfig config
var timeFormat = "2006-01-02 15:04:05"

func getDbConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbConfig.dbUser, dbConfig.dbPassword, dbConfig.dbHost,
		dbConfig.dbPort, dbConfig.dbName)
	db, err := sql.Open("mysql", connectionString)
	if err == nil {
		err = db.Ping()
	}

	return db, err
}

func getGuestList() ([]Guest, error) {
	db, err := getDbConnection()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("select uuid, date, first_name, last_name from registry")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var date string
		var firstName string
		var lastName string
		var uuid string

		err = rows.Scan(&uuid, &date, &firstName, &lastName)
		if err != nil {
			log.Println(err)
		}
		guest := Guest{
			ID:        uuid,
			Date:      date,
			Firstname: firstName,
			Lastname:  lastName,
		}
		guests = append(guests, guest)
	}

	return guests, err

}

func getGuestsEndpoint(w http.ResponseWriter, r *http.Request) {
	guests, err := getGuestList()
	response := AllGuestResponse{guests}
	statusCode := http.StatusOK
	if err != nil {
		log.Println(err)
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, err := json.MarshalIndent(&response, "", " ")
	if err != nil {
		log.Println(err)
	}
	w.Write(data)
}

func createGuestEndpoint(w http.ResponseWriter, r *http.Request) {
	var guest Guest
	_ = json.NewDecoder(r.Body).Decode(&guest)
	guest.Date = time.Now().Format(timeFormat)
	guest.ID = uuid.New().String()

	db, err := getDbConnection()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	statement, _ := db.Prepare("INSERT INTO registry(uuid, date, first_name, last_name) VALUES(?, ?, ?, ?)")
	_, err = statement.Exec(guest.ID, guest.Date, guest.Firstname, guest.Lastname)
	defer statement.Close()
	if err != nil {
		log.Println(err)
	}

	statusCode := http.StatusOK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, err := json.MarshalIndent(&guest, "", " ")
	if err != nil {
		log.Println(err)
	}
	w.Write(data)
}

func versionEndpoint(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("{version: '%s'}", version)))
}

func healthzEndpoint(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK

	db, err := getDbConnection()
	if err != nil {
		statusCode = http.StatusInternalServerError
		log.Printf("Health Check failed: %s\n", err)
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{status: 'healthy'}"))
}

func main() {
	log.Println("Starting app...")

	/*
	   Configure the application here. This could include reading ENV variables, loading
	   a configuration files, etc.
	*/
	httpAddr := os.Getenv("HTTP_ADDR")
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPasswordFile := os.Getenv("DB_PASSWORD_FILE")
	dbName := os.Getenv("DB_DATABASE")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	if dbPasswordFile != "" {
		log.Printf("DB_PASSWORD_FILE specified. Reading password from %s", dbPasswordFile)
		fileData, err := ioutil.ReadFile(dbPasswordFile)
		if err != nil {
			log.Fatalf("Could not read database password form %s: %s", dbPasswordFile, err)
		}
		dbPassword = string(fileData)
		//Trim trailing newline character
		dbPassword = strings.TrimRight(dbPassword, "\n")
	}

	dbConfig = config{
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUsername,
		dbPassword: dbPassword,
		dbName:     dbName,
	}
	if dbConfig.dbPort == 0 {
		dbConfig.dbPort = 3306
	}

	router := mux.NewRouter()
	router.HandleFunc("/version", versionEndpoint).Methods("GET")
	router.HandleFunc("/healthz", healthzEndpoint).Methods("GET")
	router.HandleFunc("/guest", getGuestsEndpoint).Methods("GET")
	router.HandleFunc("/guest", createGuestEndpoint).Methods("POST")

	/* Default Handler */
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html, hostname, version)
	}).Methods("GET")

	/*
		Start the HTTP Server
	*/
	log.Printf("HTTP Service listening on %s", httpAddr)
	httpErr := http.ListenAndServe(httpAddr, router)
	if httpErr != nil {
		log.Fatal(httpErr)
	}

}
