package main

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Carpeta donde el Dockerfile copia la web ya construida
const carpetaWeb = "/srv"

var backend *url.URL

func main() {
	configurarLogger()
	backend = leerBackend()

	rutas := http.NewServeMux()
	rutas.HandleFunc("GET /healthz", responderVivo)
	rutas.Handle("/api/", crearProxy())
	rutas.HandleFunc("/", servirWeb)

	servidor := &http.Server{Addr: ":8080", Handler: rutas, ReadHeaderTimeout: 5 * time.Second}

	slog.Info("servidor web iniciado", "puerto", "8080", "backend", backend.String())
	err := servidor.ListenAndServe()
	terminarConError(err.Error())
}

func configurarLogger() {
	manejador := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(manejador).With("servicio", "base-tecnica-grupo6-web")
	slog.SetDefault(logger)
}

func leerBackend() *url.URL {
	urlBackend := os.Getenv("BACKEND_URL")
	if urlBackend == "" {
		terminarConError("falta la variable de entorno obligatoria BACKEND_URL")
	}

	direccion, err := url.Parse(urlBackend)
	if err != nil {
		terminarConError("BACKEND_URL no es una URL válida")
	}
	return direccion
}

func responderVivo(w http.ResponseWriter, r *http.Request) {
	agregarCabeceras(w.Header())
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"estado":"ok"}`))
}

// Reenvía /api al backend: así el navegador ve un solo origen y la cookie de sesión funciona sin CORS
func crearProxy() *httputil.ReverseProxy {
	transporte := http.DefaultTransport.(*http.Transport).Clone()
	transporte.ResponseHeaderTimeout = 5 * time.Second

	proxy := &httputil.ReverseProxy{Transport: transporte}
	proxy.Rewrite = reenviarAlBackend
	proxy.ModifyResponse = agregarCabecerasAlBackend
	return proxy
}

func reenviarAlBackend(peticion *httputil.ProxyRequest) {
	peticion.SetURL(backend)
	peticion.SetXForwarded()
}

func agregarCabecerasAlBackend(respuesta *http.Response) error {
	agregarCabeceras(respuesta.Header)
	return nil
}

// Si la ruta no es un archivo, devuelve index.html para que la resuelva el router de React
func servirWeb(w http.ResponseWriter, r *http.Request) {
	agregarCabeceras(w.Header())

	archivos := os.DirFS(carpetaWeb)
	ruta := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if !esArchivo(archivos, ruta) {
		ruta = "index.html"
	}

	// Vite les pone un hash en el nombre a los archivos de assets/, así que nunca cambian
	if strings.HasPrefix(ruta, "assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	http.ServeFileFS(w, r, archivos, ruta)
}

func esArchivo(archivos fs.FS, ruta string) bool {
	informacion, err := fs.Stat(archivos, ruta)
	return err == nil && !informacion.IsDir()
}

// Cabeceras de seguridad que pide el DAS (sección 4)
func agregarCabeceras(cabeceras http.Header) {
	cabeceras.Set("X-Content-Type-Options", "nosniff")
	cabeceras.Set("Referrer-Policy", "strict-origin-when-cross-origin")
}

func terminarConError(mensaje string) {
	slog.Error("el servidor web terminó con error", "error", mensaje)
	os.Exit(1)
}
