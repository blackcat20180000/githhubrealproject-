package main

import (
	"log"
	"net/http"
	_ "fmt"
	"os"
	"realpro"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/99designs/gqlgen/handler"
	"github.com/robfig/cron"
	"realpro/api/taskschedule"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	/************
	*** cron job 2019/5/19 edition *********
	editor:kashrin alexendera
	**********/
	c := cron.New()
	c.AddFunc("@daily", func() {
		taskschedule.Mainstr("https://www.instituteforsupplymanagement.org/ISMReport/MfgROB.cfm?SSO=1",1)
		taskschedule.Mainstr("https://www.instituteforsupplymanagement.org/ISMReport/NonMfgROB.cfm?SSO=1",2)
		taskschedule.Insertumsci()
		taskschedule.Insertbuildingshit()
		taskschedule.Esi_data()
		taskschedule.Eu_data()
	 })
	/*********
	   cron job testing
	*******/
	c.AddFunc("10 * * * * *", func() {taskschedule.GenerateError("myinitfuncc")})
	c.Start()
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Playground("GraphQL playground", "/query"))
	r.HandleFunc("/query", handler.GraphQL(realpro.NewExecutableSchema(realpro.Config{Resolvers: &realpro.Resolver{}})))
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "token", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	corsed := handlers.CORS(headersOk, originsOk, methodsOk)(r)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, corsed))
}
