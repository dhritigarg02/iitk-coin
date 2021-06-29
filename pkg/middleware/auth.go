package middleware

import(
	"net/http"
	"log"

	"github.com/dhritigarg02/iitk-coin/pkg/auth"
	"github.com/dhritigarg02/iitk-coin/pkg/db"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, int)

type EnsureAuth struct {
    Handler AuthenticatedHandler
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	payload, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

    ea.Handler(w, r, payload.Rollno)
}

type EnsureAdmin struct {
    Handler http.HandlerFunc
	DBstore db.DBStore
}

func (ea *EnsureAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	payload, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	isAdmin, err := ea.DBstore.CheckAdmin(payload.Rollno)
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