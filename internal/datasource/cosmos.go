package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
        "strings"

	"github.com/aimzeter/autonotif/entity"
	fakeua "github.com/wux1an/fake-useragent"
)

const (
	getProposalsEndpoint = "/cosmos/gov/v1beta1/proposals"
	timeout              = 10 * time.Second
	maxRetry             = 5
)

type ProposalList_COSMOS_v1 struct {
	Proposals []struct {
		ProposalID string `json:"proposal_id"`
	} `json:"proposals"`
}

type Cosmos struct {
	nodepool []string
	client   *http.Client
}

func NewCosmos(nodepool []string) *Cosmos {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Cosmos{
		nodepool: nodepool,
		client:   client,
	}
}

func (c *Cosmos) GetProposalDetail(ctx context.Context, id int) (entity.Proposal, error) {
	var p entity.Proposal
	var err error

	for _, node := range c.nodepool {
		p, err = c.getProposalByID(ctx, node, id)
		if err == nil {
			break
		}
	}

	return p, err
}

type Response struct {
	STATUS string `json:"status"`
}

func (c *Cosmos) getProposalByID(ctx context.Context, nodeAdress string, id int) (entity.Proposal, error) {
	url := fmt.Sprintf("%s%s/%d", nodeAdress, getProposalsEndpoint, id)
	req, _ := prepareRequest(ctx, http.MethodGet, url, nil)
	var resp *http.Response
	var errResp error
	for attempt := 1; ; attempt++ {
		resp, errResp = c.client.Do(req)
		if errResp == nil {
			break
		}

		if attempt >= maxRetry {
			return entity.Proposal{}, fmt.Errorf("failed to http get proposal: %s", errResp.Error())
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

                bodyString := string(bodyBytes)
                bodyParse := strings.ReplaceAll(bodyString, "\\n", "\n")
                bodyFinal := strings.ReplaceAll(bodyParse, "\n", "\\n")
                bodyFinalBytes := []byte(bodyFinal)
                fmt.Println(bodyFinal)

		return entity.Proposal{
			ID:        id,
			ChainType: entity.BlockchainType_COSMOS,
			RawData:   string(bodyFinalBytes),
		}, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		lastID, err := c.getLatestProposalID(ctx, nodeAdress)
		if err != nil {
			return entity.Proposal{}, nil
		}

		if id > lastID {
			return entity.Proposal{}, entity.ErrProposalNotYetExistInDatasource
		}

		return entity.Proposal{
			ID:        id,
			ChainType: entity.BlockchainType_COSMOS,
			RawData:   entity.ProposalRawData_NOT_EXIST,
		}, nil
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return entity.Proposal{}, fmt.Errorf("got unhandled response status (%d) with body string: %s", resp.StatusCode, string(bodyBytes))
}

func (c *Cosmos) getLatestProposalID(ctx context.Context, nodeAdress string) (int, error) {
	url := fmt.Sprintf("%s%s", nodeAdress, getProposalsEndpoint)
	req, _ := prepareRequest(ctx, http.MethodGet, url, nil)
	var resp *http.Response
	var errResp error
	for attempt := 1; ; attempt++ {
		resp, errResp = c.client.Do(req)
		if errResp == nil {
			break
		}

		if attempt >= maxRetry {
			return 0, fmt.Errorf("failed to http get proposal: %s", errResp.Error())
		}
	}
	defer resp.Body.Close()

	var all ProposalList_COSMOS_v1
	err := json.NewDecoder(resp.Body).Decode(&all)
	if err != nil {
		return 0, fmt.Errorf("failed to json encode: %s", err.Error())
	}

	allId := []int{}
	for _, p := range all.Proposals {
		intID, _ := strconv.Atoi(p.ProposalID)
		allId = append(allId, intID)
	}

	sort.Ints(allId)
	return allId[len(allId)-1], nil
}

func prepareRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", fakeua.RandomType(fakeua.Chrome))
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	return req, err
}
