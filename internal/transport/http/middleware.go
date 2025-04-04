package http

import (
	"context"
	"encoding/json"
	"net/http"
	"quiz_backend_core/internal/dto"
	"quiz_backend_core/internal/service"
	"strconv"
)

func accessControlMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, Accept, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func fillContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdStr := r.Header.Get("X-User-ID")
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			response := ErrorResponse{
				Code:    codeFrom(dto.ErrInternalServerError), //TODO correct error
				Message: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(response.Code)
			json.NewEncoder(w).Encode(response)
			return
		}

		role := dto.Role(r.Header.Get("X-User-Role")) //TODO More strict check - error if no such header
		if role == "" {
			response := ErrorResponse{
				Code:    codeFrom(dto.ErrInternalServerError),
				Message: "role is empty", //TODO
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(response.Code)
			json.NewEncoder(w).Encode(response)
			return
			//return  dto.ErrBadRouting //TODO More strict check - error if no such header
		}

		r = r.WithContext(context.WithValue(r.Context(), service.ContextVariablesUserID, userId))
		r = r.WithContext(context.WithValue(r.Context(), service.ContextVariablesUserRole, role))

		h.ServeHTTP(w, r)
	})
}
