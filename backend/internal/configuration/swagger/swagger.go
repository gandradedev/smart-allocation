package swagger

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwaggerRoutes registra a rota do Swagger UI no mux.
// Disponível em: GET /swagger/
func RegisterSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
}
