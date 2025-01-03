package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/anaabdi/warga-app-go/cmd/app/config"
	pkgerr "github.com/anaabdi/warga-app-go/internal/api/error"
	"github.com/anaabdi/warga-app-go/internal/api/parser"
	"github.com/anaabdi/warga-app-go/pkg/constant"
	pkgjwt "github.com/anaabdi/warga-app-go/pkg/jwt"
)

type AuthMiddleware struct {
	responder parser.JSONResponder
	cfg       *config.Config
}

func NewAuthMiddleware(responder parser.JSONResponder, cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		responder: responder,
		cfg:       cfg,
	}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

		authorizationHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			if len(fields) == 2 {
				accessToken = fields[1]
			}
		}

		if accessToken == "" {
			m.responder.Error(w, pkgerr.ErrUnauthorized)
			return
		}

		claims, err := pkgjwt.ValidateToken(accessToken, m.cfg.Auth.AccessTokenPublicKey)
		if err != nil {
			m.responder.Error(w, pkgerr.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constant.ContextAccountIDKey, fmt.Sprint(claims["sub"]))
		ctx = context.WithValue(ctx, constant.ContextRequestScope, fmt.Sprint(claims["scope"]))

		chainedNextRequestID := fmt.Sprint(claims["nxtrid"])
		if chainedNextRequestID != "" {
			ctx = context.WithValue(ctx, constant.ContextChainNextRequestIDKey, chainedNextRequestID)
		}

		accType := fmt.Sprint(claims["acc_type"])
		splitted := strings.Split(accType, ":")

		ctx = context.WithValue(ctx, constant.ContextAccountTypeKey, splitted[0])
		if len(splitted) > 1 {
			ctx = context.WithValue(ctx, constant.ContextAccountRoleKey, splitted[1])
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
