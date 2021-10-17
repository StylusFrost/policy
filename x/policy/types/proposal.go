package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type ProposalType string

const (
	ProposalTypeStoreRego         ProposalType = "StoreRego"
	ProposalTypeInstantiatePolicy ProposalType = "InstantiatePolicy"
)

// DisableAllProposals contains no policy gov types.
var DisableAllProposals []ProposalType

// EnableAllProposals contains all policy gov types as keys.
var EnableAllProposals = []ProposalType{
	ProposalTypeStoreRego,
	ProposalTypeInstantiatePolicy,
}

// ConvertToProposals maps each key to a ProposalType and returns a typed list.
// If any string is not a valid type (in this file), then return an error
func ConvertToProposals(keys []string) ([]ProposalType, error) {
	valid := make(map[string]bool, len(EnableAllProposals))
	for _, key := range EnableAllProposals {
		valid[string(key)] = true
	}

	proposals := make([]ProposalType, len(keys))
	for i, key := range keys {
		if _, ok := valid[key]; !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "'%s' is not a valid ProposalType", key)
		}
		proposals[i] = ProposalType(key)
	}
	return proposals, nil
}

func init() { // register new content types with the sdk
	govtypes.RegisterProposalType(string(ProposalTypeStoreRego))
	govtypes.RegisterProposalType(string(ProposalTypeInstantiatePolicy))

	govtypes.RegisterProposalTypeCodec(&StoreRegoProposal{}, "policy/StoreRegoProposal")
	govtypes.RegisterProposalTypeCodec(&InstantiatePolicyProposal{}, "policy/InstantiatePolicyProposal")

}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p StoreRegoProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *StoreRegoProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p StoreRegoProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p StoreRegoProposal) ProposalType() string { return string(ProposalTypeStoreRego) }

// ValidateBasic validates the proposal
func (p StoreRegoProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.RunAs); err != nil {
		return sdkerrors.Wrap(err, "run as")
	}

	if err := validateRegoCode(p.REGOByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	if err := validateSourceURL(p.Source); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "source %s", err.Error())
	}

	if !json.Valid(p.EntryPoints) {
		return sdkerrors.Wrap(ErrInvalid, "entry points json")
	}

	if p.InstantiatePermission != nil {
		if err := p.InstantiatePermission.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "instantiate permission")
		}
	}
	return nil
}

// String implements the Stringer interface.
func (p StoreRegoProposal) String() string {
	return fmt.Sprintf(`Store Code Proposal:
  Title:       %s
  Description: %s
  Run as:      %s
  RegoCode:    %X
  EntryPoints: %q
  Source:      %s
`, p.Title, p.Description, p.RunAs, p.REGOByteCode, p.EntryPoints, p.Source)
}

// MarshalYAML pretty prints the rego byte code
func (p StoreRegoProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title                 string        `yaml:"title"`
		Description           string        `yaml:"description"`
		RunAs                 string        `yaml:"run_as"`
		REGOByteCode          string        `yaml:"rego_byte_code"`
		EntryPoints           string        `yaml:"entry_points"`
		Source                string        `yaml:"source"`
		InstantiatePermission *AccessConfig `yaml:"instantiate_permission"`
	}{
		Title:                 p.Title,
		Description:           p.Description,
		RunAs:                 p.RunAs,
		REGOByteCode:          base64.StdEncoding.EncodeToString(p.REGOByteCode),
		EntryPoints:           string(p.EntryPoints),
		Source:                p.Source,
		InstantiatePermission: p.InstantiatePermission,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p InstantiatePolicyProposal) ProposalRoute() string { return RouterKey }

// GetTitle returns the title of the proposal
func (p *InstantiatePolicyProposal) GetTitle() string { return p.Title }

// GetDescription returns the human readable description of the proposal
func (p InstantiatePolicyProposal) GetDescription() string { return p.Description }

// ProposalType returns the type
func (p InstantiatePolicyProposal) ProposalType() string {
	return string(ProposalTypeInstantiatePolicy)
}

// ValidateBasic validates the proposal
func (p InstantiatePolicyProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.RunAs); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "run as")
	}

	if p.RegoID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "rego id is required")
	}

	if err := validateLabel(p.Label); err != nil {
		return err
	}

	if !p.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	if len(p.Admin) != 0 {
		if _, err := sdk.AccAddressFromBech32(p.Admin); err != nil {
			return err
		}
	}
	if !json.Valid(p.EntryPoints) {
		return sdkerrors.Wrap(ErrInvalid, "entry points json")
	}

	return nil
}

// String implements the Stringer interface.
func (p InstantiatePolicyProposal) String() string {
	return fmt.Sprintf(`Instantiate Rego Proposal:
  Title:       %s
  Description: %s
  Run as:      %s
  Admin:       %s
  Rego id:     %d
  Label:       %s
  EntryPoints: %q
  Funds:       %s
`, p.Title, p.Description, p.RunAs, p.Admin, p.RegoID, p.Label, p.EntryPoints, p.Funds)
}

// MarshalYAML pretty prints the init message
func (p InstantiatePolicyProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		RunAs       string    `yaml:"run_as"`
		Admin       string    `yaml:"admin"`
		RegoID      uint64    `yaml:"code_id"`
		Label       string    `yaml:"label"`
		EntryPoints string    `yaml:"entry_points"`
		Funds       sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		RunAs:       p.RunAs,
		Admin:       p.Admin,
		RegoID:      p.RegoID,
		Label:       p.Label,
		EntryPoints: string(p.EntryPoints),
		Funds:       p.Funds,
	}, nil
}

func validateProposalCommons(title, description string) error {
	if strings.TrimSpace(title) != title {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title must not start/end with white spaces")
	}
	if len(title) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if len(title) > govtypes.MaxTitleLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength)
	}
	if strings.TrimSpace(description) != description {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description must not start/end with white spaces")
	}
	if len(description) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank")
	}
	if len(description) > govtypes.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength)
	}
	return nil
}
