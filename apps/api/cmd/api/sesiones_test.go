package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSesionesSimuladasNoDisponiblesEnProduccion(t *testing.T) {
	rutas := http.NewServeMux()
	crearAPI(rutas, "prod", nil)

	respuesta := httptest.NewRecorder()
	rutas.ServeHTTP(respuesta, httptest.NewRequest(http.MethodPost, "/api/sesiones", nil))
	if respuesta.Code != http.StatusNotFound {
		t.Fatalf("el inicio de sesión simulado está expuesto en producción: %d", respuesta.Code)
	}
}
