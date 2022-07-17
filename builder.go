package autonotif

import (
	"errors"

	"github.com/aimzeter/autonotif/config"
	"github.com/aimzeter/autonotif/internal/datasource"
	"github.com/aimzeter/autonotif/internal/repository"
	"github.com/aimzeter/autonotif/internal/target"
)

type Dependencies struct {
	dsAPI    DatasourceAPI
	dsStore  ProposalStore
	notifier Notifier
}

func BuildDependencies(cfg *config.Config) (*Dependencies, error) {
	dsAPI := datasource.NewCosmos(cfg.Datasource.Cosmos.Nodepool)

	dsStore, storeErr := BuildRepository(cfg.Repository)
	if storeErr != nil {
		return nil, storeErr
	}

	notifier, notifErr := target.NewTelegram(cfg.Target.Telegram.Token, cfg.Target.Telegram.ChannelID)
	if notifErr != nil {
		return nil, notifErr
	}

	return &Dependencies{
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
