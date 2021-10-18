package client

import (
	"github.com/StylusFrost/policy/x/policy/client/cli"
	"github.com/StylusFrost/policy/x/policy/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var StoreRegoProposalHandler = govclient.NewProposalHandler(cli.StoreRegoProposalCmd, rest.StoreRegoProposalHandler)
var ProposalInstantiatePolicyHandler = govclient.NewProposalHandler(cli.ProposalInstantiatePolicyCmd, rest.InstantiateProposalHandler)
var ProposalMigratePolicyHandlerHandler = govclient.NewProposalHandler(cli.ProposalMigratePolicyCmd, rest.MigrateProposalHandler)
