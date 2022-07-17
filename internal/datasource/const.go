package datasource

type ProposalErrorResponse_COSMOS struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}
