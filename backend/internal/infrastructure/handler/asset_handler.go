package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	assetusecase "smart-allocation/internal/application/usecase/asset"
	"smart-allocation/internal/application/usecase/asset/dto"
	domainerrors "smart-allocation/internal/domain/errors"
)

// AssetHandler defines the HTTP contract for asset operations.
type AssetHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByTicker(w http.ResponseWriter, r *http.Request)
	ListAssets(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type assetHandler struct {
	createUseCase assetusecase.CreateAssetUseCase
	getUseCase    assetusecase.GetAssetUseCase
	listUseCase   assetusecase.ListAssetsUseCase
	updateUseCase assetusecase.UpdateAssetUseCase
	deleteUseCase assetusecase.DeleteAssetUseCase
}

func NewAssetHandler(
	createUseCase assetusecase.CreateAssetUseCase,
	getUseCase assetusecase.GetAssetUseCase,
	listUseCase assetusecase.ListAssetsUseCase,
	updateUseCase assetusecase.UpdateAssetUseCase,
	deleteUseCase assetusecase.DeleteAssetUseCase,
) AssetHandler {
	return &assetHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func respondError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if domainerrors.IsNotFoundError(err) {
		status = http.StatusNotFound
	} else if domainerrors.IsValidationError(err) || domainerrors.IsAlreadyExistsError(err) {
		status = http.StatusBadRequest
	}

	ce := &domainerrors.CustomError{}
	if domainerrors.IsCustomError(err, ce) {
		respond(w, status, ce)
		return
	}
	respond(w, status, map[string]string{"error": err.Error()})
}

// parseTotalToInvest reads the query param ?total_to_invest=X.
// Returns 0 if absent or invalid.
func parseTotalToInvest(r *http.Request) float64 {
	s := r.URL.Query().Get("total_to_invest")
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// Create registers a new asset in the portfolio.
//
//	@Summary	Register a new asset
//	@Tags		assets
//	@Accept		json
//	@Produce	json
//	@Param		asset	body		dto.CreateAssetRequestDTO	true	"Asset data"
//	@Success	201		{object}	dto.CreateAssetResponseDTO
//	@Failure	400		{object}	domainerrors.CustomError
//	@Failure	500		{object}	domainerrors.CustomError
//	@Router		/assets [post]
func (h *assetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAssetRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	resp, err := h.createUseCase.Execute(r.Context(), &req)
	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusCreated, resp)
}

// GetByTicker fetches an asset by ticker.
//
//	@Summary	Fetch an asset by ticker
//	@Tags		assets
//	@Produce	json
//	@Param		ticker			path		string	true	"Asset ticker"						example(BBAS3)
//	@Param		total_to_invest	query		number	false	"Total amount available to invest"
//	@Success	200				{object}	dto.GetAssetResponseDTO
//	@Failure	404				{object}	domainerrors.CustomError
//	@Failure	500				{object}	domainerrors.CustomError
//	@Router		/assets/{ticker} [get]
func (h *assetHandler) GetByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(r.PathValue("ticker"))

	resp, err := h.getUseCase.Execute(r.Context(), ticker, parseTotalToInvest(r))
	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusOK, resp)
}

// ListAssets lists all assets in the portfolio with a rebalancing summary.
//
//	@Summary	List all assets in the portfolio
//	@Tags		assets
//	@Produce	json
//	@Param		total_to_invest	query		number	false	"Total amount available to invest"
//	@Success	200				{object}	dto.ListAssetsResponseDTO
//	@Failure	500				{object}	domainerrors.CustomError
//	@Router		/assets [get]
func (h *assetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	resp, err := h.listUseCase.Execute(r.Context(), parseTotalToInvest(r))
	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusOK, resp)
}

// Update updates an existing asset's data.
//
//	@Summary	Update an asset
//	@Tags		assets
//	@Accept		json
//	@Param		ticker	path	string						true	"Asset ticker"		example(BBAS3)
//	@Param		asset	body	dto.UpdateAssetRequestDTO	true	"Updated data"
//	@Success	204
//	@Failure	400	{object}	domainerrors.CustomError
//	@Failure	404	{object}	domainerrors.CustomError
//	@Failure	500	{object}	domainerrors.CustomError
//	@Router		/assets/{ticker} [put]
func (h *assetHandler) Update(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(r.PathValue("ticker"))

	var req dto.UpdateAssetRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	if err := h.updateUseCase.Execute(r.Context(), ticker, &req); err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusNoContent, nil)
}

// Delete removes an asset from the portfolio.
//
//	@Summary	Remove an asset from the portfolio
//	@Tags		assets
//	@Param		ticker	path	string	true	"Asset ticker"	example(BBAS3)
//	@Success	204
//	@Failure	404	{object}	domainerrors.CustomError
//	@Failure	500	{object}	domainerrors.CustomError
//	@Router		/assets/{ticker} [delete]
func (h *assetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(r.PathValue("ticker"))

	if err := h.deleteUseCase.Execute(r.Context(), ticker); err != nil {
		respondError(w, err)
		return
	}

	respond(w, http.StatusNoContent, nil)
}
