package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	// "path/filepath"

	"github.com/gorilla/mux"
	"github.com/kei-gnu/golang_http/logger"
)

// HTTPサーバのインスタンス Logのインスタンス保持する
type httpServer struct {
	Log *logger.Log
}

// 
type ProduceRequest struct {
	Record logger.Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record logger.Record `json:"record"`
}

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr: addr,
		Handler: r,
	}
}


func newHTTPServer() *httpServer {
	fmt.Printf("newHTTPServer:\n")
	return &httpServer{
		Log: logger.NewLog(),
	}
}

func WriteAccessLog(filename string, r *http.Request) error {
	// f, err := os.Create(filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err !=nil {
		fmt.Println(err)
		fmt.Println("fail to create file")
		return nil 
	}
	defer f.Close()
	
	request_log := fmt.Sprintf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	count, err := f.Write([]byte(request_log))

	if err !=nil {
		fmt.Println(err)
		fmt.Println("fail to write file")
		return nil 
	}

	fmt.Printf("write %d bytes\n", count)
	return nil
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Printf("*http.request: %v\n", r)
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)    // 意味わからんコードだ Decodeってなんだ？
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Printf("handleConsume,http.Request: %v\n", r)
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)    // 意味わからんコードだ Decodeってなんだ？
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	record, err := s.Log.Read(req.Offset)
	if err == logger.ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return 
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	filename := "access_log"
	fmt.Printf("filename: %v", filename)
	err = WriteAccessLog(filename, r)
	if err !=nil {
		fmt.Printf(err.Error())
	}
}