package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./client [get|set] [key] [value]")
		fmt.Println("  get <key>          - Retrieve value for key")
		fmt.Println("  set <key> <value>  - Store key-value pair")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Error: get requires a key")
			os.Exit(1)
		}
		key := os.Args[2]

		body := map[string]string{"key": key}
		bodyJSON, _ := json.Marshal(body)

		resp, err := http.Post("http://localhost:9000/get", "application/json", bytes.NewBuffer(bodyJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result["found"].(bool) {
			fmt.Println(result["value"])
		} else {
			fmt.Printf("Key '%s' not found\n", key)
		}

	case "set":
		if len(os.Args) < 4 {
			fmt.Println("Error: set requires key and value")
			os.Exit(1)
		}
		key := os.Args[2]
		value := os.Args[3]

		body := map[string]string{"key": key, "value": value}
		bodyJSON, _ := json.Marshal(body)

		resp, err := http.Post("http://localhost:9000/set", "application/json", bytes.NewBuffer(bodyJSON))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		resp.Body.Close()

		fmt.Printf("Set %s = %s\n", key, value)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
