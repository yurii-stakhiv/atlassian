package main

import (
	"encoding/json"
	"fmt"
	tokenizer "github.com/yurii-stakhiv/atlassian"
	"net/http"
	"os"
)

func showUsage() {
	fmt.Println("./server <host:port>")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
	}

	http.HandleFunc("/tokenize", TokenizeHander)
	http.ListenAndServe(os.Args[1], nil)
}

type TokenizeReq struct {
	Input string `json:"input"`
}

func TokenizeHander(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	req := &TokenizeReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad request")
		return
	}

	scanner := tokenizer.NewScanner(nil)
	resp, err := scanner.ScanBytes([]byte(req.Input))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
