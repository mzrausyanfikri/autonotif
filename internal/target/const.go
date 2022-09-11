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
💥 <strong><i>New Proposal on %s - ID: %s</i></strong>

📰 <strong>%s</strong>

Status: <strong>%s</strong>
Type: <strong>%s</strong>
Voting start time: <strong>%s</strong>
Voting end time: <strong>%s</strong>

🗳️ <strong><a href="%s/%s">View details and cast your VOTE</a></strong> 🗳️

🏛️ <strong><a href="%s">View all active proposals</a></strong> 🏛️
`
)
