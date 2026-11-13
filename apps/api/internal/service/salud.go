package service

import (
	"context"

	"base-tecnica-grupo6/internal/store"
)

type Salud struct {
	baseDatos *store.BaseDatos
}

func NuevaSalud(baseDatos *store.BaseDatos) *Salud {
	return &Salud{baseDatos: baseDatos}
}

func (s *Salud) EstaDisponible(ctx context.Context) error {
	return s.baseDatos.RevisarConexion(ctx)
}
