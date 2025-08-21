package main

import (
	"encoding/json"
	"net/http"

	"github.com/theluminousartemis/caching-proxy/internal/models"
)

func (app *application) getProducts(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(app.config.Origin)
	if err != nil {
		app.logger.Error("Failed to get products", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	var Products *models.Products
	if err := json.NewDecoder(res.Body).Decode(&Products); err != nil {
		app.logger.Error("Failed to unmarshal products", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.cache.Set(r.Context(), Products)
	WriteJSON(w, http.StatusOK, Products)
}
