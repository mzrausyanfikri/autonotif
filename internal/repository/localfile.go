package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aimzeter/autonotif/entity"
	scribble "github.com/nanobox-io/golang-scribble"
)

var collectionNotFound = "%s.json: no such file or directory"

type ProposalLocalfile struct {
	db *scribble.Driver
}

func NewProposalLocalfile(dir string) (*ProposalLocalfile, error) {
	dir += "proposals/"
	db, err := scribble.New(dir, nil)
	if err != nil {
		return nil, err
	}

	return &ProposalLocalfile{db: db}, nil
}

func (r *ProposalLocalfile) GetLastID(ctx context.Context, chainType string) (int, error) {
	collection := strings.ToLower(chainType)
	records, err := r.db.ReadAll(collection)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf(collectionNotFound, collection)) {
			return 0, nil
		}

		return 0, fmt.Errorf("db read all: %s", err.Error())
	}

	all := []file{}
	for _, rec := range records {
		var f file
		if err := json.Unmarshal([]byte(rec), &f); err != nil {
			return 0, fmt.Errorf("json unmarshall: %s", err.Error())
		}
		all = append(all, f)
	}

	if len(all) == 0 {
		return 0, nil
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt > all[j].CreatedAt
	})

	return all[0].ProposalID, nil
}

func (r *ProposalLocalfile) Set(ctx context.Context, p *entity.Proposal) error {
	collection := strings.ToLower(p.ChainType)

	f := proposalToFile(p)
	f.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	return r.db.Write(collection, f.CreatedAt, f)
}

func (r *ProposalLocalfile) RevokeLastID(ctx context.Context, chainType string, intendedLastID int) error {
	p := entity.RevokedProposal
	p.ChainType = chainType
	p.ID = intendedLastID
	return r.Set(ctx, &p)
}

type file struct {
	ProposalID  int
	ChainType   string
	ChainConfig interface{}
	Data        string
	CreatedAt   string
}

func proposalToFile(p *entity.Proposal) file {
	return file{
		ProposalID:  p.ID,
		ChainType:   p.ChainType,
		ChainConfig: p.ChainConfig,
		Data:        p.Data.String(),
	}
}
