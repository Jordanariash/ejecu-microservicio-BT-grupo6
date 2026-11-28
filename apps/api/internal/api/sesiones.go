package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const valorSesionSimulada = "simulada"

type EntradaInicioSesion struct {
	Body struct {
		Correo string `json:"correo" format:"email" maxLength:"254" example:"operador@base-tecnica.test"`
		Clave  string `json:"clave" minLength:"1" maxLength:"128" example:"demo"`
	}
}

type EntradaSesionActual struct {
	ValorCookie string `cookie:"sesion"`
}

type Sesion struct {
	ID       string   `json:"id" format:"uuid"`
	Correo   string   `json:"correo" format:"email"`
	Rol      string   `json:"rol"`
	Permisos []string `json:"permisos" nullable:"false"`
}

type RespuestaInicioSesion struct {
	Cookie string `header:"Set-Cookie"`
	Body   Sesion
}

type RespuestaSesionActual struct {
	Body Sesion
}

type RespuestaCierreSesion struct {
	Cookie string `header:"Set-Cookie"`
}

var sesionSimulada = Sesion{
	ID:       "4472ad91-292f-492b-8bba-99ebb11133b5",
	Correo:   "operador@base-tecnica.test",
	Rol:      "operador",
	Permisos: []string{},
}

func RegistrarRutasSesiones(apiHuma huma.API) {
	huma.Register(apiHuma, huma.Operation{
		OperationID: "iniciar-sesion",
		Method:      http.MethodPost,
		Path:        "/api/sesiones",
		Summary:     "Inicia una sesión",
		Tags:        []string{"Sesiones"},
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized},
	}, iniciarSesion)

	huma.Register(apiHuma, huma.Operation{
		OperationID: "consultar-sesion-actual",
		Method:      http.MethodGet,
		Path:        "/api/sesiones/actual",
		Summary:     "Consulta la sesión actual",
		Tags:        []string{"Sesiones"},
		Errors:      []int{http.StatusUnauthorized},
	}, consultarSesionActual)

	huma.Register(apiHuma, huma.Operation{
		OperationID:   "cerrar-sesion",
		Method:        http.MethodDelete,
		Path:          "/api/sesiones",
		Summary:       "Cierra la sesión",
		Tags:          []string{"Sesiones"},
		DefaultStatus: http.StatusNoContent,
	}, cerrarSesion)
}

func iniciarSesion(_ context.Context, entrada *EntradaInicioSesion) (*RespuestaInicioSesion, error) {
	if entrada.Body.Correo != sesionSimulada.Correo || entrada.Body.Clave != "demo" {
		return nil, huma.Error401Unauthorized("Correo o clave incorrectos")
	}

	return &RespuestaInicioSesion{
		Cookie: cookieSesion(valorSesionSimulada, 0),
		Body:   sesionSimulada,
	}, nil
}

func consultarSesionActual(_ context.Context, entrada *EntradaSesionActual) (*RespuestaSesionActual, error) {
	if entrada.ValorCookie != valorSesionSimulada {
		return nil, huma.Error401Unauthorized("No hay una sesión activa")
	}

	return &RespuestaSesionActual{Body: sesionSimulada}, nil
}

func cerrarSesion(_ context.Context, _ *struct{}) (*RespuestaCierreSesion, error) {
	// MaxAge negativo indica al navegador que elimine la cookie.
	return &RespuestaCierreSesion{Cookie: cookieSesion("", -1)}, nil
}

func cookieSesion(valor string, maxEdad int) string {
	cookie := http.Cookie{
		Name:     "sesion",
		Value:    valor,
		Path:     "/",
		MaxAge:   maxEdad,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie.String()
}
