package entity

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/aimzeter/autonotif/config"
)

var RevokedProposalData, _ = gabs.ParseJSON([]byte(`{"proposal": "REVOKED_PROPOSAL_DATA"}`))

var RevokedProposal = Proposal{
	ChainConfig: config.Chain{},
	Data:        RevokedProposalData,
}

type Proposal struct {
	ID          int
	ChainType   string
	ChainConfig config.Chain
	Data        *gabs.Container
}

func (p Proposal) IsShouldNotify() bool {
	return true
}

func (p Proposal) IsRevokedProposal() bool {
	return p.Data == RevokedProposalData
}
