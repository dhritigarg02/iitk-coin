package middleware

import(
	"net/http"

	"github.com/dhritigarg02/iitk-coin/pkg/auth"
)

type EnsureAuth struct {
    handler http.HandlerFunc
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	_, err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

    ea.handler(w, r)
}

func NewEnsureAuth(handlerToWrap http.HandlerFunc) *EnsureAuth {
    return &EnsureAuth{handlerToWrap}
}