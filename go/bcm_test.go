package bcm_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"log"
	"path/filepath"

	bcm "github.com/bitcache-tech/bitcache-message"
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
	files, err := GetAllFilesFromDirectory("../data/bitcache-messages")
	if err != nil {
		t.Errorf("Error getting files from directory: %v", err)
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		decoded, err := bcm.NewBitCacheMessage(data)
		if err != nil {
			log.Printf("Error decoding file %q: %v", file, err)
		}

		filename := filepath.Base(file)
		extension := filepath.Ext(file)
		filenameWithoutExt := filename[:len(filename)-len(extension)]

		expectedTxBytes, err := os.ReadFile(fmt.Sprintf("../data/%s.txn", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		if !bytes.Equal(decoded.Tx.Bytes(), expectedTxBytes) {
			t.Errorf("Decoded BitCacheMessage Tx does not match original for: " + file)

		}
		expectedJSON, err := os.ReadFile(fmt.Sprintf("../data/%s.json", filenameWithoutExt))
		if err != nil {
			t.Errorf("Error reading file %q: %v", file, err)
			continue
		}
		var expectedKUMs []*bcm.KeyUsageMetadata
		json.Unmarshal(expectedJSON, &expectedKUMs)
		if !reflect.DeepEqual(decoded.KUMs, expectedKUMs) {
			t.Errorf("Decoded BitCacheMessage KUM does not match original for: " + file)
		}

	}

}

func Test(t *testing.T) {
	prefix := "TXOKEYS."
	fmt.Println(hex.EncodeToString([]byte(prefix)))
}
