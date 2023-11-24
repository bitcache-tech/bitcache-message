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

func (bcm *BitCacheMessage) Bytes() ([]byte, error) {
	return nil, nil
}

func readLengthPrefixed(data []byte, cursor uint64) ([]byte, uint64) {
	lengthPrefixSize := uint64(2)
	length := binary.LittleEndian.Uint16(data[cursor : cursor+lengthPrefixSize])
	strStart := cursor + lengthPrefixSize
	strEnd := strStart + uint64(length)
	return data[strStart:strEnd], lengthPrefixSize + uint64(length)
}
