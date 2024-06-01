package host

import (
	"encoding/json"
	"github.com/Uikola/neo4j-golang/internal/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h Handler) Save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	event := cloudevents.NewEvent()

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Error().Err(err).Msg("failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid CloudEvent format"})
		return
	}

	var interfaces []entity.Interface
	if err := json.Unmarshal(event.Data(), &interfaces); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal cloud event data")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	host := entity.Host{
		Hostname:   r.RemoteAddr,
		Interfaces: interfaces,
	}

	if err := h.hostRepository.Save(ctx, host); err != nil {
		log.Error().Err(err).Msg("error while saving host data")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
