package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"project1/conf"
	"project1/handler/UserHandler"
)

func SetCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	conf.PGDB, err = sql.Open("postgres", conf.PostgresConnectionUrl)
	defer conf.PGDB.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hello postgres!")

	http.Handle("/user/page", SetCorsHeaders(http.HandlerFunc(UserHandler.Page)))
	http.Handle("/user/add", SetCorsHeaders(http.HandlerFunc(UserHandler.Add)))
	http.Handle("/user/count", SetCorsHeaders(http.HandlerFunc(UserHandler.Count)))
	http.Handle("/user/update", SetCorsHeaders(http.HandlerFunc(UserHandler.Update)))
	http.Handle("/user/delete", SetCorsHeaders(http.HandlerFunc(UserHandler.Delete)))

	err = http.ListenAndServe(":8081", nil)
}

/*
TODO
error modal
*/
