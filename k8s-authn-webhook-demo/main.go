package main

import (
	"encoding/json"
	"log"
	"net/http"

	authentication "k8s.io/api/authentication/v1"
)

const (
	mockUser  = "cncamp"
	mockToken = "secret101"
)

func main() {
	http.HandleFunc("/authenticate", authn)
	log.Printf("start authentication webhook at localhost:%d ...\n", 3000)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func authn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var tr authentication.TokenReview
	err := decoder.Decode(&tr)
	if err != nil {
		log.Printf("bad authn request, error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&authentication.TokenReview{
			Status: authentication.TokenReviewStatus{
				Authenticated: false,
			},
		})
		return
	}
	// mock
	log.Println("receive authn request")
	if tr.Spec.Token == mockToken {
		log.Println("authn pass as user:", mockUser)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authentication.TokenReview{
			Spec: authentication.TokenReviewSpec{},
			Status: authentication.TokenReviewStatus{
				Authenticated: true,
				User: authentication.UserInfo{
					Username: mockUser,
				},
			},
		})
		return
	}
	log.Println("authn failed with token:", tr.Spec.Token)
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(authentication.TokenReview{
		Spec: authentication.TokenReviewSpec{},
		Status: authentication.TokenReviewStatus{
			Authenticated: false,
		},
	})
}
