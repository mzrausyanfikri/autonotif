package entity

type BlockchainType int

const (
	BlockchainType_OTHER BlockchainType = iota
	BlockchainType_COSMOS
)

var AllBlockchainType = []BlockchainType{
	BlockchainType_COSMOS,
}

type Proposal struct {
	ID        int
	ChainType BlockchainType
	RawData   string
}

func (e BlockchainType) String() string {
	return [...]string{"OTHER", "OSMOSIS"}[e]
}

func (e BlockchainType) EnumIndex() int {
	return int(e)
}

func (p Proposal) IsShouldNotify() bool {
	return p.RawData != ProposalRawData_NOT_EXIST
}

func BlockchainTypeFromString(str string) BlockchainType {
	for _, chainType := range AllBlockchainType {
		if chainType.String() == str {
			return chainType
		}
	}

	return BlockchainType_OTHER
}
