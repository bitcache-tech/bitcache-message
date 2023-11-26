import { BitCacheMessage, KeyUsageMetadata } from './bcm';
import * as fs from 'fs';
import * as path from 'path';

describe('BitCacheMessage', () => {
    const testDataDir = '../data/bitcache-messages';
    const jsonAndTxDataDir = '../data';

    const removeFileExtension = (filename: string): string => {
        return filename.replace(path.extname(filename), '');
    };

    test('Deserialize binary files and check with JSON and transaction files', () => {
        const binaryFiles = fs.readdirSync(testDataDir)
            .filter(file => file.endsWith('.bin'))
            .map(file => path.join(testDataDir, file));

        for (const file of binaryFiles) {
            const inputBCMBytes = fs.readFileSync(file);
            const decodedBCM = BitCacheMessage.fromBytes(inputBCMBytes);

            const baseFilename = removeFileExtension(path.basename(file));
            const expectedTxBytes = fs.readFileSync(path.join(jsonAndTxDataDir, baseFilename + '.txn'));
            const expectedTxHex = expectedTxBytes.toString('hex');
            const expectedKUMs: KeyUsageMetadata[] = JSON.parse(fs.readFileSync(path.join(jsonAndTxDataDir, baseFilename + '.json'), 'utf-8'));

            expect(decodedBCM.Tx).toEqual(expectedTxHex);
            expect(decodedBCM.KUMs).toEqual(expectedKUMs);
        }
    });

    test('Encode JSON and transaction files to binary and compare with binary files', () => {
        const jsonFiles = fs.readdirSync(jsonAndTxDataDir)
            .filter(file => file.endsWith('.json'))
            .map(file => path.join(jsonAndTxDataDir, file));

        for (const file of jsonFiles) {
            const baseFilename = removeFileExtension(path.basename(file));
            const inputTxBytes = fs.readFileSync(path.join(jsonAndTxDataDir, baseFilename + '.txn'));
            const inputKUMs: KeyUsageMetadata[] = JSON.parse(fs.readFileSync(file, 'utf-8'));

            const bcm = new BitCacheMessage(inputTxBytes.toString('hex'), inputKUMs);
            const encodedData = bcm.toBytes();

            const expectedBCMBytes = fs.readFileSync(path.join(testDataDir, baseFilename + '.bin'));
            expect(encodedData).toEqual(expectedBCMBytes);
        }
    });
});
