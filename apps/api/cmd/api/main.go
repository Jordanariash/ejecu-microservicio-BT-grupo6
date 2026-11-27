package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"base-tecnica-grupo6/internal/api"
	"base-tecnica-grupo6/internal/service"
	"base-tecnica-grupo6/internal/store"
)

// Se reemplaza al compilar con -ldflags (ver Dockerfile).
var version = "dev"

func main() {
	slog.SetDefault(crearLogger())

	// "api openapi" imprime el contrato y termina. Lo usa pnpm codegen.
	if len(os.Args) > 1 && os.Args[1] == "openapi" {
		err := imprimirContrato()
		terminarSiHayError(err)
		return
	}

	err := ejecutarServidor()
	terminarSiHayError(err)
}

func ejecutarServidor() error {
	configuracion, err := leerConfiguracion()
	if err != nil {
		return err
	}

	ctx, detener := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer detener()

	baseDatos, err := store.Conectar(ctx, configuracion.URLBaseDatos)
	if err != nil {
		return err
	}
	defer baseDatos.Cerrar()

	salud := service.NuevaSalud(baseDatos)

	rutas := http.NewServeMux()
	crearAPI(rutas, configuracion.Entorno, salud)

	servidor := &http.Server{Addr: ":" + configuracion.Puerto, Handler: rutas, ReadHeaderTimeout: 5 * time.Second}

	erroresServidor := make(chan error, 1)
	go func() {
		slog.Info("servidor escuchando", "puerto", configuracion.Puerto, "entorno", configuracion.Entorno)
		erroresServidor <- servidor.ListenAndServe()
	}()

	select {
	case err := <-erroresServidor:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		slog.Info("señal de término recibida, apagando")
	}

	return apagarServidor(servidor)
}

func crearAPI(rutas *http.ServeMux, entorno string, salud *service.Salud) huma.API {
	opcionesHuma := huma.DefaultConfig("Base técnica · Grupo 6", version)
	opcionesHuma.OpenAPIPath = "/api/openapi"
	opcionesHuma.DocsPath = "/api/docs"
	// Sin el campo "$schema" en las respuestas
	opcionesHuma.SchemasPath = ""
	opcionesHuma.CreateHooks = nil

	// En prod no se publica la documentación (DAS sección 5)
	if entorno == "prod" {
		opcionesHuma.OpenAPIPath = ""
		opcionesHuma.DocsPath = ""
	}

	apiHuma := humago.New(rutas, opcionesHuma)
	api.RegistrarRutasSalud(apiHuma, salud)
	return apiHuma
}

func imprimirContrato() error {
	// Sin base de datos: solo se arman las rutas para describirlas
	apiHuma := crearAPI(http.NewServeMux(), "dev", nil)

	contrato, err := apiHuma.OpenAPI().YAML()
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(contrato)
	return err
}

func apagarServidor(servidor *http.Server) error {
	ctx, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()

	return servidor.Shutdown(ctx)
}

func crearLogger() *slog.Logger {
	manejador := slog.NewJSONHandler(os.Stdout, nil)
	return slog.New(manejador).With("servicio", "base-tecnica-grupo6", "version", version)
}

func terminarSiHayError(err error) {
	if err != nil {
		slog.Error("el servicio terminó con error", "error", err)
		os.Exit(1)
	}
}
