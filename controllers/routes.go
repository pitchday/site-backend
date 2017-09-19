package controllers

import (
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/pitchday/site-backend/controllers/api"
	"github.com/pitchday/site-backend/middleware"
	"log"
	"net/http"
	"os"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func CreateRouter() http.Handler {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/{apiPrefix}").Subrouter()
	apiRouter = apiRouter.StrictSlash(true)
	apiRouter.HandleFunc("/", Use(api.API)).Methods("GET")

	apiRouter.HandleFunc("/announcements", Use(api.Get_Announcements)).Methods("GET")

	apiRouter.HandleFunc("/bots/{botType}/{botId}", Use(api.Bot_Hook)).Methods("POST")
	apiRouter.HandleFunc("/bots/userCounts", Use(api.Get_Group_Member_Count)).Methods("GET")

	apiRouter.HandleFunc("/contributors", Use(api.Get_Contributors)).Methods("GET")

	apiRouter.HandleFunc("/newsletter", Use(api.Newsletter_Subscribe)).Methods("POST")

	apiRouter.HandleFunc("/wallets", Use(api.Get_Wallet_Balance)).Methods("GET")


	// Setup CSRF Protection
	csrfHandler := nosurf.New(router)

	// Exempt API routes and Static files
	csrfHandler.ExemptGlob("/*/announcements")
	csrfHandler.ExemptGlob("/*/bots/*/*")
	csrfHandler.ExemptGlob("/*/newsletter")
	csrfHandler.ExemptGlob("/*/contributors")
	csrfHandler.ExemptGlob("/*/contributors/*")
	csrfHandler.ExemptGlob("/*/wallets")

	return Use(csrfHandler.ServeHTTP, middleware.GetContext)

	//These are here to be used by the admin panel
	//apiRouter.HandleFunc("/contributors", Use(api.Add_Contributor)).Methods("POST")
	//apiRouter.HandleFunc("/contributors/{contributorId}", Use(api.Remove_Contributor)).Methods("DELETE")

}

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
