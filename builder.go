package autonotif

import (
	"context"
	"errors"

	"github.com/aimzeter/autonotif/config"
	"github.com/aimzeter/autonotif/entity"
	"github.com/aimzeter/autonotif/internal/datasource"
	"github.com/aimzeter/autonotif/internal/repository"
	"github.com/aimzeter/autonotif/internal/target"
)

type ProposalStore interface {
	Set(ctx context.Context, p *entity.Proposal) error
	GetLastID(ctx context.Context, chainType string) (int, error)
	RevokeLastID(ctx context.Context, chainType string, lastID int) error
}

type DatasourceAPI interface {
	GetProposalDetail(ctx context.Context, p *entity.Proposal) (*entity.Proposal, error)
}

type Notifier interface {
	SendMessage(ctx context.Context, p *entity.Proposal) error
}

type Dependencies struct {
	conf     *config.Config
	dsAPI    DatasourceAPI
	dsStore  ProposalStore
	notifier Notifier
}

func BuildDependencies(conf *config.Config) (*Dependencies, error) {
	dsAPI := datasource.NewDatasource()

	dsStore, storeErr := BuildRepository(conf.Repository)
	if storeErr != nil {
		return nil, storeErr
	}

	notifier, notifErr := target.NewTelegram(conf.Target.Telegram.Token, conf.Target.Telegram.ChannelID)
	if notifErr != nil {
		return nil, notifErr
	}

	return &Dependencies{
		conf:     conf,
		dsAPI:    dsAPI,
		dsStore:  dsStore,
		notifier: notifier,
	}, nil
}

func BuildAutonotif(d *Dependencies) *Autonotif {
	return &Autonotif{
		d: d,
	}
}

func BuildRepository(cfg config.Repository) (ProposalStore, error) {
	switch {
	case cfg.Localfile.Enable:
		return repository.NewProposalLocalfile(cfg.Localfile.Dir)
	case cfg.Postgresql.Enable:
		return repository.NewProposalPSQL(cfg.Postgresql)
	}

	return nil, errors.New("no repository enabled")
}
