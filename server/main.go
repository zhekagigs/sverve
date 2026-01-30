package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"kvstore/server/db"
)

var index *db.KeyIndex

func jsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type GetReq struct {
	Key string `json:"key"`
}

type GetResp struct {
	Value string `json:"value"`
	Found bool   `json:"found"`
}

type SetReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResp struct {
	Success bool `json:"success"`
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	var req GetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	t1 := time.Now()

	keyHash := db.HashKey(req.Key)
	offset, exists := index.Get(keyHash)
	var value string
	if exists {
		value, _ = db.DbGet(req.Key, offset)
	}
	t2 := time.Now()

	fmt.Printf("GET %s: decode=%v query=%v total=%v\n", req.Key, t1.Sub(t0), t2.Sub(t1), t2.Sub(t0))

	jsonResponse(w, GetResp{Value: value, Found: exists})
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	var req SetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := db.DbSet(req.Key, req.Value); err != nil {
		http.Error(w, "Failed to set key", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, SetResp{Success: true})
}

func main() {
	var err error

	// Open database file for reads
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	tStart := time.Now()
	index, err = db.DbScan()
	if err != nil {
		panic(err)
	}
	tEnd := time.Now()
	timeDiff := tEnd.Sub(tStart)
	fmt.Printf("Database index loaded in %v\n", timeDiff)

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/set", handleSet)

	fmt.Println("Database server listening on :9000")
	http.ListenAndServe(":9000", nil)
}
