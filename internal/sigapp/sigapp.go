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
				UsageText: "signdocs sign [document.pdf] - Signs the specified document",
				Action:    SignCommand,
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

// 	subCmd := SUB_HELP
// 	if len(args) > 0 {
// 		subCmd = strToSubCommand(args[1])
// 	}
// 	switch subCmd {
// 	case SUB_HELP:
// 		help()
// 		return
// 	case SUB_SIGN:
// 		if len(args) < 2 {
// 			help()
// 			return
// 		}
// 		filename := args[2]
// 		fileData, err := os.ReadFile(filename)
// 		if err != nil {
// 			log.Fatalf("Failed to read file: %v", err)
// 			os.Exit(1)
// 		}
// 		p := tea.NewProgram(initialModel(filename, fileData))
// 		if _, err := p.Run(); err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			os.Exit(1)
// 		}
// 	}

// }
