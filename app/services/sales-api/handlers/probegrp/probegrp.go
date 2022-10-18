package probegrp

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return json.NewEncoder(w).Encode(status)
}

func (h Handlers) Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return json.NewEncoder(w).Encode(status)
}
