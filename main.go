package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
)

var addr string

type Chi struct {
	Type    string `json:"type`
	Version string `json:"version"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "cemgurhan"
	password = "115115"
	dbname   = "chi-test1-db"
)

func init() {

	pflag.StringVarP(&addr, "address", "a", ":3000", "the address for the API to listen to")
	pflag.Parse()

}

func OpenConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to PSQL!")

	return db

}

func postNewChi(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()
	var c Chi
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO chi (type, version) VALUES ($1,$2)`
	_, err = db.Exec(sqlStatement, c.Type, c.Version)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func main() {

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Post("/create-chi", postNewChi)

	})

	fmt.Println("Listening on port", addr)
	http.ListenAndServe(":3000", r)

}
