package target

import (
	"time"
	_ "time/tzdata"
)

const (
	timeTZ   = "Asia/Jakarta"
	msgLimit = 4096
)

var (
	timeLoc, _      = time.LoadLocation(timeTZ)
	messageTemplate = `
💥 <strong><i>New Proposal on Cosmos - ID: %s</i></strong>

📰 <strong>%s</strong>

Type: <strong>%s</strong>
Voting start time: <strong>%s</strong>
Voting end time: <strong>%s</strong>

<i>%s</i>

🗳️ <strong><a href="https://www.mintscan.io/cosmos/proposals/%s">view details and cast your VOTE</a></strong> 🗳️

🏛️ <strong><a href="https://www.mintscan.io/cosmos/proposals">view all active proposals</a></strong> 🏛️
`
)
