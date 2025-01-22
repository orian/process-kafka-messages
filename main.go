package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/orian/process-kafka-messages/message"
	"github.com/tailscale/hujson"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dirPtr := flag.String("dir", "./message/testdata", "Directory to search for JSON files")
	failOnFirst := flag.Bool("fail-on-first", false, "Fail on first failure")
	onlyFailures := flag.Bool("only-failures", false, "Only failures")

	// Parse the flags
	flag.Parse()

	// Walk the directory and process each JSON file
	callback := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			// Read the file contents
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			clean, err := hujson.Standardize(data)
			if err != nil {
				return fmt.Errorf("failed to normalize file %s: %w", path, err)
			}
			// Parse the JSON data
			var msg message.Message
			dec := json.NewDecoder(bytes.NewBuffer(clean))
			dec.DisallowUnknownFields()
			if err = dec.Decode(&msg); err != nil {
				return fmt.Errorf("failed to unmarshal JSON %s: %w", path, err)
			}

			// Process the parsed data
			if !*onlyFailures {
				fmt.Printf("Successfully loaded: %s\n", path)
			}
		}
		return nil
	}
	if !*failOnFirst {
		oldCallBack := callback
		callback = func(path string, info os.FileInfo, err error) error {
			if err := oldCallBack(path, info, err); err != nil {
				log.Printf("processing error: %s", err)
			}
			return nil
		}
	}
	err := filepath.Walk(*dirPtr, callback)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
