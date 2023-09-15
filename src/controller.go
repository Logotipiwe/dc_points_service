package main

import (
	"fmt"
	"github.com/logotipiwe/dc_go_utils/src/config"
	"net/http"
)

func main() {
	err := initializeApp()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	//user auth
	http.Handle("/get-points", appHandler(getPoints))

	//m2m auth
	http.Handle("/change-points", appHandler(changePoints))

	err = http.ListenAndServe(":"+config.GetConfig("CONTAINER_PORT"), nil)
	println("Server up!")
	if err != nil {
		panic("Lol server fell")
	}
}

type appError struct {
	Error   error
	Message string
	Code    int
}

func (a *appError) ErrorJson() string {
	return fmt.Sprintf("{\"error\":\"%s\"}", a.Message)
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appErr := fn(w, r)
	if appErr != nil {
		println(appErr.Message + " -- " + appErr.Error.Error())
		http.Error(w, appErr.ErrorJson(), appErr.Code)
	}
}

func initializeApp() error {
	config.LoadDcConfig()
	err := InitDb()
	if err != nil {
		panic(err)
	}
	return err
}
