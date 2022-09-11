package entity

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/aimzeter/autonotif/config"
)

var (
	RevokedProposalData, _ = gabs.ParseJSON([]byte(`{"proposal": "REVOKED_PROPOSAL_DATA"}`))
	DeletedProposalData, _ = gabs.ParseJSON([]byte(`{"proposal": "DELETED_PROPOSAL_DATA"}`))
)

var RevokedProposal = Proposal{
	ChainConfig: config.Chain{},
	Data:        RevokedProposalData,
}

var DeletedProposal = Proposal{
	ChainConfig: config.Chain{},
	Data:        DeletedProposalData,
}

type Proposal struct {
	ID          int
	ChainType   string
	ChainConfig config.Chain
	Data        *gabs.Container
}

func (p Proposal) IsShouldNotify() bool {
	return p.Data != DeletedProposal.Data
}

func (p Proposal) IsRevokedProposal() bool {
	return p.Data == RevokedProposalData
}
