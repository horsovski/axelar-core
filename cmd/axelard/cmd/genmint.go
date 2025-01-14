package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/spf13/cobra"
)

const (
	flagInflationMax = "inflation-max"
	flagInflationMin = "inflation-min"
)

// SetGenesisMintCmd returns set-genesis-mint cobra Command.
func SetGenesisMintCmd(defaultNodeHome string) *cobra.Command {
	var (
		inflationMin string
		inflationMax string
	)

	cmd := &cobra.Command{
		Use:   "set-genesis-mint",
		Short: "Set the genesis parameters for the mint module",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			var genesis mintTypes.GenesisState
			if appState[mintTypes.ModuleName] != nil {
				cdc.MustUnmarshalJSON(appState[mintTypes.ModuleName], &genesis)
			}

			if inflationMin != "" {
				min, err := sdk.NewDecFromStr(inflationMin)
				if err != nil {
					return err
				}

				genesis.Params.InflationMin = min
			}

			if inflationMax != "" {
				max, err := sdk.NewDecFromStr(inflationMax)
				if err != nil {
					return err
				}

				genesis.Params.InflationMax = max
			}

			if err := genesis.Params.Validate(); err != nil {
				return err
			}

			genesisBz, err := cdc.MarshalJSON(&genesis)
			if err != nil {
				return fmt.Errorf("failed to marshal genesis state: %w", err)
			}
			appState[mintTypes.ModuleName] = genesisBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().StringVar(&inflationMin, flagInflationMin, "", "Minimum inflation rate")
	cmd.Flags().StringVar(&inflationMax, flagInflationMax, "", "Maximum inflation rate")

	return cmd
}
