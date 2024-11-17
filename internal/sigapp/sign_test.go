package sigapp

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

func loadPk() string {
	viper.SetConfigFile("../../.env")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("could not load .env file" + err.Error())
	}
	return viper.GetString("PK")
}

func TestDeriveAddress(t *testing.T) {
	pk := loadPk()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	pk2, err := crypto.HexToECDSA(strings.TrimSpace(pk))
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	addr, err := deriveAddress(pk2)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	fmt.Println(addr)
}
