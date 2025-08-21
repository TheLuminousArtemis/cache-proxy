package main

import (
	"context"
)

func (app *application) clearCache() {
	app.logger.Info("Clearing cache")
	if err := app.cache.Del(context.Background(), key); err != nil {
		app.logger.Error("Failed to clear cache", "error", err)
	}
	app.logger.Info("Cache cleared")
}
