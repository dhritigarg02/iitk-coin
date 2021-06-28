package api

import (
	"net/http"

	"github.com/dhritigarg02/iitk-coin/pkg/db"
	"github.com/dhritigarg02/iitk-coin/pkg/middleware"
)

type Server struct {

	DBstore db.DBStore
	Router *http.ServeMux
}

func NewEnsureAuth(handlerToWrap http.HandlerFunc) *middleware.EnsureAuth {
    return &middleware.EnsureAuth{
		Handler : handlerToWrap,
	}
}

func (server *Server) NewEnsureAdmin(handlerToWrap http.HandlerFunc) *middleware.EnsureAdmin {
    return &middleware.EnsureAdmin{
		Handler : handlerToWrap, 
		Dbstore : server.DBstore,
	}
}

func (server *Server) SetupRouter() {

	router := http.NewServeMux()

	router.HandleFunc("/", server.HelloHandler)
	router.HandleFunc("/login", server.Login)
	router.HandleFunc("/signup", server.Signup)
	router.Handle("/reward", server.NewEnsureAdmin(server.RewardCoins))
	router.Handle("/transfer", NewEnsureAuth(server.TransferCoins))
	router.Handle("/getbalance", NewEnsureAuth(server.GetBalance))

	server.Router = router
}