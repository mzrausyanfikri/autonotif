package repository

import (
	"context"
	"time"

	"github.com/aimzeter/autonotif/config"
	"github.com/aimzeter/autonotif/entity"
	repopgx "github.com/aimzeter/autonotif/internal/repository/pgx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	lastIdQuery       = "SELECT proposal_id FROM proposals WHERE chain_type=$1 ORDER BY proposal_id DESC LIMIT 1"
	insertQuery       = "INSERT INTO proposals (proposal_id, chain_type, raw_data, created_at) VALUES ($1, $2, $3, $4)"
	revokeLastIDQuery = "DELETE FROM proposals WHERE chain_type=$1 AND proposal_id > $2"
)

type ProposalPSQL struct {
	connPool *pgxpool.Pool
}

func NewProposalPSQL(cfg config.Postgresql) (*ProposalPSQL, error) {
	connPool, err := repopgx.NewPool(cfg)
	return &ProposalPSQL{connPool: connPool}, err
}

func (r *ProposalPSQL) GetLastID(ctx context.Context, chainType entity.BlockchainType) (int, error) {
	row := r.connPool.QueryRow(ctx, lastIdQuery, chainType.String())
	var lastID int
	err := row.Scan(&lastID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return lastID, err
}

func (r *ProposalPSQL) Set(ctx context.Context, p entity.Proposal) error {
	_, err := r.connPool.Exec(ctx, insertQuery, p.ID, p.ChainType.String(), p.RawData, time.Now().UTC())
	return err
}

func (r *ProposalPSQL) RevokeLastID(ctx context.Context, chainType entity.BlockchainType, lastID int) error {
	_, err := r.connPool.Exec(ctx, revokeLastIDQuery, chainType.String(), lastID)
	return err
}
