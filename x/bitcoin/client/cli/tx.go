package cli

import (
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authUtils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	balance "github.com/axelarnetwork/axelar-core/x/balance/exported"

	"github.com/spf13/cobra"

	"github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/x/bitcoin/keeper"
	"github.com/axelarnetwork/axelar-core/x/bitcoin/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	btcTxCmd := &cobra.Command{
		Use:                        "bitcoin",
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		TraverseChildren:           true,
		RunE:                       client.ValidateCmd,
	}

	btcTxCmd.AddCommand(flags.PostCommands(
		GetCmdTrackAddress(cdc),
		GetCmdVerifyTx(cdc),
		GetCmdSignRawTx(cdc),
		GetCmdLink(cdc),
	)...)

	return btcTxCmd
}

func GetCmdTrackAddress(cdc *codec.Codec) *cobra.Command {
	var rescan bool
	addrCmd := &cobra.Command{
		Use:   "track [address]",
		Short: "Make the axelar network aware of a specific address on Bitcoin",
		Long:  "Make the axelar network aware of a specific address on Bitcoin. Use --rescan to rescan the entire Bitcoin history for past transactions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, txBldr := utils.PrepareCli(cmd.InOrStdin(), cdc)

			msg := types.NewMsgTrackAddress(cliCtx.GetFromAddress(), args[0], rescan)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authUtils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	addRescanFlag(addrCmd, &rescan)
	return addrCmd
}

// GetCmdVerifyTx returns the transaction verification command
func GetCmdVerifyTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "verifyTx [txInfo json]",
		Short: "Verify a Bitcoin transaction",
		Long: fmt.Sprintf(
			"Verify that a transaction happened on the Bitcoin network so it can be processed on axelar. "+
				"Get the json string by using the %s query", keeper.QueryOutInfo),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, txBldr := utils.PrepareCli(cmd.InOrStdin(), cdc)

			var out types.OutPointInfo
			cliCtx.Codec.MustUnmarshalJSON([]byte(args[0]), &out)

			msg := types.MsgVerifyTx{Sender: cliCtx.GetFromAddress(), OutPointInfo: out}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authUtils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSignRawTx returns the command to sign a raw Bitcoin transaction
func GetCmdSignRawTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "signTx [txID:voutIdx] [tx json]",
		Short: "Sign raw spending transaction with utxo of [txID]",
		Long:  fmt.Sprintf("Sign raw transaction. Get raw transaction by querying %s", keeper.QueryRawTx),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, txBldr := utils.PrepareCli(cmd.InOrStdin(), cdc)
			var tx *wire.MsgTx
			types.ModuleCdc.MustUnmarshalJSON([]byte(args[1]), &tx)

			outpoint, err := types.OutPointFromStr(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSignTx(cliCtx.GetFromAddress(), outpoint, tx)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authUtils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLink links a cross chain address to a bitcoin address created by Axelar
func GetCmdLink(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "link [chain] [address]",
		Short: "Link a cross chain address to a bitcoin address created by Axelar",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, txBldr := utils.PrepareCli(cmd.InOrStdin(), cdc)

			chain := balance.ChainFromString(args[0])
			address := balance.CrossChainAddress{Chain: chain, Address: args[1]}

			if err := address.Validate(); err != nil {
				return err
			}

			msg := types.MsgLink{Sender: cliCtx.GetFromAddress(), Recipient: address}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authUtils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func addRescanFlag(cmd *cobra.Command, rescan *bool) {
	cmd.Flags().BoolVarP(rescan, "rescan", "r", false,
		"Rescan the entire Bitcoin blockchain for previous transactions to this address")
}