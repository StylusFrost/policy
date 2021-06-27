package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	regoUtils "github.com/StylusFrost/policy/x/policy/client/utils"

	"github.com/StylusFrost/policy/x/policy/types"
	flag "github.com/spf13/pflag"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	flagSource                 = "source"
	flagInstantiateByEverybody = "instantiate-everybody"
	flagInstantiateByAddress   = "instantiate-only-address"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Policy transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		StoreRegoCmd(),
	)
	return txCmd
}

// StoreRegoCmd will upload code to be reused.
func StoreRegoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store [rego_file] [json_encoded_entry_points] --source [source] ",
		Short: "Upload a Rego Policy code",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			msg, err := parseStoreRegoArgs(args[0], args[1], clientCtx.GetFromAddress(), cmd.Flags())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(flagSource, "", "A valid URI reference to the policy's source code, optional")
	cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a policy from the rego code, optional")
	cmd.Flags().String(flagInstantiateByAddress, "", "Only this address can instantiate a policy instance from the rego code, optional")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseStoreRegoArgs(file string, entry_points string, sender sdk.AccAddress, flags *flag.FlagSet) (types.MsgStoreRego, error) {
	rego, err := ioutil.ReadFile(file)
	if err != nil {
		return types.MsgStoreRego{}, err
	}

	// gzip the rego file
	if !regoUtils.IsGzip(rego) {
		rego, err = regoUtils.GzipIt(rego)

		if err != nil {
			return types.MsgStoreRego{}, err
		}
	}

	var perm *types.AccessConfig
	onlyAddrStr, err := flags.GetString(flagInstantiateByAddress)
	if err != nil {
		return types.MsgStoreRego{}, fmt.Errorf("instantiate by address: %s", err)
	}
	if onlyAddrStr != "" {
		allowedAddr, err := sdk.AccAddressFromBech32(onlyAddrStr)
		if err != nil {
			return types.MsgStoreRego{}, sdkerrors.Wrap(err, flagInstantiateByAddress)
		}
		x := types.AccessTypeOnlyAddress.With(allowedAddr)
		perm = &x
	} else {
		everybodyStr, err := flags.GetString(flagInstantiateByEverybody)
		if err != nil {
			return types.MsgStoreRego{}, fmt.Errorf("instantiate by everybody: %s", err)
		}
		if everybodyStr != "" {
			ok, err := strconv.ParseBool(everybodyStr)
			if err != nil {
				return types.MsgStoreRego{}, fmt.Errorf("boolean value expected for instantiate by everybody: %s", err)
			}
			if ok {
				perm = &types.AllowEverybody
			}
		}
	}

	// build and sign the transaction, then broadcast to Tendermint
	source, err := flags.GetString(flagSource)
	if err != nil {
		return types.MsgStoreRego{}, fmt.Errorf("source: %s", err)
	}

	if !json.Valid([]byte(entry_points)) {
		return types.MsgStoreRego{}, fmt.Errorf("entry_points: invalid json")
	}

	var entryPointsArr []string
	errJson := json.Unmarshal([]byte(entry_points), &entryPointsArr)

	if errJson != nil {
		return types.MsgStoreRego{}, fmt.Errorf("entry_points: invalid json array")
	}

	msg := types.MsgStoreRego{
		Sender:                sender.String(),
		REGOByteCode:          rego,
		Source:                source,
		EntryPoints:           []byte(entry_points),
		InstantiatePermission: perm,
	}
	return msg, nil
}
