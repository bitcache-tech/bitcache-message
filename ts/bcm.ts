import { Buffer } from 'buffer';

const TxChunkPrefix = "TXBYTES.";
const KUMChunkPrefix = "TXOKEYS.";

interface KeyUsageMetadata {
    vout: number;
    script_type: string;
    key_fingerprint: string;
    key_derivation: string;
}

class BitCacheMessage {
    Tx: string;
    KUMs: KeyUsageMetadata[];

    constructor(tx: string, kums: KeyUsageMetadata[]) {
        this.Tx = tx;
        this.KUMs = kums;
    }

    static fromBytes(bcmBytes: Buffer): BitCacheMessage {
        if (bcmBytes.slice(0, 8).toString() !== TxChunkPrefix) {
            throw new Error("Invalid tx chunk prefix");
        }

        const txChunkSize = bcmBytes.readUInt32LE(8);
        const txBytes = bcmBytes.slice(16, 16 + txChunkSize);

        let cursor = 16 + txChunkSize;
        if (bcmBytes.slice(cursor, cursor + 8).toString() !== KUMChunkPrefix) {
            throw new Error("Invalid KUM chunk prefix");
        }

        cursor += 8;
        const kumChunkSize = bcmBytes.readUInt32LE(cursor);
        cursor += 8;

        const kumBytes = bcmBytes.slice(cursor, cursor + kumChunkSize);

        let kcursor = 0;
        const kums: KeyUsageMetadata[] = [];
        while (kcursor < kumChunkSize) {
            const vout = kumBytes.readUInt32LE(kcursor);
            kcursor += 4;

            const [scriptType, bytesReadST] = readLengthPrefixed(kumBytes, kcursor);
            kcursor += bytesReadST;

            const [kfp, bytesReadKFP] = readLengthPrefixed(kumBytes, kcursor);
            kcursor += bytesReadKFP;

            const [kd, bytesReadKD] = readLengthPrefixed(kumBytes, kcursor);
            kcursor += bytesReadKD;

            kums.push({
                vout: vout,
                script_type: scriptType.toString(),
                key_fingerprint: kfp.toString('hex'),
                key_derivation: kd.toString(),
            });
        }

        return new BitCacheMessage(txBytes.toString('hex'), kums);
    }

    toBytes(): Buffer {
        let buffer = Buffer.alloc(1000); // Start with a buffer of 1000 bytes
        let offset = 0;

        offset += buffer.write(TxChunkPrefix, offset);
        const txLengthBigInt = BigInt(this.Tx.length / 2); // Convert hex string length to byte length
        buffer = ensureBufferSize(buffer, offset, 8);
        offset = writeUInt64LE(buffer, txLengthBigInt, offset);

        const tx = Buffer.from(this.Tx, 'hex');
        buffer = ensureBufferSize(buffer, offset, tx.length);
        offset += tx.copy(buffer, offset);

        offset += buffer.write(KUMChunkPrefix, offset);
        const kumSizePos = offset;
        offset += 8; // Reserve space for KUM chunk size

        this.KUMs.forEach(kum => {
            buffer = ensureBufferSize(buffer, offset, 4);
            buffer.writeUInt32LE(kum.vout, offset);
            offset += 4;

            const scriptTypeBuffer = Buffer.from(kum.script_type);
            buffer = ensureBufferSize(buffer, offset, scriptTypeBuffer.length + 2);
            offset = writeLengthPrefixed(buffer, scriptTypeBuffer, offset);

            const kfpBuffer = Buffer.from(kum.key_fingerprint, 'hex');
            buffer = ensureBufferSize(buffer, offset, kfpBuffer.length + 2);
            offset = writeLengthPrefixed(buffer, kfpBuffer, offset);

            const kdBuffer = Buffer.from(kum.key_derivation);
            buffer = ensureBufferSize(buffer, offset, kdBuffer.length + 2);
            offset = writeLengthPrefixed(buffer, kdBuffer, offset);
        });

        const kumChunkSize = offset - kumSizePos - 8;
        buffer = ensureBufferSize(buffer, kumSizePos, 4);
        buffer.writeUInt32LE(kumChunkSize, kumSizePos); // Write the actual KUM chunk size

        return buffer.slice(0, offset);
    }
}

function readLengthPrefixed(data: Buffer, cursor: number): [Buffer, number] {
    const length = data.readUInt16LE(cursor);
    const strStart = cursor + 2;
    const strEnd = strStart + length;
    return [data.slice(strStart, strEnd), length + 2];
}

function writeLengthPrefixed(buffer: Buffer, data: Buffer, offset: number): number {
    buffer.writeUInt16LE(data.length, offset);
    offset += 2;
    data.copy(buffer, offset);
    return offset + data.length;
}

function writeUInt64LE(buffer: Buffer, value: bigint, offset: number): number {
    const MAX_UINT32 = BigInt(0xFFFFFFFF);
    const bigValue = value / MAX_UINT32;
    const lowValue = value % MAX_UINT32;

    buffer.writeUInt32LE(Number(lowValue), offset); // Write the low 32 bits
    buffer.writeUInt32LE(Number(bigValue), offset + 4); // Write the high 32 bits

    return offset + 8; // Return the new offset
}

function ensureBufferSize(buffer: Buffer, currentOffset: number, additionalLength: number): Buffer {
    if (currentOffset + additionalLength > buffer.length) {
        const newBuffer = Buffer.alloc(buffer.length + additionalLength); // Increase the size
        buffer.copy(newBuffer);
        return newBuffer;
    }
    return buffer;
}

export { BitCacheMessage, KeyUsageMetadata };
