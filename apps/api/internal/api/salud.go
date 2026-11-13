package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"base-tecnica-grupo6/internal/service"
)

type RespuestaSalud struct {
	Body struct {
		Estado string `json:"estado" example:"ok" doc:"Siempre \"ok\" cuando la respuesta es 200"`
	}
}

type ManejadorSalud struct {
	salud *service.Salud
}

func RegistrarRutasSalud(apiHuma huma.API, salud *service.Salud) {
	manejador := &ManejadorSalud{salud: salud}

	vida := huma.Operation{OperationID: "revisar-vida", Method: http.MethodGet, Path: "/healthz"}
	vida.Summary = "Indica si el proceso está vivo"
	vida.Tags = []string{"Salud"}
	huma.Register(apiHuma, vida, manejador.revisarVida)

	disponibilidad := huma.Operation{OperationID: "revisar-disponibilidad", Method: http.MethodGet, Path: "/readyz"}
	disponibilidad.Summary = "Indica si el servicio puede atender: responde 200 solo si PostgreSQL contesta"
	disponibilidad.Tags = []string{"Salud"}
	disponibilidad.Errors = []int{http.StatusServiceUnavailable}
	huma.Register(apiHuma, disponibilidad, manejador.revisarDisponibilidad)
}

func (m *ManejadorSalud) revisarVida(ctx context.Context, entrada *struct{}) (*RespuestaSalud, error) {
	return respuestaOK(), nil
}

func (m *ManejadorSalud) revisarDisponibilidad(ctx context.Context, entrada *struct{}) (*RespuestaSalud, error) {
	err := m.salud.EstaDisponible(ctx)
	if err != nil {
		slog.Warn("la base de datos no responde", "error", err)
		return nil, huma.Error503ServiceUnavailable("La base de datos no responde")
	}

	return respuestaOK(), nil
}

func respuestaOK() *RespuestaSalud {
	respuesta := &RespuestaSalud{}
	respuesta.Body.Estado = "ok"
	return respuesta
}
