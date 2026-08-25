package main

import (
	"context"
	store "go-social/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	ctx := context.Background()

	user, err := app.store.Users.GetById(ctx, userId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, err)

		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
