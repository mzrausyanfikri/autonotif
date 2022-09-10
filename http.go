package autonotif

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aimzeter/autonotif/entity"
)

func (a *Autonotif) HealthHandler(w http.ResponseWriter, r *http.Request) {
	err := a.HealthCheck()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, "NOT OK\n")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK\n")
}

func (a *Autonotif) ForceLastIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	chainType, lastID, ok := parseForceLastIDHeader(w, r)
	if !ok {
		return
	}

	err := a.ForceLastID(ctx, chainType, lastID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error: %s\n", err.Error()))
		return
	}

	_, _ = io.WriteString(w, fmt.Sprintf("%s: %d revoked successfully\n", chainType.String(), lastID))
}

func (a *Autonotif) ForceLastID(ctx context.Context, chainType entity.BlockchainType, lastID int) error {
	// start proposal from zero
	if lastID == -1 {
		err := a.d.dsStore.RevokeLastID(ctx, chainType, 0)
		if err != nil {
			return fmt.Errorf("dsStore.RevokeLastID: %s", err.Error())
		}
		return nil
	}

	p, err := a.d.dsAPI.GetProposalDetail(ctx, lastID)
	if err != nil {
		return fmt.Errorf("dsAPI.GetProposalDetail: %s", err.Error())
	}

	err = a.d.dsStore.Set(ctx, p)
	if err != nil {
		return fmt.Errorf("dsStore.Set: %s", err.Error())
	}

	err = a.d.dsStore.RevokeLastID(ctx, chainType, lastID)
	if err != nil {
		return fmt.Errorf("dsStore.RevokeLastID: %s", err.Error())
	}

	return nil
}

func parseForceLastIDHeader(w http.ResponseWriter, r *http.Request) (entity.BlockchainType, int, bool) {
	ctStr := r.Header.Get("chain")
	if ctStr == "" {
		return handleBadRequest(w, "Invalid value: chain empty\n")
	}

	chainType := entity.BlockchainTypeFromString(ctStr)
	if chainType == entity.BlockchainType_OTHER {
		return handleBadRequest(w, "Invalid value: chain unknown\n")
	}

	lastIDStr := r.Header.Get("lastId")
	lastID, err := strconv.Atoi(lastIDStr)
	if err != nil || lastID == 0 {
		return handleBadRequest(w, "Invalid value: lastId\n")
	}

	return chainType, lastID, true
}

func handleBadRequest(w http.ResponseWriter, msg string) (entity.BlockchainType, int, bool) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = io.WriteString(w, msg)
	return 0, 0, false
}
