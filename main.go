package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/gddo/httputil/header"
)

const (
	Addition       string = "addition"
	Subtraction    string = "subtraction"
	Multiplication string = "multiplication"
	Unknown        string = "unknown"
)

type ArithmeticBody struct {
	OperationType string `json:"operation_type"`
	X             int32  `json:"x"`
	Y             int32  `json:"y"`
}

func bio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"slackUsername": "Kingsley Nwankwo",
			"backend":       true,
			"age":           28,
			"bio":           "I am Kingsley and I build backend stuff using golang.",
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func arithmetic(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var body ArithmeticBody
		err := dec.Decode(&body)
		if err != nil {
			msg := "Request body is not supported"
			http.Error(w, msg, http.StatusBadRequest)
		}
		var result int64
		var operation_type string
		switch body.OperationType {
		case Addition:
			result = int64(body.X) + int64(body.Y)
			operation_type = Addition
		case Subtraction:
			result = int64(body.X) - int64(body.Y)
			operation_type = Subtraction
		case Multiplication:
			result = int64(body.X) * int64(body.Y)
			operation_type = Multiplication
		default:
			http.Error(w, fmt.Sprintln("Invalid Operation Type"), http.StatusBadRequest)
			return
		}

		err2 := json.NewEncoder(w).Encode(map[string]interface{}{
			"slackUsername":  "Kingsley Nwankwo",
			"operation_type": operation_type,
			"result":         result,
		})
		if err2 != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
		// log.Fatal("$PORT must be set")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", bio)
	mux.HandleFunc("/arithmetic", arithmetic)

	err := http.ListenAndServe(":"+port, mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
