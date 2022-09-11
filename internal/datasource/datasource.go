package datasource

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/aimzeter/autonotif/entity"
)

const existOffset = 7

type Datasource struct {
	client *http.Client
}

type getProposalParams struct {
	ProposalID int
	ChainType  string
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
	var errMsgs []string

	conf := p.ChainConfig.ChainAPI
	for _, nodeAddr := range conf.Nodepool {
		params := getProposalParams{
			ProposalID: p.ID,
			ChainType:  p.ChainType,
			URL:        nodeAddr + conf.Endpoint,
			MaxRetry:   conf.Retry,
			Timeout:    conf.Timeout,
		}

		data, err := d.getProposalDetail(ctx, params)
		p.Data = data

		if err == nil {
			return p, nil
		}

		errMsgs = append(errMsgs, err.Error())
	}

	return nil, errors.New(strings.Join(errMsgs, " <> "))
}

func (d *Datasource) getProposalDetail(ctx context.Context, params getProposalParams) (*gabs.Container, error) {
	resp, err := d.callProposalDetailAPI(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("callProposalDetailAPI: %s", err.Error())
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

func (d *Datasource) callProposalDetailAPI(ctx context.Context, params getProposalParams) (*http.Response, error) {
	reqCtx, cancel := context.WithTimeout(context.Background(), params.Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%d", params.URL, params.ProposalID)

	log.Printf("INFO | call %s \n", url)
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

	return resp, nil
}

func (d *Datasource) handleNotFound(ctx context.Context, params getProposalParams) (*gabs.Container, error) {
	exist, err := d.isNextProposalExist(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("isNextProposalExist: %s", err.Error())
	}

	if !exist {
		return nil, entity.ErrProposalNotYetExistInDatasource
	}

	p := entity.DeletedProposal
	return p.Data, nil
}

func (d *Datasource) handleNonOK(params getProposalParams, resp *http.Response) (*gabs.Container, error) {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return nil, fmt.Errorf("proposal id %d got unhandled response status (%d) from %s with body string: %s", params.ProposalID, resp.StatusCode, params.URL, string(bodyBytes))
}

func (d *Datasource) isNextProposalExist(ctx context.Context, params getProposalParams) (bool, error) {
	limitID := params.ProposalID + existOffset

	for currentID := params.ProposalID + 1; ; currentID++ {
		params.ProposalID = currentID
		resp, err := d.callProposalDetailAPI(ctx, params)
		if err != nil {
			return false, fmt.Errorf("callProposalDetailAPI: %s", err.Error())
		}

		_ = resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true, nil
		}

		if currentID >= limitID {
			return false, nil
		}
	}
}
