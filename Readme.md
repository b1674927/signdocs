# Signdocs

A command line program to cryptographically sign documents and recover the address from a signature and hash.

The app creates a Keccak256 hash of the document and creates an ECDSA signature
`\x19Ethereum Signed Message:\n` is _not_ prefixed to the data.

## Sign

```
$ signdocs sign mydoc.pdf
```

You enter your private-key as prompted and signdocs will show you the hash of 'mydoc.pdf' and your signature of the hash.

## Recover

Type

```
$ signdocs recover
```

You enter the hash and signature as prompted and signdocs will show you the signer's address.

## Help

```
$ signdocs help
NAME:
   signdocs - ECDSA signature tool

USAGE:
   signdocs [global options] command [command options]

COMMANDS:
   sign     sign a file
   recover  recover an address from signature and hash
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

# Install using make

- clone this repo
- ensure go version >=1.23.3 is installed
  ```
  $ go version
  go version go1.23.3 linux/amd64
  ```
- run make
  ```bash
  make
  sudo cp ./build/signdocs /usr/local/bin
  ```
