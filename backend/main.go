// @title          Smart Allocation API
// @version        1.0
// @description    REST API for managing and rebalancing a variable income investment portfolio.

// @contact.name   Gabriel Andrade

// @host      localhost:8080
// @BasePath  /

package main

import (
	"log"
	"net/http"

	_ "smart-allocation/docs"
	"smart-allocation/internal/application/usecase/asset"
	"smart-allocation/internal/configuration/database"
	"smart-allocation/internal/configuration/swagger"
	"smart-allocation/internal/infrastructure/handler"
	infrarepo "smart-allocation/internal/infrastructure/repository"
	"smart-allocation/internal/infrastructure/routes"
)

func main() {
	db, err := database.New("portfolio.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	repo := infrarepo.NewAssetRepository(db)

	createUC := asset.NewCreateAssetUseCase(repo)
	getUC := asset.NewGetAssetUseCase(repo)
	listUC := asset.NewListAssetsUseCase(repo)
	updateUC := asset.NewUpdateAssetUseCase(repo)
	deleteUC := asset.NewDeleteAssetUseCase(repo)

	h := handler.NewAssetHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mux := http.NewServeMux()
	routes.RegisterAssetRoutes(mux, h)
	swagger.RegisterSwaggerRoutes(mux)

	log.Println("Server running at http://localhost:8080")
	log.Println("Swagger UI at http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
