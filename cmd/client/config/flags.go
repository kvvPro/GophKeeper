package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ClientFlags struct {
	Address   string `env:"ADDRESS" json:"address"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config    string `env:"CONFIG" json:"config"`
}

func Initialize(agentFlags *ClientFlags) error {

	// try to get vars from Flags
	pflag.StringVarP(&agentFlags.Address, "addr", "a", "localhost:8080", "Net address host:port")
	pflag.StringVarP(&agentFlags.CryptoKey, "crypto-key", "e", "/workspaces/gophkeeper/cmd/keys/key.pub", "Path to public key RSA to encrypt messages")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	fmt.Println("\nFLAGS-----------")
	fmt.Printf("ADDRESS=%v", agentFlags.Address)
	fmt.Printf("\nCRYPTO_KEY=%v", agentFlags.CryptoKey)
	fmt.Printf("\nCONFIG=%v", agentFlags.Config)
	fmt.Println()

	// try to get vars from env
	if err := env.Parse(agentFlags); err != nil {
		return err
	}

	fmt.Println("ENV-----------")
	fmt.Printf("ADDRESS=%v", agentFlags.Address)
	fmt.Printf("\nCRYPTO_KEY=%v", agentFlags.CryptoKey)
	fmt.Printf("\nCONFIG=%v", agentFlags.Config)

	return nil
}

func ReadConfig() (*ClientFlags, error) {
	flags := new(ClientFlags)

	pflag.StringVarP(&flags.Config, "config", "c", "/workspaces/gophkeeper/cmd/server/config/config.json", "Path to server config file")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	fmt.Printf("CONFIG=%v", flags.Config)

	data, err := os.ReadFile(flags.Config)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(&flags); err != nil {
		return nil, err
	}

	return flags, nil
}
