package middleware

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"pkg.agungdwiprasetyo.com/candi/candihelper"
	"pkg.agungdwiprasetyo.com/candi/candishared"
	"pkg.agungdwiprasetyo.com/candi/wrapper"

	"github.com/labstack/echo"
)

const (
	// Bearer constanta
	Bearer = "bearer"
)

// Bearer token validator
func (m *Middleware) Bearer(ctx context.Context, tokenString string) (*candishared.TokenClaim, error) {
	tokenClaim, err := m.tokenValidator.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return tokenClaim, nil
}

// HTTPBearerAuth http jwt token middleware
func (m *Middleware) HTTPBearerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authorization := c.Request().Header.Get(echo.HeaderAuthorization)
			tokenValue, err := extractAuthType(Bearer, authorization)
			if err != nil {
				return wrapper.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c.Response())
			}

			tokenClaim, err := m.Bearer(c.Request().Context(), tokenValue)
			if err != nil {
				return wrapper.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c.Response())
			}

			c.Set(candihelper.TokenClaimKey, tokenClaim)
			return next(c)
		}
	}
}

// GraphQLBearerAuth for graphql resolver
func (m *Middleware) GraphQLBearerAuth(ctx context.Context) *candishared.TokenClaim {
	headers := ctx.Value(candishared.ContextKey("headers")).(http.Header)
	authorization := headers.Get(echo.HeaderAuthorization)

	authValues := strings.Split(authorization, " ")
	authType := strings.ToLower(authValues[0])
	if authType != Bearer || len(authValues) != 2 {
		panic("Invalid authorization")
	}

	tokenClaim, err := m.Bearer(ctx, authValues[1])
	if err != nil {
		panic(err)
	}

	return tokenClaim
}

// GRPCBearerAuth method
func (m *Middleware) GRPCBearerAuth(ctx context.Context) *candishared.TokenClaim {

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		panic("missing context metadata")
	}

	authorizationMap := meta["authorization"]
	if len(authorizationMap) != 1 {
		panic(grpc.Errorf(codes.Unauthenticated, "Invalid authorization"))
	}

	authorization := authorizationMap[0]
	tokenClaim, err := m.Bearer(ctx, authorization)
	if err != nil {
		panic(err)
	}

	return tokenClaim
}