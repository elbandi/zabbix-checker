package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/guarilha/go-gnosispay"

	"github.com/urfave/cli/v2"
)

var (
	addressFlag = cli.StringFlag{
		Name:     "address",
		Usage:    "Address",
		Required: true,
		EnvVars:  []string{"ADDRESS"},
		Action: func(ctx *cli.Context, v string) error {
			if len(v) == 0 {
				return cli.Exit("Flag 'address' cannot be empty", 1)
			}
			return nil
		},
	}
	privateKeyFlag = cli.StringFlag{
		Name:    "private_key",
		Usage:   "Private key",
		EnvVars: []string{"PRIVATE_KEY"},
		Action: func(ctx *cli.Context, v string) error {
			if len(v) == 0 {
				return cli.Exit("Flag 'private_key' cannot be empty", 1)
			}
			return nil
		},
	}
	privateKeyFileFlag = cli.StringFlag{
		Name:    "private_key_file",
		Usage:   "Private key file",
		EnvVars: []string{"PRIVATE_KEY_FILE"},
		Action: func(ctx *cli.Context, v string) error {
			if len(v) == 0 {
				return cli.Exit("Flag 'private_key_file' cannot be empty", 1)
			}
			return nil
		},
	}
)

var walletCommand = cli.Command{
	Name:  "wallet",
	Usage: "wallet data",
	Flags: []cli.Flag{
		&addressFlag,
		&privateKeyFlag,
		&privateKeyFileFlag,
	},
	Action: cmdWallet,
}

func cmdWallet(ctx *cli.Context) error {
	client, err := gnosispay.New(nil,
		gnosispay.SetBaseURL("https://api.gnosispay.com"),
		gnosispay.SetSIWEParams("https://app.gnosispay.com"),
	)
	if err != nil {
		return err
	}

	var pk *ecdsa.PrivateKey
	if ctx.IsSet(privateKeyFlag.Name) {
		pk, err = crypto.HexToECDSA(ctx.String(privateKeyFlag.Name))
	} else if ctx.IsSet(privateKeyFileFlag.Name) {
		pk, err = crypto.LoadECDSA(ctx.String(privateKeyFileFlag.Name))
	} else {
		return cli.Exit("Either 'private_key' or 'private_key_file' must be specified", 1)
	}

	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	address := common.HexToAddress(ctx.String(addressFlag.Name))

	// Auth
	_, err = client.Auth.AuthenticateWithPrivateKey(ctx.Context, address, pk)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// User
	user, err := client.User.Get(ctx.Context)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	// Balance
	balance, err := client.Account.GetBalances(ctx.Context)
	if err != nil {
		return fmt.Errorf("failed to get balances: %v", err)
	}
	// Trim last 10 digits from balance.Total and balance.Spendable
	balance.Total = balance.Total[:len(balance.Total)-10]
	balance.Spendable = balance.Spendable[:len(balance.Spendable)-10]
	totalBalance, err := strconv.ParseInt(balance.Total, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse total balance: %v", err)
	}
	spendableBalance, err := strconv.ParseInt(balance.Spendable, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse spendable balance: %v", err)
	}

	data, err := json.Marshal(struct {
		SiwCount         int     `json:"siw_count"`
		SwCount          int     `json:"sw_count"`
		KycStatus        string  `json:"kyc_status"`
		CardCount        int     `json:"card_count"`
		TotalBalance     float64 `json:"total_balance"`
		SpendableBalance float64 `json:"spendable_balance"`
	}{
		SiwCount:         len(user.SignInWallets),
		SwCount:          len(user.SafeWallets),
		KycStatus:        string(*user.KycStatus),
		CardCount:        len(user.Cards),
		TotalBalance:     float64(totalBalance) / 1e8,
		SpendableBalance: float64(spendableBalance) / 1e8,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %v", err)
	}
	//	log.Printf("User: %+v, %s\n", user, string(data))
	fmt.Println(string(data))

	return nil
}
