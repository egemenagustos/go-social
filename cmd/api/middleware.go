package main

import (
	"context"
	"encoding/base64"
	"fmt"
	store "go-social/internal/storage"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedErrorReponse(w, r, fmt.Errorf("authorization header is missing!"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if parts[0] != "Bearer" {
			app.unauthorizedErrorReponse(w, r, fmt.Errorf("authorization header is malformed!"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorReponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userId, ok := claims["sub"]
		if !ok {
			app.unauthorizedErrorReponse(w, r, fmt.Errorf("invalid token claims!"))
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId.(string))
		if err != nil {
			app.unauthorizedErrorReponse(w, r, fmt.Errorf("invalid token claims!"))
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) AuthMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				app.unauthorizedBasicErrorReponse(w, r, fmt.Errorf("authorization header is missing!"))
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorReponse(w, r, fmt.Errorf("authorization header is malformed!"))
				return
			}

			decoed, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorReponse(w, r, err)
				return
			}

			creds := strings.SplitN(string(decoed), ":", 2)

			username := app.config.auth.basic.user
			password := app.config.auth.basic.password

			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicErrorReponse(w, r, fmt.Errorf("invalid credentials!"))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func (app *application) CheckPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := getUserFromContext(r)
		post := getPostFromCtx(r)

		if post.UserId == user.Id {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePredence(r.Context(), user, role)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePredence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level > role.Level, nil
}
