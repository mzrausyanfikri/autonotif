package autonotif

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	chainType, lastID, ok := a.parseForceLastIDHeader(w, r)
	if !ok {
		return
	}

	err := a.ForceLastID(ctx, chainType, lastID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("Error: %s\n", err.Error()))
		return
	}

	_, _ = io.WriteString(w, fmt.Sprintf("%s: %d revoked successfully\n", chainType, lastID))
}

func (a *Autonotif) ForceLastID(ctx context.Context, chainType string, lastID int) error {
	// start proposal from zero
	if lastID == -1 {
		err := a.d.dsStore.RevokeLastID(ctx, chainType, 0)
		if err != nil {
			return fmt.Errorf("dsStore.RevokeLastID: %s", err.Error())
		}
		return nil
	}

	p := &entity.Proposal{
		ID:          lastID,
		ChainType:   chainType,
		ChainConfig: a.d.conf.ChainList[chainType],
	}

	p, err := a.d.dsAPI.GetProposalDetail(ctx, p)
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

func (a *Autonotif) parseForceLastIDHeader(w http.ResponseWriter, r *http.Request) (string, int, bool) {
	chainType := r.Header.Get("chain")
	if chainType == "" {
		return handleBadRequest(w, "Invalid value: chain empty\n")
	}

	chainType = strings.ToLower(chainType)
	_, ok := a.d.conf.ChainList[chainType]
	if !ok {
		handleBadRequest(w, "Invalid value: chain unknown\n")
	}

	lastIDStr := r.Header.Get("lastId")
	lastID, err := strconv.Atoi(lastIDStr)
	if err != nil || lastID == 0 {
		return handleBadRequest(w, "Invalid value: lastId\n")
	}

	return chainType, lastID, true
}

func handleBadRequest(w http.ResponseWriter, msg string) (string, int, bool) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = io.WriteString(w, msg)
	return "", 0, false
}
