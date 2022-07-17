package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aimzeter/autonotif/entity"
	scribble "github.com/nanobox-io/golang-scribble"
)

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

func (r *ProposalLocalfile) GetLastID(ctx context.Context, chainType entity.BlockchainType) (int, error) {
	collection := strings.ToLower(chainType.String())
	records, err := r.db.ReadAll(collection)
	if err != nil {
		return 0, fmt.Errorf("db read all: %s", err.Error())
	}

	all := []entity.Proposal{}
	for _, rec := range records {
		var p entity.Proposal
		if err := json.Unmarshal([]byte(rec), &p); err != nil {
			return 0, fmt.Errorf("json unmarshall: %s", err.Error())
		}
		all = append(all, p)
	}

	allIds := []int{}
	for _, p := range all {
		allIds = append(allIds, p.ID)
	}

	if len(allIds) == 0 {
		return 0, nil
	}

	sort.Ints(allIds)
	return allIds[len(allIds)-1], nil
}

func (r *ProposalLocalfile) Set(ctx context.Context, p entity.Proposal) error {
	collection := strings.ToLower(p.ChainType.String())
	return r.db.Write(collection, strconv.Itoa(p.ID), p)
}

func (r *ProposalLocalfile) RevokeLastID(ctx context.Context, chainType entity.BlockchainType, intendedLastID int) error {
	collection := strings.ToLower(chainType.String())
	lastID, err := r.GetLastID(ctx, chainType)
	if err != nil {
		return fmt.Errorf("GetLastID: %s", err.Error())
	}

	i := intendedLastID + 1
	for {
		if i > lastID {
			break
		}

		r.db.Delete(collection, strconv.Itoa(i))
		i++
	}

	return nil
}
