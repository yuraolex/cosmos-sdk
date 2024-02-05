package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gov"

	// StoreKey is the store key string for gov
	StoreKey = ModuleName

	// RouterKey is the message route for gov
	RouterKey = ModuleName
)

var (
	ProposalsKeyPrefix               = collections.NewPrefix(0)  // ProposalsKeyPrefix stores the proposals raw bytes.
	ProposalIDKey                    = collections.NewPrefix(3)  // ProposalIDKey stores the sequence representing the next proposal ID.
	DepositsKeyPrefix                = collections.NewPrefix(16) // DepositsKeyPrefix stores deposits.
	VotesKeyPrefix                   = collections.NewPrefix(32) // VotesKeyPrefix stores the votes of proposals.
	ParamsKey                        = collections.NewPrefix(48) // ParamsKey stores the module's params.
	ConstitutionKey                  = collections.NewPrefix(49) // ConstitutionKey stores a chain's constitution.
	ProposalVoteOptionsKeyPrefix     = collections.NewPrefix(50) // ProposalVoteOptionsKeyPrefix stores the vote options of proposals.
	MessageBasedParamsKey            = collections.NewPrefix(51) // MessageBasedParamsKey stores the message based gov params.
	ProposalEndDepositPeriodIndexKey = collections.NewPrefix(52) // ProposalEndDepositPeriodIndexKey stores the proposal index.
	ProposalEndVotingPeriodIndexKey  = collections.NewPrefix(53) // ProposalEndVotingPeriodIndexKey stores the proposal index.
)

// Reserved kvstore keys
var (
	_ = collections.NewPrefix(1)
	_ = collections.NewPrefix(2)
	_ = collections.NewPrefix(4)
)
