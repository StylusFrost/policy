package cli

import (
	"encoding/json"
	"errors"
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
	flagAmount                 = "amount"
	flagInstantiateByEverybody = "instantiate-everybody"
	flagInstantiateByAddress   = "instantiate-only-address"
	flagLabel                  = "label"
	flagAdmin                  = "admin"
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
		InstantiatePolicyCmd(),
		UpdatePolicyAdminCmd(),
		ClearPolicyAdminCmd(),
		MigratePolicyCmd(),
	)
	return txCmd
}

// StoreRegoCmd will upload rego to be reused.
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

// InstantiatePolicyCmd will instantiate a policy from previously uploaded rego code.
func InstantiatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate [rego_id_int64] [json_encoded_entry_points] --label [text] --admin [address,optional] --amount [coins,optional]",
		Short: "Instantiate a rego policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)

			msg, err := parseInstantiateArgs(args[0], args[1], clientCtx.GetFromAddress(), cmd.Flags())
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(flagAmount, "", "Coins to send to the policy during instantiation")
	cmd.Flags().String(flagLabel, "", "A human-readable name for this policy in lists")
	cmd.Flags().String(flagAdmin, "", "Address of an admin")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseInstantiateArgs(rawRegoID, entry_points string, sender sdk.AccAddress, flags *flag.FlagSet) (types.MsgInstantiatePolicy, error) {
	// get the id of the rego to instantiate
	regoID, err := strconv.ParseUint(rawRegoID, 10, 64)
	if err != nil {
		return types.MsgInstantiatePolicy{}, err
	}

	amountStr, err := flags.GetString(flagAmount)
	if err != nil {
		return types.MsgInstantiatePolicy{}, fmt.Errorf("amount: %s", err)
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return types.MsgInstantiatePolicy{}, fmt.Errorf("amount: %s", err)
	}
	label, err := flags.GetString(flagLabel)
	if err != nil {
		return types.MsgInstantiatePolicy{}, fmt.Errorf("label: %s", err)
	}
	if label == "" {
		return types.MsgInstantiatePolicy{}, errors.New("label is required on all policys")
	}
	adminStr, err := flags.GetString(flagAdmin)
	if err != nil {
		return types.MsgInstantiatePolicy{}, fmt.Errorf("admin: %s", err)
	}

	var entryPointsArr []types.EntryPoint
	errJson := json.Unmarshal([]byte(entry_points), &entryPointsArr)

	if errJson != nil {
		return types.MsgInstantiatePolicy{}, fmt.Errorf("entry_points: invalid json array")
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.MsgInstantiatePolicy{
		Sender:      sender.String(),
		RegoID:      regoID,
		Label:       label,
		Funds:       amount,
		EntryPoints: []byte(entry_points),
		Admin:       adminStr,
	}
	return msg, nil
}

// UpdatePolicyAdminCmd sets an new admin for a policy
func UpdatePolicyAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-policy-admin [policy_addr_bech32] [new_admin_addr_bech32]",
		Short: "Set new admin for a policy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg, err := parseUpdatePolicyAdminArgs(args, clientCtx)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseUpdatePolicyAdminArgs(args []string, cliCtx client.Context) (types.MsgUpdateAdmin, error) {
	msg := types.MsgUpdateAdmin{
		Sender:   cliCtx.GetFromAddress().String(),
		Policy:   args[0],
		NewAdmin: args[1],
	}
	return msg, nil
}

// ClearPolicyAdminCmd clears an admin for a policy
func ClearPolicyAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-policy-admin [policy_addr_bech32]",
		Short: "Clears admin for a policy to prevent further migrations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgClearAdmin{
				Sender: clientCtx.GetFromAddress().String(),
				Policy: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// MigratePolicyCmd will migrate a policy to a new code version
func MigratePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate [policy_addr_bech32] [new_rego_id_int64] [json_encoded_migration_args]",
		Short:   "Migrate a rego policy to a new code version",
		Aliases: []string{"update", "mig", "m"},
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg, err := parseMigratePolicyArgs(args, clientCtx)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return nil
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseMigratePolicyArgs(args []string, cliCtx client.Context) (types.MsgMigratePolicy, error) {
	// get the id of the code to instantiate
	regoID, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return types.MsgMigratePolicy{}, sdkerrors.Wrap(err, "rego id")
	}

	entry_points := args[2]

	var entryPointsArr []types.EntryPoint
	errJson := json.Unmarshal([]byte(entry_points), &entryPointsArr)

	if errJson != nil {
		return types.MsgMigratePolicy{}, fmt.Errorf("entry_points: invalid json array")
	}

	msg := types.MsgMigratePolicy{
		Sender:      cliCtx.GetFromAddress().String(),
		Policy:      args[0],
		RegoID:      regoID,
		EntryPoints: []byte(entry_points),
	}
	return msg, nil
}
