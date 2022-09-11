package autonotif

import "github.com/Jeffail/gabs/v2"

type Proposal struct {
	ID        int
	ChainType string
	Data      *gabs.Container
}
