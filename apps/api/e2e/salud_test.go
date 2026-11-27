package e2e

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var actualizar = flag.Bool("actualizar", false, "reescribe el archivo golden con el resultado actual")

const imagenPostgres = "postgres:17.11@sha256:ff87138949ad9a46f32ff23bd6852b52ce024a13dffe9b4087be98292cdd3b4f"

func TestSaludDelServicio(t *testing.T) {
	binario := os.Getenv("SERVICIO_BIN")
	if binario == "" {
		t.Skip("es una prueba e2e: se corre con pnpm test:e2e:api")
	}

	ctx := context.Background()
	baseDatos := levantarPostgres(t, ctx)
	urlBaseDatos, err := baseDatos.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	servicio, direccion, logs := levantarServicio(t, binario, urlBaseDatos)

	var reporte strings.Builder
	reporte.WriteString("# E2E: salud del servicio\n\n")
	anotar(&reporte, "GET /healthz con la base de datos arriba", pedir(t, direccion+"/healthz"))
	anotar(&reporte, "GET /readyz con la base de datos arriba", pedir(t, direccion+"/readyz"))

	err = baseDatos.Stop(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	anotar(&reporte, "GET /readyz con la base de datos detenida", pedir(t, direccion+"/readyz"))
	anotar(&reporte, "GET /healthz con la base de datos detenida", pedir(t, direccion+"/healthz"))

	anotar(&reporte, "Apagado con SIGTERM", apagar(t, servicio, logs))

	compararConGolden(t, reporte.String())
}

func levantarPostgres(t *testing.T, ctx context.Context) *postgres.PostgresContainer {
	baseDeDatos := postgres.WithDatabase("base_tecnica")
	usuario := postgres.WithUsername("prueba")
	clave := postgres.WithPassword("prueba")
	contenedor, err := postgres.Run(ctx, imagenPostgres, baseDeDatos, usuario, clave, postgres.BasicWaitStrategies())

	// Se registra antes de revisar el error para que el contenedor se borre aunque falle al arrancar
	testcontainers.CleanupContainer(t, contenedor)
	if err != nil {
		t.Fatal(err)
	}
	return contenedor
}

func levantarServicio(t *testing.T, binario string, urlBaseDatos string) (*exec.Cmd, string, *bytes.Buffer) {
	puerto := puertoLibre(t)
	direccion := fmt.Sprintf("http://127.0.0.1:%d", puerto)

	servicio := exec.Command(binario)
	servicio.Env = os.Environ()
	servicio.Env = append(servicio.Env, fmt.Sprintf("API_PORT=%d", puerto))
	servicio.Env = append(servicio.Env, "API_ENV=test")
	servicio.Env = append(servicio.Env, "DATABASE_URL="+urlBaseDatos)
	servicio.Env = append(servicio.Env, "FRONTEND_URL=http://localhost:5160")
	servicio.Env = append(servicio.Env, "TZ=America/Santiago")

	logs := &bytes.Buffer{}
	servicio.Stdout = logs
	servicio.Stderr = logs

	err := servicio.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { servicio.Process.Kill() })

	esperarQueResponda(t, direccion+"/healthz")
	return servicio, direccion, logs
}

func puertoLibre(t *testing.T) int {
	escucha, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer escucha.Close()
	return escucha.Addr().(*net.TCPAddr).Port
}

func esperarQueResponda(t *testing.T, url string) {
	for intento := 0; intento < 50; intento++ {
		respuesta, err := http.Get(url)
		if err == nil {
			respuesta.Body.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("el servicio no respondió en %s", url)
}

// Devuelve "código tipo cuerpo", por ejemplo: 200 application/json {"estado":"ok"}
func pedir(t *testing.T, url string) string {
	cliente := http.Client{Timeout: 5 * time.Second}
	respuesta, err := cliente.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer respuesta.Body.Close()

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		t.Fatal(err)
	}

	tipo := respuesta.Header.Get("Content-Type")
	return fmt.Sprintf("%d %s %s", respuesta.StatusCode, tipo, strings.TrimSpace(string(cuerpo)))
}

func apagar(t *testing.T, servicio *exec.Cmd, logs *bytes.Buffer) string {
	err := servicio.Process.Signal(syscall.SIGTERM)
	if err != nil {
		t.Fatal(err)
	}

	codigo := 0
	err = servicio.Wait()
	var errorSalida *exec.ExitError
	if errors.As(err, &errorSalida) {
		codigo = errorSalida.ExitCode()
	}

	registroApagado := "no registra la señal de término"
	if strings.Contains(logs.String(), "señal de término recibida") {
		registroApagado = "registra la señal de término"
	}

	return fmt.Sprintf("código de salida %d, %s", codigo, registroApagado)
}

func anotar(reporte *strings.Builder, paso string, resultado string) {
	fmt.Fprintf(reporte, "- %s: %s\n", paso, resultado)
}

func compararConGolden(t *testing.T, actual string) {
	ruta := filepath.Join("testdata", "salud.golden.md")

	if *actualizar {
		err := os.WriteFile(ruta, []byte(actual), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	esperado, err := os.ReadFile(ruta)
	if err != nil {
		t.Fatal(err)
	}

	if actual != string(esperado) {
		t.Errorf("el resultado no coincide con %s\n\n--- esperado:\n%s\n--- obtenido:\n%s", ruta, esperado, actual)
	}
}
