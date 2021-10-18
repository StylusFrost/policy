package cli

import (
	"fmt"

	"github.com/StylusFrost/policy/x/policy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
)

func StoreRegoProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy-store [rego_file] [json_encoded_entry_points] --source [source] --title [text] --description [text] --run-as [address]",
		Short: "Submit a rego binary proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src, err := parseStoreRegoArgs(args[0], args[1], clientCtx.FromAddress, cmd.Flags())
			if err != nil {
				return err
			}
			runAs, err := cmd.Flags().GetString(flagRunAs)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(runAs) == 0 {
				return fmt.Errorf("run-as address is required")
			}
			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.StoreRegoProposal{
				Title:                 proposalTitle,
				Description:           proposalDescr,
				RunAs:                 runAs,
				REGOByteCode:          src.REGOByteCode,
				EntryPoints:           src.EntryPoints,
				Source:                src.Source,
				InstantiatePermission: src.InstantiatePermission,
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagRunAs, "", "The address that is stored as rego creator")
	cmd.Flags().String(flagSource, "", "A valid URI reference to the policy's source code, optional")
	cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a policy from the rego code, optional")
	cmd.Flags().String(flagInstantiateByAddress, "", "Only this address can instantiate a policy instance from the rego code, optional")

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")

	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: Permission of proposal, types: policy-store/migrate-policy/instantiate-policy/set-policy-admin/clear-policy-admin/text/parameter_change/community-pool-spend/software_upgrade/cancel-software-upgrade")
	return cmd
}

func ProposalInstantiatePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate-policy [rego_id_int64] [json_encoded_entry_points] --label [text] --title [text] --description [text] --run-as [address] --admin [address,optional] --amount [coins,optional]",
		Short: "Submit an instantiate rego policy proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src, err := parseInstantiateArgs(args[0], args[1], clientCtx.FromAddress, cmd.Flags())
			if err != nil {
				return err
			}

			runAs, err := cmd.Flags().GetString(flagRunAs)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(runAs) == 0 {
				return fmt.Errorf("run-as address is required")
			}
			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.InstantiatePolicyProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				RunAs:       runAs,
				Admin:       src.Admin,
				RegoID:      src.RegoID,
				Label:       src.Label,
				EntryPoints: src.EntryPoints,
				Funds:       src.Funds,
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(flagAmount, "", "Coins to send to the policy during instantiation")
	cmd.Flags().String(flagLabel, "", "A human-readable name for this policy in lists")
	cmd.Flags().String(flagAdmin, "", "Address of an admin")
	cmd.Flags().String(flagRunAs, "", "The address that pays the init funds. It is the creator of the policy and passed to the policy as sender on proposal execution")

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: Permission of proposal, types: policy-store/migrate-policy/instantiate-policy/set-policy-admin/clear-policy-admin/text/parameter_change/community-pool-spend/software_upgrade/cancel-software-upgrade")
	return cmd
}

func ProposalMigratePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-policy [policy_addr_bech32] [new_rego_id_int64] [json_encoded_migration_args]",
		Short: "Submit a migrate rego policy to a new code version proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src, err := parseMigratePolicyArgs(args, clientCtx)
			if err != nil {
				return err
			}

			runAs, err := cmd.Flags().GetString(flagRunAs)
			if err != nil {
				return fmt.Errorf("run-as: %s", err)
			}
			if len(runAs) == 0 {
				return fmt.Errorf("run-as address is required")
			}
			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.MigratePolicyProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Policy:      src.Policy,
				RegoID:      src.RegoID,
				EntryPoints: src.EntryPoints,
				RunAs:       runAs,
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(flagRunAs, "", "The address that is passed as sender to the policy on proposal execution")

	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: Permission of proposal, types: policy-store/migrate-policy/instantiate-policy/set-policy-admin/clear-policy-admin/text/parameter_change/community-pool-spend/software_upgrade/cancel-software-upgrade")
	return cmd
}

func ProposalUpdatePolicyAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-policy-admin [policy_addr_bech32] [new_admin_addr_bech32]",
		Short: "Submit a new admin for a policy proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			src, err := parseUpdatePolicyAdminArgs(args, clientCtx)
			if err != nil {
				return err
			}

			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return fmt.Errorf("deposit: %s", err)
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.UpdateAdminProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Policy:      src.Policy,
				NewAdmin:    src.NewAdmin,
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: Permission of proposal, types: policy-store/migrate-policy/instantiate-policy/set-policy-admin/clear-policy-admin/text/parameter_change/community-pool-spend/software_upgrade/cancel-software-upgrade")
	return cmd
}

func ProposalClearPolicyAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-policy-admin [policy_addr_bech32]",
		Short: "Submit a clear admin for a policy to prevent further migrations proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalTitle, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			depositArg, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return fmt.Errorf("deposit: %s", err)
			}
			deposit, err := sdk.ParseCoinsNormalized(depositArg)
			if err != nil {
				return err
			}

			content := types.ClearAdminProposal{
				Title:       proposalTitle,
				Description: proposalDescr,
				Policy:      args[0],
			}

			msg, err := govtypes.NewMsgSubmitProposal(&content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	// proposal flags
	cmd.Flags().String(cli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(cli.FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	// type values must match the "ProposalHandler" "routes" in cli
	cmd.Flags().String(flagProposalType, "", "Permission of proposal, types: store-code/instantiate/migrate/set-policy-admin/clear-policy-admin/text/parameter_change/software_upgrade")
	return cmd
}
