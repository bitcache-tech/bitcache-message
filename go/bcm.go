package bcm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bt/v2"
)

const TxChunkPrefix = "TXBYTES."
const KUMChunkPrefix = "TXOKEYS."

type KeyUsageMetadata struct {
	Vout           uint32 `json:"vout"`
	ScriptType     string `json:"script_type"`
	KeyFingerprint string `json:"key_fingerprint"`
	KeyDerivation  string `json:"key_derivation"`
}

type BitCacheMessage struct {
	Tx   *bt.Tx              `json:"tx"`
	KUMs []*KeyUsageMetadata `json:"key_map"`
}

func NewBitCacheMessage(bcmBytes []byte) (*BitCacheMessage, error) {
	txChunkPrefix := bcmBytes[:8]
	if !bytes.Equal(txChunkPrefix, []byte(TxChunkPrefix)) {
		return nil, fmt.Errorf("invalid tx chunk prefix")
	}

	txChunkSize := binary.LittleEndian.Uint64(bcmBytes[8:16])

	txBytes := bcmBytes[16 : 16+txChunkSize]
	tx, err := bt.NewTxFromBytes(txBytes)
	if err != nil {
		return nil, err
	}

	cursor := 16 + txChunkSize
	kumChunkPrefix := bcmBytes[cursor : cursor+8]
	if !bytes.Equal(kumChunkPrefix, []byte(KUMChunkPrefix)) {
		return nil, fmt.Errorf("invalid KUM chunk prefix")
	}

	cursor += 8
	kumChunkSize := binary.LittleEndian.Uint64(bcmBytes[cursor : cursor+8])

	cursor += 8
	kumBytes := bcmBytes[cursor : cursor+kumChunkSize]

	kums := []*KeyUsageMetadata{}
	kcursor := uint64(0)
	for kcursor < kumChunkSize {
		vout := binary.LittleEndian.Uint32(kumBytes[kcursor : kcursor+4])
		kcursor += 4

		scriptType, bytesRead := readLengthPrefixed(kumBytes, kcursor)
		kcursor += bytesRead
		kfp, bytesRead := readLengthPrefixed(kumBytes, kcursor)
		kcursor += bytesRead
		kd, bytesRead := readLengthPrefixed(kumBytes, kcursor)
		kcursor += bytesRead

		kum := &KeyUsageMetadata{
			Vout:           vout,
			ScriptType:     string(scriptType),
			KeyFingerprint: hex.EncodeToString(kfp),
			KeyDerivation:  string(kd),
		}
		kums = append(kums, kum)
	}

	return &BitCacheMessage{
		Tx:   tx,
		KUMs: kums,
	}, nil
}

func readLengthPrefixed(data []byte, cursor uint64) ([]byte, uint64) {
	lengthPrefixSize := uint64(2)
	length := binary.LittleEndian.Uint16(data[cursor : cursor+lengthPrefixSize])
	strStart := cursor + lengthPrefixSize
	strEnd := strStart + uint64(length)
	return data[strStart:strEnd], lengthPrefixSize + uint64(length)
}

func (bcm *BitCacheMessage) Bytes() ([]byte, error) {
	var buffer bytes.Buffer

	// Write TxChunkPrefix
	buffer.Write([]byte(TxChunkPrefix))

	// Serialize the transaction and write its length and bytes
	txBytes := bcm.Tx.Bytes()
	txChunkSize := uint64(len(txBytes))
	txSizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(txSizeBytes, txChunkSize)
	buffer.Write(txSizeBytes)
	buffer.Write(txBytes)

	// Write KUMChunkPrefix
	buffer.Write([]byte(KUMChunkPrefix))

	// Calculate KUM chunk size and write it later
	kumSizePos := buffer.Len()
	buffer.Write(make([]byte, 8)) // Placeholder for KUM chunk size

	// Serialize each KeyUsageMetadata and write to buffer
	for _, kum := range bcm.KUMs {
		voutBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(voutBytes, kum.Vout)
		buffer.Write(voutBytes)

		writeLengthPrefixed(&buffer, []byte(kum.ScriptType))
		kfpBytes, err := hex.DecodeString(kum.KeyFingerprint)
		if err != nil {
			return nil, err
		}
		writeLengthPrefixed(&buffer, kfpBytes)
		writeLengthPrefixed(&buffer, []byte(kum.KeyDerivation))
	}

	// Update KUM chunk size
	kumChunkSize := uint64(buffer.Len() - kumSizePos - 8)
	kumSizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(kumSizeBytes, kumChunkSize)
	copy(buffer.Bytes()[kumSizePos:], kumSizeBytes)

	return buffer.Bytes(), nil
}

func writeLengthPrefixed(buffer *bytes.Buffer, data []byte) {
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, uint16(len(data)))
	buffer.Write(lengthBytes)
	buffer.Write(data)
}
