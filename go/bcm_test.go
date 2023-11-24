package bcm_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"log"
	"path/filepath"

	bcm "github.com/bitcache-tech/bitcache-message"
	"github.com/libsv/go-bt/v2"
)

func GetAllFilesFromDirectory(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func TestNewBitCacheMessage(t *testing.T) {
	// Get filename list for test data BitCache messages in binary format
	files, err := GetAllFilesFromDirectory("../data/bitcache-messages")
	if err != nil {
		t.Errorf("Error getting files from directory: %v", err)
	}

	for _, file := range files {
		// Decode BitCacheMessage from bytes
		inputBCMBytes, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		decodedBCM, err := bcm.NewBitCacheMessage(inputBCMBytes)
		if err != nil {
			log.Printf("Error decoding file %q: %v", file, err)
		}

		// Get expected Tx and KUMs from JSON and tx bytes files
		filenameWithoutExt := filepath.Base(file)[:len(filepath.Base(file))-len(filepath.Ext(file))]
		expectedTxBytes, err := os.ReadFile(fmt.Sprintf("../data/%s.txn", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		expectedJSON, err := os.ReadFile(fmt.Sprintf("../data/%s.json", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		var expectedKUMs []*bcm.KeyUsageMetadata
		json.Unmarshal(expectedJSON, &expectedKUMs)

		// Compare decoded BitCacheMessage with expected values
		if !bytes.Equal(decodedBCM.Tx.Bytes(), expectedTxBytes) {
			t.Errorf("Decoded BitCacheMessage Tx does not match original for: " + file)

		}
		if !reflect.DeepEqual(decodedBCM.KUMs, expectedKUMs) {
			t.Errorf("Decoded BitCacheMessage KUM does not match original for: " + file)
		}
	}

}

func TestBitCacheMessageBytes(t *testing.T) {
	// Get filename list for test data BitCache messages in binary format
	files, err := GetAllFilesFromDirectory("../data/bitcache-messages")
	if err != nil {
		t.Errorf("Error getting files from directory: %v", err)
	}

	for _, file := range files {
		// Get input Tx and KUMs from JSON and tx bytes files
		filenameWithoutExt := filepath.Base(file)[:len(filepath.Base(file))-len(filepath.Ext(file))]
		inputTxBytes, err := os.ReadFile(fmt.Sprintf("../data/%s.txn", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		inputJSON, err := os.ReadFile(fmt.Sprintf("../data/%s.json", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}

		// Build BitCacheMessage object from Tx and KUMs
		inputTx, err := bt.NewTxFromBytes(inputTxBytes)
		if err != nil {
			t.Errorf("Error creating Tx from bytes for file %q: %v", file, err)
			continue
		}
		var inputKUMs []*bcm.KeyUsageMetadata
		json.Unmarshal(inputJSON, &inputKUMs)
		inputBCM := &bcm.BitCacheMessage{
			Tx:   inputTx,
			KUMs: inputKUMs,
		}

		// Encode BitCacheMessage to bytes
		encodedData, err := inputBCM.Bytes()
		if err != nil {
			t.Errorf("Error regenerating bytes from BitCacheMessage for file %q: %v", file, err)
			continue
		}

		// Get original bytes from file
		originalData, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}

		// Compare encoded bytes with original bytes
		if !bytes.Equal(originalData, encodedData) {
			t.Errorf("Encoded bytes do not match original bytes for file: %q", file)
		}
	}
}
