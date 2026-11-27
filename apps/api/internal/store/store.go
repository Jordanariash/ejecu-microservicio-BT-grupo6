package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseDatos struct {
	conexiones *pgxpool.Pool
}

// No falla si la base de datos aún no está arriba: las conexiones se abren al usarlas.
func Conectar(ctx context.Context, urlBaseDatos string) (*BaseDatos, error) {
	conexiones, err := pgxpool.New(ctx, urlBaseDatos)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL no es válida: %w", err)
	}

	return &BaseDatos{conexiones: conexiones}, nil
}

func (b *BaseDatos) RevisarConexion(ctx context.Context) error {
	ctx, cancelar := context.WithTimeout(ctx, 2*time.Second)
	defer cancelar()

	return b.conexiones.Ping(ctx)
}

func (b *BaseDatos) Cerrar() {
	b.conexiones.Close()
}
