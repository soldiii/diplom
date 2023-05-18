package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
)

type contextKey string

const (
	ctxKeyID   contextKey = "userID"
	ctxKeyRole contextKey = "userRole"
)

const (
	agent      = 0
	supervisor = 1
	user       = 2
)

type IDAndRoleFromContext struct {
	ID   int
	Role string
}

func ParseFromContext(r *http.Request, role int) (*IDAndRoleFromContext, error) {
	ID := r.Context().Value(ctxKeyID)
	Role := r.Context().Value(ctxKeyRole)
	uRole, ok := Role.(string)
	if !ok {
		err := errors.New("ошибка в поле userRole токена")
		return nil, err
	}
	switch role {
	case 0:
		if uRole != "agent" && uRole != "Agent" {
			err := errors.New("userRole != agent || Agent")
			return nil, err
		}
	case 1:
		if uRole != "supervisor" && uRole != "Supervisor" {
			err := errors.New("userRole != supervisor || Supervisor")
			return nil, err
		}
	case 2:
		break
	}
	uID, ok := ID.(int)
	if !ok {
		err := errors.New("ошибка в поле userID токена")
		return nil, err
	}
	returnStruct := &IDAndRoleFromContext{ID: uID, Role: uRole}
	return returnStruct, nil
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

		claims, err := h.services.Authorization.ParseToken(headerParts[1], true)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		flag := h.services.Authorization.IsTokenExpired(claims.ExpiresAt)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		if flag {
			NewErrorResponse(w, http.StatusUnauthorized, "истекло время access-токена")
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyID, claims.UserID)
		ctx = context.WithValue(ctx, ctxKeyRole, claims.UserRole)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
