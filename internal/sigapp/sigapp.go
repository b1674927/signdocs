package sigapp

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type SubCommand int

const (
	SUB_HELP SubCommand = iota
	SUB_SIGN
	SUB_RECOVER
)

type DocumentSignature struct {
	Metadata  Metadata `json:"metadata"`
	FileHash  string   `json:"fileHash"`  // The keccak256 hash of the document
	Signature string   `json:"signature"` // ECDSA signature of the file hash
	Signer    string   `json:"signer"`    // Ethereum address of the signer
}

type Metadata struct {
	Name        string `json:"name"`        // File name
	Description string `json:"description"` // Description of the file
	Timestamp   string `json:"timestamp"`   // ISO 8601 timestamp
}

func Run(args []string) {
	app := &cli.App{
		Name:  "signdocs",
		Usage: "ECDSA signature tool",
		CommandNotFound: func(ctx *cli.Context, s string) {
			fmt.Printf("Unknown command %s\n", s)
		},
		Commands: []*cli.Command{

			{
				Name:      "sign",
				Usage:     "sign a file",
				UsageText: "signdocs sign [document.pdf] - Signs the specified document\nsigndocs sign --file out.json [document.pdf] - Signs the specified document and writes the output to out.json",
				Action:    SignCommand,
				Flags: []cli.Flag{&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Value:   "",
					Usage:   "output file",
				}},
			},
			{
				Name:      "recover",
				Usage:     "recover an address from signature and hash",
				UsageText: "signdocs recover",
				Action:    RecoverCommand,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
