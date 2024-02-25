package cli

import (
	"strconv"

	"auction/x/auction/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-auction [name] [start-price] [duration] [created-at] [current-highest-bid-id] [highest-bid-exists] [ended]",
		Short: "Create a new auction",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]
			argStartPrice, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			argDuration, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argCreatedAt, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argCurrentHighestBidId, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}
			argHighestBidExists, err := cast.ToBoolE(args[5])
			if err != nil {
				return err
			}
			argEnded, err := cast.ToBoolE(args[6])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAuction(clientCtx.GetFromAddress().String(), argName, argStartPrice, argDuration, argCreatedAt, argCurrentHighestBidId, argHighestBidExists, argEnded)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-auction [id] [name] [start-price] [duration] [created-at] [current-highest-bid-id] [highest-bid-exists] [ended]",
		Short: "Update a auction",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argName := args[1]

			argStartPrice, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			argDuration, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argCreatedAt, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argCurrentHighestBidId, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			argHighestBidExists, err := cast.ToBoolE(args[6])
			if err != nil {
				return err
			}

			argEnded, err := cast.ToBoolE(args[7])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAuction(clientCtx.GetFromAddress().String(), id, argName, argStartPrice, argDuration, argCreatedAt, argCurrentHighestBidId, argHighestBidExists, argEnded)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-auction [id]",
		Short: "Delete a auction by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteAuction(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
