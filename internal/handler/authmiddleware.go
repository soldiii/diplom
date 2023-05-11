package handler

import (
	"net/http"
	"strings"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
)

type tokenClaims struct {
	UserID   int    `json:"userID"`
	UserRole string `json:"userRole"`
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(AUTHORIZATION_HEADER)
		if tokenString == "" {
			NewErrorResponse(w, http.StatusUnauthorized, "пустой auth header")
			return
		}
		headerParts := strings.Split(tokenString, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			NewErrorResponse(w, http.StatusUnauthorized, "неправильный auth header")
			return
		}

		if len(headerParts[1]) == 0 {
			NewErrorResponse(w, http.StatusUnauthorized, "токен пуст")
			return
		}

		/*token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		})*/

	}

}
