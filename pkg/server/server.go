package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	Port  int
	Mux   *mux.Router
	Store map[uuid.UUID]string
	srv   http.Server
}

// NewServer
func NewServer(port int) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/ping", pingHandler).Methods("GET")
	return &Server{
		Port:  port,
		Mux:   router,
		Store: make(map[uuid.UUID]string),
		srv: http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: router,
		},
	}
}

func NewID() uuid.UUID {
	return uuid.New()
}

func (S *Server) Listen() error {
	return S.srv.ListenAndServe()
}

func (S *Server) Register() {
	S.Mux.HandleFunc("/upload", S.uploadHandler).Methods("POST")
	S.Mux.HandleFunc("/getImage", S.getHandler).Methods("GET")
}
func pingHandler(w http.ResponseWriter, r *http.Request) {
	pongResponse := []byte("pong")
	w.Write(pongResponse)
}

func (S *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	id := NewID()
	ioutil.WriteFile(id.String()+".png", contents, 0644)
	S.Store[id] = id.String() + ".png"
	w.Write([]byte(id.String()))
}

func (S *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		fmt.Println(err)
	}

	if _, ok := S.Store[uid]; !ok {
		errorMessage := []byte("Image not present")
		w.Write(errorMessage)
	}

	image, err := ioutil.ReadFile(id + ".png")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(image)
}
