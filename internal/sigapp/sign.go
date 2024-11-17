package sigapp

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

type SignModel struct {
	state      SignState
	privateKey *ecdsa.PrivateKey
	ethAddress string
	input      textinput.Model
	file       []byte
	fileName   string
	filehash   [32]byte
	signature  []byte
}

type SignState int

const (
	enterPrivateKey SignState = iota
	showSign
)

func SignCommand(cCtx *cli.Context) error {
	filename := cCtx.Args().Get(0)
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}
	p := tea.NewProgram(initialSignModel(filename, fileData))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error: %v\n", err)
	}
	return nil
}

func initialSignModel(filename string, file []byte) SignModel {

	m := SignModel{
		state:    enterPrivateKey,
		input:    textinput.New(),
		file:     file,
		fileName: filename,
	}
	m.input.CharLimit = 66
	m.input.Placeholder = "CaFe..."
	m.input.Focus()
	m.input.EchoMode = textinput.EchoPassword
	return m
}

func (m SignModel) Init() tea.Cmd {
	return textinput.Blink
}

// Derive the Ethereum address from a private key
func deriveAddress(privateKey *ecdsa.PrivateKey) (string, error) {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKey).Hex(), nil
}

func (m SignModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case enterPrivateKey:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				pk := m.input.Value()
				pk = strings.TrimSpace(pk)
				pk = strings.TrimPrefix(pk, "0x")
				var err error
				m.privateKey, err = crypto.HexToECDSA(strings.TrimSpace(pk))
				if err != nil {
					m.input.SetValue("")
					m.input.Placeholder = "Invalid key, try again"
					return m, nil
				}
				m.ethAddress, err = deriveAddress(m.privateKey)
				if err != nil {
					m.input.SetValue("")
					m.input.Placeholder = "Invalid key, try again"
					return m, nil
				}
				err = m.processFile()
				if err != nil {
					m.input.SetValue("")
					m.input.Placeholder = "Error:" + err.Error()
					return m, nil
				}
				m.state = showSign
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case showSign:
		return m, tea.Quit
	}
	return m, nil
}

func (m SignModel) View() string {
	switch m.state {
	case enterPrivateKey:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Render("Enter your Ethereum private key:") +
			m.input.View()
	case showSign:
		msg := fmt.Sprintf("File processed: %s\n", m.fileName)
		msg = msg + fmt.Sprintf("Your signer address is:\n%s\n", m.ethAddress)
		msg = msg + fmt.Sprintf("The SHA-256 hash of your file is:\n%x\n", m.filehash)
		msg = msg + fmt.Sprintf("Your signature of the hash is:\n%x\n", m.signature)
		return msg
	}
	return ""
}

// Dummy file processing function
func (m *SignModel) processFile() error {
	hash := sha256.Sum256(m.file)
	m.filehash = hash
	// Sign the hash
	var err error
	m.signature, err = crypto.Sign(hash[:], m.privateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign the hash: %v", err)
	}
	return nil
}
