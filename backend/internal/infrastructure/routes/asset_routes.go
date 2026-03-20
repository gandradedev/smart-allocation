package routes

import (
	"net/http"

	"smart-allocation/internal/infrastructure/handler"
)

func RegisterAssetRoutes(mux *http.ServeMux, h handler.AssetHandler) {
	mux.HandleFunc("GET /assets", h.ListAssets)
	mux.HandleFunc("POST /assets", h.Create)
	mux.HandleFunc("GET /assets/{ticker}", h.GetByTicker)
	mux.HandleFunc("PUT /assets/{ticker}", h.Update)
	mux.HandleFunc("DELETE /assets/{ticker}", h.Delete)
}
