package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/theluminousartemis/caching-proxy/internal/models"
)

var key string = "products"

func (app *application) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.getFromCache(r.Context())
	if err != nil || products == nil {
		if err != nil {
			app.logger.Error("error fetching from redis", "error", err)
		} else if products == nil {
			app.logger.Info("cache empty, fetching from origin", "origin", app.config.Origin)
		}
		products = app.fetchFromOrigin(w)
		app.cache.Set(r.Context(), products)
		w.Header().Set("X-Cache", "MISS")
	} else {
		w.Header().Set("X-Cache", "HIT")
	}
	WriteJSON(w, http.StatusOK, products)
}

func (app *application) getFromCache(ctx context.Context) (*models.Products, error) {
	return app.cache.Get(ctx, key)
}

func (app *application) fetchFromOrigin(w http.ResponseWriter) *models.Products {
	res, err := http.Get(app.config.Origin)
	if err != nil {
		app.logger.Error("Failed to get products", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
	defer res.Body.Close()
	var products *models.Products
	if err := json.NewDecoder(res.Body).Decode(&products); err != nil {
		app.logger.Error("Failed to unmarshal products", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
	return products
}
