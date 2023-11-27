import io, json, os
from typing import cast, TypedDict

from bitcache import read_bitcache_message, write_bitcache_transaction_message

def locate_top_path(path: str) -> str:
    return locate_top_path(os.path.join(path, "..")) \
        if not os.path.exists(os.path.join(path, "data")) else path
TOP_PATH = os.path.abspath(locate_top_path(os.path.dirname(os.path.abspath(__file__))))

DATA_PATH = os.path.join(TOP_PATH, "data")
assert os.path.exists(DATA_PATH), "Unable to find test input metadata"

MESSAGE_PATH = os.path.join(DATA_PATH, "bitcache-messages")
assert os.path.exists(DATA_PATH), "Unable to find test message data"

class JSONTxoKeyUsage(TypedDict):
    vout: int
    script_type: str
    key_fingerprint: str
    key_derivation: str

def test_read_bitcache_messages() -> None:
    assert os.path.exists(MESSAGE_PATH)
    for filename in sorted(os.listdir(MESSAGE_PATH)):
        message_path = os.path.join(MESSAGE_PATH, filename)
        with open(message_path, "rb") as f:
            message_data = f.read()
            f.seek(0, os.SEEK_SET)
            message = read_bitcache_message(f)
        # Verify that the message has correct transaction data.
        filename_prefix, _filename_suffix = os.path.splitext(filename)
        transaction_path = os.path.join(DATA_PATH, filename_prefix+".txn")
        assert os.path.exists(transaction_path)
        assert message.tx_data == open(transaction_path, "rb").read()
        # Verify that the message has correct metadata.
        metadata_path = os.path.join(DATA_PATH, filename_prefix+".json")
        assert os.path.exists(metadata_path)
        with open(metadata_path, "r") as f:
            for i, o in enumerate(cast(list[JSONTxoKeyUsage], json.load(f))):
                k = message.key_data[i]
                assert o["vout"] == k.txo_index
                assert o["script_type"] == k.script_type
                assert bytes.fromhex(o["key_fingerprint"]) == k.parent_key_fingerprint
                assert o["key_derivation"] == k.derivation_text
        # Verify that reconstructing the message gives matching bytes.
        output_stream = io.BytesIO()
        write_bitcache_transaction_message(output_stream, message)
        assert message_data == output_stream.getbuffer()
