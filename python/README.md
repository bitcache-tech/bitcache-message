# Bitcache Python reference implementation

This package implements key functions and data structures that can be used to serialise and
deserialise bitcache messages. It was extracted from the ElectrumSV `develop` branch implementation.

## Usage

### Parsing a serialised message

A serialised message is required to be a stream of bytes. The stream is passed into the
`read_bitcache_message` function, and a `BitcacheMessage` object is returned. Currently this is
targeted around a message being the transaction bytes followed by metadata about any key usage
in the outputs of that transaction.

```python
import io
from bitcache import BitcacheMessage, read_bitcache_message

def parse_single_message_from_stream(text: str) -> None:
    stream = io.BytesIO(bytes.fromhex(text))
    message = read_bitcache_message(stream)
    print(f"Transaction size: {len(message.tx_data)} bytes")
    for keydata in message.key_data:
        print(f"Key vout:        {keydata.txo_index}")
        print(f"Key script type: {keydata.script_type}")
        print(f"Key fingerprint: {keydata.parent_key_fingerprint.hex()}")
        print(f"Key derivation:  {keydata.derivation_text}")

parse_single_message_from_stream('545842595445532e58000000000000000200000001000000000000000000000'
    '0000000000000000000000000000000000000000000ffffffff03510101ffffffff0100f2052a010000001976a91'
    '4e7d468699f0f2acbca887376739551053866617d88ac0000000054584f4b4559532e1e000000000000000000000'
    '005007032706b6804009341cb4c0b0062697033323a6d2f302f30')
```

Executing the above code gives the following result:

```text
Transaction size: 88 bytes
Key vout:        0
Key script type: p2pkh
Key fingerprint: 9341cb4c
Key derivation:  bip32:m/0/0
```

### Serialising message data

```python
import io
from bitcache import BitcacheMessage, BitcacheTxoKeyUsage, write_bitcache_transaction_message

TX_BYTES = bytes.fromhex('02000000010000000000000000000000000000000000000000000000000000000000000'
    '000ffffffff03510101ffffffff0100f2052a010000001976a914e7d468699f0f2acbca887376739551053866617'
    'd88ac00000000')

message = BitcacheMessage(TX_BYTES, [])
message.key_data.append(BitcacheTxoKeyUsage(0, "p2pkh", bytes.fromhex("9341cb4c"), "bip32:m/0/0"))
stream = io.BytesIO()
message_bytes = write_bitcache_transaction_message(stream, message)
assert stream.getvalue() == bytes.fromhex('545842595445532e58000000000000000200000001000000000000'
    '0000000000000000000000000000000000000000000000000000ffffffff03510101ffffffff0100f2052a010000'
    '001976a914e7d468699f0f2acbca887376739551053866617d88ac0000000054584f4b4559532e1e000000000000'
    '000000000005007032706b6804009341cb4c0b0062697033323a6d2f302f30')
```

## Development instructions

### Setup

Windows:
```cmd
py -3.10 -m pip install pipenv
py -3.10 -m pipenv --rm
py -3.10 -m pipenv --python 3.10
py -3.10 -m pipenv install -r requirements-dev.txt
```

Linux/MacOS:
```bash
python3.10 -m pip install pipenv
python3.10 -m pipenv --rm
python3.10 -m pipenv --python 3.10
python3.10 -m pipenv install -r requirements-dev.txt
```

### Running tests

Windows:
```cmd
py -3.10 -m pipenv run pytest bitcache\tests
```

Linux/MacOS:
```bash
python3.10 -m pipenv run pytest bitcache/tests
```
