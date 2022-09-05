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
	Type    string `json:"type"`
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
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func getChiByVersion(w http.ResponseWriter, r *http.Request) {

	db := OpenConnection()

	versionParam := chi.URLParam(r, "version")

	sqlStatement := `SELECT * FROM chi WHERE version = $1`

	rows, err := db.Query(sqlStatement, versionParam)

	if err != nil {

		w.WriteHeader(442)
		w.Write([]byte(fmt.Sprintf("error fetching chi %s: %v", versionParam, err)))
		return

	}

	var chis []Chi

	for rows.Next() {

		var chi Chi
		rows.Scan(&chi.Type, &chi.Version)
		chis = append(chis, chi)

	}

	chiBytes, _ := json.MarshalIndent(chis, "", "\t")

	w.Header().Set("content-type", "application/json")
	w.Write(chiBytes)

	defer rows.Close()
	defer db.Close()

}

func getAllChi(w http.ResponseWriter, r *http.Request) {

}

func main() {

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Post("/create-chi", postNewChi)
		r.Route("/get-chi", func(r chi.Router) {

			r.Get("/", getAllChi)

			r.Route("/{version}", func(r chi.Router) {

				r.Get("/", getChiByVersion)

			})

		})
		r.Get("/get-chi/{version}", getChiByVersion) // broken

	})

	fmt.Println("Listening on port", addr)
	http.ListenAndServe(":3000", r)

}
