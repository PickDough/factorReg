package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Input struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Output struct {
	A int `json:"a"`
	B int `json:"b"`
}

func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/calculate", checkCalculateInput(calculateHandler))
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":8989", buildRouter()))
}

func checkCalculateInput(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var input Input

		b := bytes.NewBuffer(nil)
		tee := io.TeeReader(r.Body, b)

		err := json.NewDecoder(tee).Decode(&input)
		if err != nil {
			http.Error(w, "Incorrect input", http.StatusBadRequest)
			return
		}

		if input.A < 0 || input.B < 0 {
			http.Error(w, "Incorrect input", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(b)

		next(w, r, ps)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	aChan := make(chan int)
	bChan := make(chan int)

	defer close(aChan)
	defer close(bChan)

	go chanFactorial(input.A, aChan)
	go chanFactorial(input.B, bChan)

	output := Output{
		A: <-aChan,
		B: <-bChan,
	}

	json.NewEncoder(w).Encode(output)
}

func chanFactorial(n int, c chan int) {
	c <- factorial(n)
}

func factorial(n int) int {
	fac := 1

	for n > 1 {
		fac *= n
		n--
	}

	return fac
}
