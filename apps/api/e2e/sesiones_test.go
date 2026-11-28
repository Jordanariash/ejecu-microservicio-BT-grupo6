package e2e

import (
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"testing"
	"time"
)

func probarSesiones(t *testing.T, direccion string, reporte *strings.Builder) {
	t.Helper()

	cookiesCliente, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	cliente := http.Client{Jar: cookiesCliente, Timeout: 5 * time.Second}
	pedir := func(paso, metodo, ruta, cuerpo string) []*http.Cookie {
		t.Helper()
		peticion, err := http.NewRequest(metodo, direccion+ruta, strings.NewReader(cuerpo))
		if err != nil {
			t.Fatal(err)
		}
		if cuerpo != "" {
			peticion.Header.Set("Content-Type", "application/json")
		}
		respuesta, err := cliente.Do(peticion)
		if err != nil {
			t.Fatal(err)
		}
		defer respuesta.Body.Close()
		anotar(reporte, paso, strconv.Itoa(respuesta.StatusCode))
		return respuesta.Cookies()
	}

	const rutaSesiones = "/api/sesiones"
	const rutaSesionActual = rutaSesiones + "/actual"
	const credencialesInvalidas = `{"correo":"operador@base-tecnica.test","clave":"incorrecta"}`
	const credencialesValidas = `{"correo":"operador@base-tecnica.test","clave":"demo"}`

	pedir("Consulta sin sesión", http.MethodGet, rutaSesionActual, "")
	pedir("Inicio con clave incorrecta", http.MethodPost, rutaSesiones, credencialesInvalidas)
	cookies := pedir("Inicio de sesión", http.MethodPost, rutaSesiones, credencialesValidas)
	if len(cookies) != 1 {
		t.Fatalf("se esperaba una cookie de sesión, llegaron %d", len(cookies))
	}
	sesion := cookies[0]
	if sesion.Name != "sesion" || !sesion.HttpOnly || !sesion.Secure || sesion.SameSite != http.SameSiteLaxMode {
		t.Fatalf("atributos de la cookie de sesión incorrectos: %s", sesion)
	}

	pedir("Consulta con sesión activa", http.MethodGet, rutaSesionActual, "")
	pedir("Cierre de sesión", http.MethodDelete, rutaSesiones, "")
	pedir("Consulta después del cierre", http.MethodGet, rutaSesionActual, "")
}
