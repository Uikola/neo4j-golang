package host

import (
	"context"
	"github.com/Uikola/neo4j-golang/internal/entity"
)

type Repository interface {
	Save(ctx context.Context, host entity.Host) error
}

type Handler struct {
	hostRepository Repository
}

func NewHandler(hostRepository Repository) *Handler {
	return &Handler{hostRepository: hostRepository}
}
