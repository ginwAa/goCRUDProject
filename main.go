package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"project1/conf"
	"project1/entity"
	"project1/handler/UserHandler"
	"time"
)

func SetCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Method, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Claims struct {
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		cookie, err := r.Cookie("TKARL")
		if err != nil || (!cookie.Expires.IsZero() && cookie.Expires.Before(time.Now())) {
			log.Println("cookie failed :" + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return entity.SecretKey, nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if err != nil || !ok || !token.Valid ||
			(fmt.Sprint(claims["role"]) != "admin" && fmt.Sprint(claims["role"]) != "normal") {
			log.Println("jwt failed :" + tokenStr + " : " + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		cookie, err := r.Cookie("TKARL")
		if err != nil || (!cookie.Expires.IsZero() && cookie.Expires.Before(time.Now())) {
			log.Println("cookie failed :" + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return entity.SecretKey, nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if err != nil || !ok || !token.Valid {
			log.Println("jwt failed :" + tokenStr + " : " + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if fmt.Sprint(claims["role"]) != "admin" {
			log.Println("role not admin")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
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

	http.Handle("/user/login", SetCorsHeaders(http.HandlerFunc(UserHandler.Login)))
	http.Handle("/user/register", SetCorsHeaders(http.HandlerFunc(UserHandler.Register)))
	http.Handle("/user/account-unique", SetCorsHeaders(http.HandlerFunc(UserHandler.AccountUnique)))
	http.Handle("/user/page", SetCorsHeaders(Auth(http.HandlerFunc(UserHandler.Page))))
	http.Handle("/user/count", SetCorsHeaders(Auth(http.HandlerFunc(UserHandler.Count))))
	http.Handle("/user/add", SetCorsHeaders(AuthAdmin(http.HandlerFunc(UserHandler.Add))))
	http.Handle("/user/update", SetCorsHeaders(AuthAdmin(http.HandlerFunc(UserHandler.Update))))
	http.Handle("/user/delete", SetCorsHeaders(AuthAdmin(http.HandlerFunc(UserHandler.Delete))))

	err = http.ListenAndServe(":8081", nil)
}

/*
TODO
新vue拆包
1、节流防抖

*/
