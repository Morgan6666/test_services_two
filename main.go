package main

// title Orders API
// @version 1.0
// @description This is a first service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email dbairamkulow@mail.ru
// @host localhost:8081
// @BasePath /

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"net/http/httputil"
	"os"
	_ "os"
	"time"
	_ "time"
)

type NGS struct {
	Id  string `json:"Id"`
	DNA string `json:"DNA"`
	RNA string `json:"RNA"`
}

var ngs []NGS

// returnAllArticles godoc
// @SAccept json
// @Router /all

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		select {
		case <-time.After(10 * time.Second):
			log.Println(fmt.Sprintf("%q", x))
			f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fmt.Fprintln(f, fmt.Sprintf("%q", x))
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("Endpoint Hit: returnAllArticles")
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("Save")

	}
	json.NewEncoder(w).Encode(ngs)

}

// homePage godoc
// @SAccept json
// @Router /

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcom to the HomePage")
	fmt.Println("Endpoint Hit: homePage")
}

// createNewArticle godoc
// @SAccept json
// @Router /article [post]
func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var seq NGS
	json.Unmarshal(reqBody, &seq)
	ngs = append(ngs, seq)
	json.NewEncoder(w).Encode(seq)
}

// deleteArticle godoc
// @SAccept json
// @Router /article/{id} [delete]
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range ngs {
		if article.Id == id {
			ngs = append(ngs[:index], ngs[index+1:]...)
		}
	}
}

// returnSingleArticle godoc
// @SAccept json
// @Router /article/{id}
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	for _, seq := range ngs {
		if seq.Id == key {
			json.NewEncoder(w).Encode(seq)
		}
	}

}

func handleRequests() {
	mainRouter := mux.NewRouter().StrictSlash(true)

	mainRouter.HandleFunc("/", logHandler(homePage))
	mainRouter.HandleFunc("/seq", logHandler(returnAllArticles))
	mainRouter.HandleFunc("/data/{id}", logHandler(returnSingleArticle))
	mainRouter.HandleFunc("/data", logHandler(createNewArticle)).Methods("POST")
	mainRouter.HandleFunc("/data/{id}", logHandler(deleteArticle)).Methods("Delete")

	log.Fatal(http.ListenAndServe(":8081", mainRouter))

}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	ngs = []NGS{
		NGS{Id: "1", DNA: "ATGC", RNA: "UTGC"},
		NGS{Id: "2", DNA: "TGCA", RNA: "TGCU"},
		NGS{Id: "3", DNA: "AAAAA", RNA: "UUUUU"},
	}

	handleRequests()
}
