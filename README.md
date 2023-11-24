# BitCache Message

`bitcache-message` is a Go package for serializing and deserializing BitCache messages to and from bytes. This library, part of the `bcm` (BitCache Message) package, provides a simple and efficient way to handle BitCache messages, making it easier to work with them in Go applications.

## Installation

To install the `bitcache-message` package, use the `go get` command:

```bash
go get github.com/bitcache-tech/bitcache-message
```

This will download and install the bitcache-message package, ready for use in your Go projects.

## Usage

Here is a basic example of how to use the bitcache-message package to encode and decode BitCache messages.

### Encoding a BitCache Message

```go
package main

import (
    "github.com/bitcache-tech/bitcache-message/bcm"
    "log"
)

func main() {
    // Create a BitCacheMessage instance
    bcmInstance := bcm.BitCacheMessage{
        // Initialize your BitCacheMessage here
    }

    // Encode the BitCacheMessage to bytes
    encodedBytes, err := bcmInstance.Bytes()
    if err != nil {
        log.Fatalf("Error encoding BitCacheMessage: %v", err)
    }

    // Use encodedBytes as needed
}
```

### Decoding a BitCache Message

```go
package main

import (
    "github.com/bitcache-tech/bitcache-message/bcm"
    "log"
)

func main() {
    // Assume you have a byte slice `data` that represents a serialized BitCacheMessage
    data := []byte{ /* ... */ }

    // Decode the byte slice into a BitCacheMessage
    bcmInstance, err := bcm.NewBitCacheMessage(data)
    if err != nil {
        log.Fatalf("Error decoding BitCacheMessage: %v", err)
    }

    // Use bcmInstance as needed
}
```

## License

This library is distributed under the ISC License. See LICENSE file in the repository for more information.

