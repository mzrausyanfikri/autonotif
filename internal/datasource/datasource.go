package datasource

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/aimzeter/autonotif/entity"
)

type Datasource struct {
	client *http.Client
}

type getProposalParams struct {
	ProposalID int
	URL        string
	MaxRetry   int
	Timeout    time.Duration
}

func NewDatasource() *Datasource {
	client := &http.Client{}

	return &Datasource{
		client: client,
	}
}

func (d *Datasource) GetProposalDetail(ctx context.Context, p *entity.Proposal) (*entity.Proposal, error) {
	var data *gabs.Container
	var err error

	conf := p.ChainConfig.ChainAPI
	for _, nodeAddr := range conf.Nodepool {
		params := getProposalParams{
			ProposalID: p.ID,
			URL:        nodeAddr + conf.Endpoint,
			MaxRetry:   conf.Retry,
			Timeout:    conf.Timeout,
		}

		data, err = d.getProposalDetail(ctx, nodeAddr, params)
		p.Data = data

		if err == nil {
			break
		}
	}

	return p, err
}

func (d *Datasource) getProposalDetail(ctx context.Context, nodeAddr string, params getProposalParams) (*gabs.Container, error) {
	reqCtx, cancel := context.WithTimeout(ctx, params.Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%d", params.URL, params.ProposalID)
	req, _ := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)

	var resp *http.Response
	var errResp error

	for attempt := 1; ; attempt++ {
		resp, errResp = d.client.Do(req)
		if errResp == nil {
			break
		}

		if attempt >= params.MaxRetry {
			return nil, fmt.Errorf("failed to http get proposal: %s", errResp.Error())
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return d.handleNotFound(ctx, params)
	}

	if resp.StatusCode != http.StatusOK {
		return d.handleNonOK(params, resp)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err.Error())
	}

	jsonParsed, err := gabs.ParseJSON(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json bytes: %s", err.Error())
	}

	return jsonParsed, nil
}

func (d *Datasource) handleNotFound(ctx context.Context, params getProposalParams) (*gabs.Container, error) {
	return nil, entity.ErrProposalNotYetExistInDatasource
}

func (d *Datasource) handleNonOK(params getProposalParams, resp *http.Response) (*gabs.Container, error) {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return nil, fmt.Errorf("proposal id %d got unhandled response status (%d) from %s with body string: %s", params.ProposalID, resp.StatusCode, params.URL, string(bodyBytes))
}
