package web

import (
	"log"
	"net/http"
)

type StartParams struct {
	Addr             string
	DataFilePath     string
	AuthUserFilePath string
}

func Start(params StartParams) {
	loginHandler := LoginHandler{usersFilePath: params.AuthUserFilePath}
	dataHandler := DataHandler{dataFilePath: params.DataFilePath}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler.ServeHTTP)
	http.HandleFunc("/data", authMiddleware(dataHandler.ServeHTTP))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server successfuly started on addr:", params.Addr)
	if err := http.ListenAndServe(params.Addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
