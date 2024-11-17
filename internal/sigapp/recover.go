package sigapp

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

type RecoverModel struct {
	state      RecoverState
	ethAddress string
	input      textinput.Model
	filehash   [32]byte
	signature  []byte
}

type RecoverState int

const (
	enterHash RecoverState = iota
	enterSig
	showRecover
)

func RecoverCommand(cCtx *cli.Context) error {
	p := tea.NewProgram(initialRecoverModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Error: %v\n", err)
	}
	return nil
}

func (m RecoverModel) Init() tea.Cmd {
	return textinput.Blink
}

func initialRecoverModel() RecoverModel {

	m := RecoverModel{
		state: enterHash,
		input: textinput.New(),
	}
	m.input.Focus()
	return m
}

func (m RecoverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case enterHash:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				hash := strings.TrimSpace(m.input.Value())
				hash = strings.TrimPrefix(hash, "0x")
				var err error
				m.filehash, err = stringToHash(hash)
				m.input.SetValue("")
				if err != nil {
					m.input.Placeholder = "Invalid hash, try again"
					return m, nil
				}
				m.state = enterSig
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case enterSig:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				sgn := strings.TrimSpace(m.input.Value())
				sgn = strings.TrimPrefix(sgn, "0x")
				var err error
				m.signature, err = hex.DecodeString(sgn)
				if err != nil {
					m.input.SetValue("")
					m.input.Placeholder = "Invalid signature, try again"
					return m, nil
				}
				m.recover()
				m.state = showRecover
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case showRecover:
		return m, tea.Quit
	}
	return m, nil
}

func stringToHash(input string) ([32]byte, error) {
	var hash [32]byte

	// Decode the hex string into a byte slice
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return hash, fmt.Errorf("invalid hex string: %w", err)
	}

	// Ensure the length is exactly 32 bytes
	if len(bytes) != 32 {
		return hash, fmt.Errorf("hash must be 32 bytes, got %d bytes", len(bytes))
	}

	// Copy the bytes into the fixed array
	copy(hash[:], bytes)

	return hash, nil
}

func (m RecoverModel) View() string {
	switch m.state {
	case enterHash:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Render("Enter the hash that has been signed:") +
			m.input.View()
	case enterSig:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Render("Enter the signature of the hash:") +
			m.input.View()
	case showRecover:
		msg := fmt.Sprintf("Recovered address: %s\n", m.ethAddress)
		return msg
	}
	return ""
}

func (m *RecoverModel) recover() {
	hash := m.filehash
	// Recover the public key from the signature
	publicKey, err := crypto.SigToPub(hash[:], m.signature)
	if err != nil {
		log.Fatalf("Failed to recover public key: %v", err)
	}

	// Get the Ethereum address from the public key
	m.ethAddress = crypto.PubkeyToAddress(*publicKey).Hex()
}
