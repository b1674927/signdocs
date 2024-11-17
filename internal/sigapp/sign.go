package sigapp

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
	fileOut    string
	filehash   []byte
	signature  []byte
	descr      string
}

type SignState int

const (
	enterPrivateKey SignState = iota
	enterDescription
	showSign
)

func SignCommand(cCtx *cli.Context) error {
	filename := cCtx.Args().Get(0)
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}
	fileOut := cCtx.String("file")
	p := tea.NewProgram(initialSignModel(filename, fileData, fileOut))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error: %v\n", err)
	}
	return nil
}

func initialSignModel(filename string, file []byte, fileOut string) SignModel {

	m := SignModel{
		state:    enterPrivateKey,
		input:    textinput.New(),
		file:     file,
		fileName: filename,
		fileOut:  fileOut,
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
				m.input.SetValue("")
				if err != nil {
					m.input.Placeholder = "Invalid key, try again"
					return m, nil
				}
				m.ethAddress, err = deriveAddress(m.privateKey)
				if err != nil {
					m.input.Placeholder = "Invalid key, try again"
					return m, nil
				}
				err = m.processFile()
				if err != nil {
					m.input.Placeholder = "Error:" + err.Error()
					return m, nil
				}
				m.input.EchoMode = textinput.EchoNormal
				m.input.Placeholder = "signdocs signed document"
				m.state = enterDescription
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case enterDescription:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.descr = m.input.Value()
				m.state = showSign
				return m, nil
			case "ctrl+c":
				m.state = showSign
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case showSign:
		switch msg.(type) {
		case tea.KeyMsg:
			return m, tea.Quit
		}
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
	case enterDescription:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Render("Enter description of the file:") +
			m.input.View()
	case showSign:
		txtJson := m.ToDocSignature()
		var msg string
		if m.fileOut != "" {
			err := writeFile(txtJson, m.fileOut)
			if err != nil {
				msg = fmt.Sprintf("failed to write outputfile: %v", err)
			} else {
				msg = fmt.Sprintf("result written to: %s", m.fileOut)
			}
		} else {
			msg = string(txtJson)
		}
		return string(msg)
	}
	return ""
}

func writeFile(txtJson []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer file.Close()
	_, err = file.Write(txtJson)
	if err != nil {
		return fmt.Errorf("Failed to write JSON to file: %v", err)
	}
	return nil
}

// Dummy file processing function
func (m *SignModel) processFile() error {
	hash := crypto.Keccak256(m.file)
	m.filehash = hash
	// Sign the hash
	var err error
	m.signature, err = crypto.Sign(hash[:], m.privateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign the hash: %v", err)
	}
	return nil
}

func (m *SignModel) ToDocSignature() []byte {
	docSig := DocumentSignature{
		Metadata: Metadata{
			Name:        m.fileName,
			Description: m.descr,
			Timestamp:   time.Now().UTC().Format(time.RFC3339), // ISO 8601
		},
		FileHash:  fmt.Sprintf("0x%x", m.filehash),
		Signature: fmt.Sprintf("0x%x", m.signature),
		Signer:    fmt.Sprintf("%s", m.ethAddress),
	}
	// Convert struct to JSON
	jsonData, _ := json.MarshalIndent(docSig, "", "  ")
	return jsonData
}
