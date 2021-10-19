package rest

import (
	"encoding/json"
	"net/http"

	"github.com/StylusFrost/policy/x/policy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type StoreRegoProposalJsonReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description" yaml:"description"`
	Proposer    string    `json:"proposer" yaml:"proposer"`
	Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`

	RunAs string `json:"run_as" yaml:"run_as"`
	// REGOByteCode can be raw or gzip compressed
	REGOByteCode []byte `json:"policy_byte_rego" yaml:"policy_byte_rego"`
	// Source is a valid absolute HTTPS URI to the policy's source code, optional
	Source string `json:"source" yaml:"source"`
	// Valid entry points json encoded
	EntryPoints json.RawMessage `json:"entry_points"  yaml:"entry_points"`
	// InstantiatePermission to apply on policy creation, optional
	InstantiatePermission *types.AccessConfig `json:"instantiate_permission" yaml:"instantiate_permission"`
}

func (s StoreRegoProposalJsonReq) Content() govtypes.Content {
	return &types.StoreRegoProposal{
		Title:                 s.Title,
		Description:           s.Description,
		RunAs:                 s.RunAs,
		REGOByteCode:          s.REGOByteCode,
		EntryPoints:           s.EntryPoints,
		Source:                s.Source,
		InstantiatePermission: s.InstantiatePermission,
	}
}
func (s StoreRegoProposalJsonReq) GetProposer() string {
	return s.Proposer
}
func (s StoreRegoProposalJsonReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s StoreRegoProposalJsonReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func StoreRegoProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "policy_store_rego",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req StoreRegoProposalJsonReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type InstantiateProposalJsonReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	RunAs string `json:"run_as" yaml:"run_as"`
	// Admin is an optional address that can execute migrations
	Admin       string          `json:"admin,omitempty" yaml:"admin"`
	Rego        uint64          `json:"rego_id" yaml:"rego_id"`
	Label       string          `json:"label" yaml:"label"`
	EntryPoints json.RawMessage `json:"entry_points" yaml:"entry_points"`
	Funds       sdk.Coins       `json:"funds" yaml:"funds"`
}

func (s InstantiateProposalJsonReq) Content() govtypes.Content {
	return &types.InstantiatePolicyProposal{
		Title:       s.Title,
		Description: s.Description,
		RunAs:       s.RunAs,
		Admin:       s.Admin,
		RegoID:      s.Rego,
		Label:       s.Label,
		EntryPoints: s.EntryPoints,
		Funds:       s.Funds,
	}
}
func (s InstantiateProposalJsonReq) GetProposer() string {
	return s.Proposer
}
func (s InstantiateProposalJsonReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s InstantiateProposalJsonReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func InstantiateProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "policy_instantiate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req InstantiateProposalJsonReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type MigrateProposalJsonReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Policy      string          `json:"policy" yaml:"policy"`
	RegoID      uint64          `json:"rego_id" yaml:"rego_id"`
	EntryPoints json.RawMessage `json:"entry_points" yaml:"entry_points"`
	// RunAs is the role that is passed to the policy's environment
	RunAs string `json:"run_as" yaml:"run_as"`
}

func (s MigrateProposalJsonReq) Content() govtypes.Content {
	return &types.MigratePolicyProposal{
		Title:       s.Title,
		Description: s.Description,
		Policy:      s.Policy,
		RegoID:      s.RegoID,
		EntryPoints: s.EntryPoints,
		RunAs:       s.RunAs,
	}
}
func (s MigrateProposalJsonReq) GetProposer() string {
	return s.Proposer
}
func (s MigrateProposalJsonReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s MigrateProposalJsonReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func MigrateProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "policy_migrate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req MigrateProposalJsonReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type UpdateAdminJsonReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	NewAdmin string `json:"new_admin" yaml:"new_admin"`
	Policy   string `json:"policy" yaml:"policy"`
}

func (s UpdateAdminJsonReq) Content() govtypes.Content {
	return &types.UpdateAdminProposal{
		Title:       s.Title,
		Description: s.Description,
		Policy:      s.Policy,
		NewAdmin:    s.NewAdmin,
	}
}
func (s UpdateAdminJsonReq) GetProposer() string {
	return s.Proposer
}
func (s UpdateAdminJsonReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s UpdateAdminJsonReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func UpdatePolicyAdminProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "policy_update_admin",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req UpdateAdminJsonReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type ClearAdminJsonReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Policy string `json:"policy" yaml:"policy"`
}

func (s ClearAdminJsonReq) Content() govtypes.Content {
	return &types.ClearAdminProposal{
		Title:       s.Title,
		Description: s.Description,
		Policy:      s.Policy,
	}
}
func (s ClearAdminJsonReq) GetProposer() string {
	return s.Proposer
}
func (s ClearAdminJsonReq) GetDeposit() sdk.Coins {
	return s.Deposit
}
func (s ClearAdminJsonReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}
func ClearPolicyAdminProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "policy_clear_admin",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req ClearAdminJsonReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type policyProposalData interface {
	Content() govtypes.Content
	GetProposer() string
	GetDeposit() sdk.Coins
	GetBaseReq() rest.BaseReq
}

func toStdTxResponse(cliCtx client.Context, w http.ResponseWriter, data policyProposalData) {
	proposerAddr, err := sdk.AccAddressFromBech32(data.GetProposer())
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, err := govtypes.NewMsgSubmitProposal(data.Content(), data.GetDeposit(), proposerAddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	baseReq := data.GetBaseReq().Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}
	tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
}