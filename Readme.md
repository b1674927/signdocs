# Signdocs

A command line program to (1) cryptographically sign documents and (2) verify signatures

- You need an Ethereum private key to sign
- The app creates a Keccak256 hash of the document you provide and creates an ECDSA signature using your private key
- Alternatively, the app recovers the address from a signature and hash, i.e., verifies the signature
- Keys that you enter are not stored and not visible in the console
- `\x19Ethereum Signed Message:\n` is _not_ prefixed to the data

## Sign

```
$ signdocs sign mydoc.pdf
```

You enter your private-key as prompted and signdocs will show you the hash of 'mydoc.pdf' and your signature of the hash.

You can create a metadata file with the signature of the document and its hash with the option --file or -f (see signdocs sign help)
$ signdocs sign -f outfile.json mydoc.pdf
Example:

```
{
  "metadata": {
    "name": "document.pdf",
    "description": "my document available at source",
    "timestamp": "2024-11-17T21:47:14Z"
  },
  "fileHash": "0xdc07ae4dd8c0975e27284314d4f078efc3567b5bd482786c5db3cb550fc9055a",
  "signature": "0xff457afdfe05fe84d5380bf573d84d84baaae3b30980974b35c05ed9f9d26db42bfcf8853bd8a768af80620f1f8f5c27498b8d4e90dcf086c4f9bae53029d89300",
  "signer": "0x7Fee781C115AB16a0065E5023ee63f4950F2F06c"
}
```

Note that the metadata is not signed, only the hash.

## Recover

```
$ signdocs verify
```

You enter the hash and signature as prompted and signdocs will show you the signer's address.

## Help

- `$ signdocs help`
- `$ signdocs sign help`
- `$ signdocs verify help`

```
$ signdocs help
NAME:
   signdocs - ECDSA signature tool

USAGE:
   signdocs [global options] command [command options]

COMMANDS:
   sign     sign a file
   verify   recover an address from signature and hash
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

# Install using make

- clone this repo
- building the app requires go, see https://go.dev/doc/install
- ensure go version >=1.23.3 is installed

  ```
  $ go version
  go version go1.23.3 linux/amd64
  ```

- run make in the folder you cloned
  ```bash
  $cd signdocs
  $make
  ```
  Copy the executable into folder in your executable path, for example (Linux):
  ```
  sudo cp ./build/signdocs /usr/local/bin
  ```
