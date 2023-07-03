package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quickpolls/database"

	"github.com/gorilla/mux"
)

func signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user database.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.PasswordHash, _ = HashPassword(user.PasswordHash)
	database.DB.Create(&user)

	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user database.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	var userFound database.User
	database.DB.First(&userFound, "Username = ?", user.Username)

	err := VerifyPasswordHash(user.PasswordHash, userFound.PasswordHash)

	if err == nil {
		tokenString, _ := CreateToken(&userFound)
		w.Write([]byte(tokenString))
	} else {
		w.Write([]byte(fmt.Sprint(err)))
	}
}

func checkTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		tokenClaims, err := ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Token missing: "+fmt.Sprint(err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", tokenClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("token").(*Claims)

	var poll database.Poll
	_ = json.NewDecoder(r.Body).Decode(&poll)
	poll.CreatedByID = tokenClaims.ID

	database.DB.Create(&poll)

	json.NewEncoder(w).Encode(poll)
}

func getPolls(w http.ResponseWriter, r *http.Request) {
	// tokenClaims := r.Context().Value("token").(*Claims)

	var polls []database.Poll
	database.DB.
		Preload("CreatedBy").
		Preload("Questions").
		Preload("Questions.Options").Find(&polls)

	for i := range polls {
		polls[i].CreatedBy.PasswordHash = ""
	}

	json.NewEncoder(w).Encode(polls)
}

func sendVote(w http.ResponseWriter, r *http.Request) {
	tokenClaims := r.Context().Value("token").(*Claims)

	var vote database.Vote
	_ = json.NewDecoder(r.Body).Decode(&vote)
	vote.UserID = tokenClaims.ID

	database.DB.Create(&vote)

	json.NewEncoder(w).Encode(vote)
}

func main() {
	database.Connect()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/signin", signin).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")

	s := router.PathPrefix("/api").Subrouter()
	s.Use(checkTokenMiddleware)
	s.HandleFunc("/createPoll", createPoll).Methods("POST")
	s.HandleFunc("/getPolls", getPolls).Methods("GET")
	s.HandleFunc("/sendVote", sendVote).Methods("POST")

	http.ListenAndServe(":8000", router)
}
