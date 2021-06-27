package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	flag "github.com/spf13/pflag"

	// "strings"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StylusFrost/policy/x/policy/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the policy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdQueryCode(),
		GetCmdListCode(),
	)
	return queryCmd
}

// GetCmdQueryCode returns the bytecode for a given contract
func GetCmdQueryCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rego [rego_id] [output filename]",
		Short: "Downloads rego bytecode for given rego id",
		Long:  "Downloads rego bytecode for given rego id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			regoID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Rego(
				context.Background(),
				&types.QueryRegoRequest{
					RegoId: regoID,
				},
			)
			if err != nil {
				return err
			}
			if len(res.Data) == 0 {
				return fmt.Errorf("rego code not found")
			}

			fmt.Printf("Downloading rego code to %s\n", args[1])
			return ioutil.WriteFile(args[1], res.Data, 0644)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdListCode lists all wasm code uploaded
func GetCmdListCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-rego",
		Short: "List all rego bytecode on the chain",
		Long:  "List all rego bytecode on the chain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(withPageKeyDecoded(cmd.Flags()))
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Regos(
				context.Background(),
				&types.QueryRegosRequest{
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "list rego codes")
	return cmd
}

// sdk ReadPageRequest expects binary but we encoded to base64 in our marshaller
func withPageKeyDecoded(flagSet *flag.FlagSet) *flag.FlagSet {
	encoded, err := flagSet.GetString(flags.FlagPageKey)
	if err != nil {
		panic(err.Error())
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err.Error())
	}
	flagSet.Set(flags.FlagPageKey, string(raw))
	return flagSet
}
