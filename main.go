package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func dbSet(key, value string) error {
	f, err := os.OpenFile("database.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(key + "," + value + "\n")
	return err
}

func dbGet(key string, indexMap map[string]int64) (string, error) {
	f, err := os.Open("database.csv")
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Seek(indexMap[key], 0)
	scanner := bufio.NewScanner(f)
	var value string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+",") {
			parts := strings.SplitN(line, ",", 2)
			if len(parts) == 2 {
				value = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return value, nil
}

func dbScan() (map[string]int64, error) {
	indexMap := make(map[string]int64)
	f, err := os.Open("database.csv")
	if err != nil {
		return indexMap, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	byteOffset := int64(0)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)

		if len(parts) == 2 {
			indexMap[parts[0]] = byteOffset
		}
		byteOffset += int64(len(line) + 1)  // +1 for newline
	}

	if err := scanner.Err(); err != nil {
		return indexMap, err
	}

	return indexMap, nil
}

func simple() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main [get|set] [key] [value]")
		fmt.Println("  get <key>          - Retrieve value for key")
		fmt.Println("  set <key> <value>  - Store key-value pair")
		os.Exit(1)
	}

	command := os.Args[1]
	indexMap, err := dbScan()
	if err != nil {
		panic(err)
	}
	switch command {
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Error: get requires a key")
			os.Exit(1)
		}
		key := os.Args[2]
		result, err := dbGet(key, indexMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if result != "" {
			fmt.Println(result)
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
		err := dbSet(key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set %s = %s\n", key, value)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
