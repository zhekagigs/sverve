package db

import (
	"bufio"
	"os"
	"strings"
)

const dbfile = "database.csv"

var dbFile *os.File

func Init() error {
	var err error
	dbFile, err = os.Open(dbfile)
	return err
}

func Close() error {
	if dbFile != nil {
		return dbFile.Close()
	}
	return nil
}

func DbSet(key, value string) error {
	f, err := os.OpenFile(dbfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(key + "," + value + "\n")
	return err
}

func DbGet(key string, offset int64) (string, error) {
	if dbFile == nil {
		return "", nil
	}

	dbFile.Seek(offset, 0)

	// Read line manually with fixed buffer
	buf := make([]byte, 8192)
	n, err := dbFile.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	if n == 0 {
		return "", nil
	}

	// Find newline
	line := string(buf[:n])
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}

	if strings.HasPrefix(line, key+",") {
		parts := strings.SplitN(line, ",", 2)
		if len(parts) == 2 {
			return parts[1], nil
		}
	}

	return "", nil
}

func DbScan() (*KeyIndex, error) {
	f, err := os.Open(dbfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}

	index := &KeyIndex{
		hashes:  make([]uint64, count),
		offsets: make([]int64, count),
		count:   0,
	}

	f, err = os.Open(dbfile)
	if err != nil {
		return index, err
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)
	byteOffset := int64(0)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)

		if len(parts) == 2 {
			keyHash := HashKey(parts[0])
			index.hashes[index.count] = keyHash
			index.offsets[index.count] = byteOffset
			index.count++
		}
		byteOffset += int64(len(line) + 1)
	}

	// Sort for binary search
	index.Sort()

	return index, scanner.Err()
}
