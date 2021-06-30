package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	jsonutils "github.com/cloudlifter/go-utils/json"
)

func (ah *StocksdashHander) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		ah.logger.Debug("validating access token")

		token, err := extractToken(r)
		if err != nil {
			ah.logger.Error("Token not provided or malformed")
			w.WriteHeader(http.StatusBadRequest)
			// datastore.ToJSON(&GenericError{Error: err.Error()}, w)
			jsonutils.ToJSON(&GenericResponse{Status: false, Message: "Authentication failed. Token not provided or malformed"}, w)
			return
		}
		ah.logger.Debug("token present in header", token)

		userID, err := ah.stocksService.ValidateAccessToken(token)
		if err != nil {
			ah.logger.Error("token validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			// datastore.ToJSON(&GenericError{Error: err.Error()}, w)
			jsonutils.ToJSON(&GenericResponse{Status: false, Message: "Authentication failed. Invalid token"}, w)
			return
		}
		ah.logger.Debug("access token validated")

		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("token not provided or malformed")
	}
	return authHeaderContent[1], nil
}
