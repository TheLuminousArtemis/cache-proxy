package main

import (
	"context"
	"io"
	"net/http"
)

func (app *application) getContent(w http.ResponseWriter, r *http.Request) {
	res, err := app.getFromCache(r.Context())
	if err != nil || res == nil {
		if err != nil {
			app.logger.Error("error fetching from redis", "error", err)
		} else if res == nil {
			app.logger.Info("cache empty, fetching from origin", "origin", app.config.Origin)
		}
		resp := app.fetchFromOrigin(w)
		app.cache.Set(r.Context(), app.config.Origin, resp)
		w.Header().Set("X-Cache", "MISS")
		WriteJSON(w, http.StatusOK, map[string]any{"Response": string(resp)})
	} else {
		w.Header().Set("X-Cache", "HIT")
		WriteJSON(w, http.StatusOK, map[string]any{"Response": string(res)})
	}
}

func (app *application) getFromCache(ctx context.Context) ([]byte, error) {
	return app.cache.Get(ctx, app.config.Origin)
}

func (app *application) fetchFromOrigin(w http.ResponseWriter) []byte {
	res, err := http.Get(app.config.Origin)
	if err != nil {
		app.logger.Error("Failed to get content", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		app.logger.Error("Failed to read response", "error", err)
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
	return data
}
