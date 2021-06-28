package middleware

import(
	"net/http"
	"log"

	"github.com/dhritigarg02/iitk-coin/pkg/auth"
	"github.com/dhritigarg02/iitk-coin/pkg/db"
)

type EnsureAuth struct {
    Handler http.HandlerFunc
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

    ea.Handler(w, r)
}

type EnsureAdmin struct {
    Handler http.HandlerFunc
	Dbstore db.DBStore
}

func (ea *EnsureAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	payload, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	isAdmin, err := ea.Dbstore.CheckAdmin(payload.Rollno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[EnsureAdmin] [ERROR] : %v\n", err)
		return
	}

	if !isAdmin {
		http.Error(w, "You do not have sufficient permissions", http.StatusUnauthorized)
		return
	}

    ea.Handler(w, r)
}