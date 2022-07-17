package autonotif

import (
	"context"
	"fmt"
	"log"

	"github.com/aimzeter/autonotif/entity"
)

type ProposalStore interface {
	GetLastID(ctx context.Context, chainType entity.BlockchainType) (int, error)
	Set(ctx context.Context, p entity.Proposal) error
	RevokeLastID(ctx context.Context, chainType entity.BlockchainType, lastID int) error
}

type DatasourceAPI interface {
	GetProposalDetail(ctx context.Context, id int) (entity.Proposal, error)
}

type Notifier interface {
	SendMessage(ctx context.Context, p entity.Proposal) error
}

type Autonotif struct {
	d *Dependencies
}

func (a *Autonotif) Terminate() error {
	return nil
}

func (a *Autonotif) HealthCheck() error {
	return nil
}

func (a *Autonotif) Run() error {
	for _, chainType := range entity.AllBlockchainType {
		log.Printf("INFO | chain %s running...\n", chainType)
		err := a.notifyRecentProposal(chainType)
		if err != nil {
			log.Printf("ERROR | chain %s runner.notifyRecentProposal: %s\n", chainType, err)
			continue
		}

		log.Printf("INFO | chain %s run successfully\n", chainType)
	}

	return nil
}

func (a *Autonotif) notifyRecentProposal(chainType entity.BlockchainType) error {
	ctx := context.Background()

	lastID, err := a.d.dsStore.GetLastID(ctx, chainType)
	if err != nil {
		return fmt.Errorf("dsStore.GetLastID: %s", err.Error())
	}

	nextID := lastID + 1

	p, err := a.d.dsAPI.GetProposalDetail(ctx, nextID)
	if err == entity.ErrProposalNotYetExistInDatasource {
		return nil
	}

	if err != nil {
		return fmt.Errorf("dsAPI.GetProposalDetail: %s", err.Error())
	}

	err = a.notify(ctx, p)
	if err != nil {
		return fmt.Errorf("notify: %s", err.Error())
	}

	err = a.d.dsStore.Set(ctx, p)
	if err != nil {
		return fmt.Errorf("dsStore.Set: %s", err.Error())
	}

	return nil
}

func (a *Autonotif) notify(ctx context.Context, p entity.Proposal) error {
	if !p.IsShouldNotify() {
		return nil
	}

	err := a.d.notifier.SendMessage(ctx, p)
	if err != nil {
		return fmt.Errorf("notifier.SendMessage: %s", err.Error())
	}

	return nil
}
