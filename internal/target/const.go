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
💥 <strong><i>New Proposal on Osmosis - ID: %s</i></strong>

📰 <strong>%s</strong>

Type: <strong>%s</strong>
Voting start time: <strong>%s</strong>
Voting end time: <strong>%s</strong>


🗳️ <strong><a href="https://www.mintscan.io/osmosis/proposals/%s">View details and cast your VOTE</a></strong> 🗳️

`
)
