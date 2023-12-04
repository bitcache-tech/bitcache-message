# BitCache Message Library

The BitCache Message Library is a TypeScript/JavaScript library designed for serializing and deserializing BitCache messages. This library provides an efficient and simple way to handle BitCache message formats, making it easier to work with them in TypeScript and JavaScript environments.

## Features

- Serialize and deserialize BitCache messages to and from binary format.
- Handle complex data structures with ease.
- Efficient memory management with dynamic buffer resizing.

## Installation

To install the BitCache Message Library, use the following command:

```bash
npm install bitcache-message
```

This will download and install the library, making it ready for use in your project.

## Usage

### Importing the Library

First, import the library into your TypeScript or JavaScript file:

```typescript
import { BitCacheMessage } from 'bitcache-message';
```

## Deserializing a BitCache Message

To deserialize a BitCache message from a binary format:

```typescript
const binaryData = ...; // Your binary data here
const message = BitCacheMessage.fromBytes(binaryData);
```

## Serializing a BitCache Message

To serialize a BitCache message to a binary format:

```typescript
const message = new BitCacheMessage("0200000001fffd773dd10fc8e8a556839338643f08382b5a9933555dea9ca0bc4bd49d25070000000049483045022100e204926ff2f9628f15b883e2728383495cc5ac9fdfd9a2a0c50cac5acdd7165e022058c639810ac28595171420eb3e11eca0dcd8a1fa9cd13eeb404a367dd1153ae941ffffffff01804a5d05000000001976a91418392a59fc1f76ad6a3c7ffcea20cfcb17bda9eb88ac00000000",
    [{
      vout: 0,
      script_type: "p2pkh",
      key_fingerprint: "9341cb4c",
      key_derivation: "bip32:m/1/20"
    },
    {
      vout: 1,
      script_type: "p2pkh",
      key_fingerprint: "9341cb4c",
      key_derivation: "bip32:m/1/21"
    }]);
const binaryData = message.toBytes();
```

## API Reference

The library provides the following main methods:

- `BitCacheMessage.fromBytes(buffer: Buffer): BitCacheMessage`
    - Deserializes a given binary buffer into a BitCacheMessage object.
- `BitCacheMessage.toBytes(): Buffer`
    - Serializes the BitCacheMessage object into a binary buffer.

## Contributing

Contributions to the BitCache Message Library are welcome! Please feel free to submit issues, pull requests, or enhancements to improve the library.

## License

This library is distributed under the ISC License. See LICENSE file in the repository for more information.
