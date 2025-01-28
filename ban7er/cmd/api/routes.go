package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/keys/generate", app.GenerateEncryptionKeysHandler)
	mux.HandleFunc("/api/v1/keys/public", app.GetPublicEncryptionKeyHandler)
	mux.HandleFunc("/api/v1/keys/private", app.GetPrivateEncryptionKeyHandler)

	mux.HandleFunc("/api/v1/ws", app.HandleWebSocket)

	mux.HandleFunc("/", app.notFoundResponse)

	return app.metrics(
		app.recoverPanic(
			app.enableCORS(
				app.rateLimit(
					app.validateAPIHeaders(
						mux,
					),
				),
			),
		),
	)
}
