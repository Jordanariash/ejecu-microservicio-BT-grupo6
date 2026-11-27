package main

import (
	"fmt"
	"os"
	"strings"
)

type Configuracion struct {
	Puerto, Entorno, URLBaseDatos, URLFrontend, ZonaHoraria string
}

func leerConfiguracion() (Configuracion, error) {
	var configuracion Configuracion
	configuracion.Puerto = os.Getenv("API_PORT")
	configuracion.Entorno = os.Getenv("API_ENV")
	configuracion.URLBaseDatos = os.Getenv("DATABASE_URL")
	configuracion.URLFrontend = os.Getenv("FRONTEND_URL")
	configuracion.ZonaHoraria = os.Getenv("TZ")

	faltantes := variablesFaltantes(configuracion)
	if len(faltantes) > 0 {
		return Configuracion{}, fmt.Errorf("faltan variables de entorno obligatorias: %s", strings.Join(faltantes, ", "))
	}

	if !entornoValido(configuracion.Entorno) {
		return Configuracion{}, fmt.Errorf("API_ENV debe ser dev, test o prod, y llegó %q", configuracion.Entorno)
	}

	return configuracion, nil
}

func variablesFaltantes(configuracion Configuracion) []string {
	var faltantes []string

	if configuracion.Puerto == "" {
		faltantes = append(faltantes, "API_PORT")
	}
	if configuracion.Entorno == "" {
		faltantes = append(faltantes, "API_ENV")
	}
	if configuracion.URLBaseDatos == "" {
		faltantes = append(faltantes, "DATABASE_URL")
	}
	if configuracion.URLFrontend == "" {
		faltantes = append(faltantes, "FRONTEND_URL")
	}
	if configuracion.ZonaHoraria == "" {
		faltantes = append(faltantes, "TZ")
	}

	return faltantes
}

func entornoValido(entorno string) bool {
	return entorno == "dev" || entorno == "test" || entorno == "prod"
}
