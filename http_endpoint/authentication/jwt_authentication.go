package httpAuth

import (
	"context"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
	"visualization-client/http_endpoint/common"
)

const contextJWTProperty = "AuthToken"

// CustomClaims defines what data would be stored in jwt token
type CustomClaims struct {
	IsAdmin   bool   `json:"isAdmin"`
	ProjectID string `json:"orgId"`
	jwt.StandardClaims
}

// JWTTokenFromParams creates jwt token given claims data
func JWTTokenFromParams(secret string, isAdmin bool, projectID string,
	expiresAt time.Time) (string, error) {

	claims := CustomClaims{
		isAdmin,
		projectID,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Error signing JWTToken %s", err)
	}
	return result, err
}

func parseJWTTokenClaims(tokenString, secret string) (*CustomClaims, error) {
	tokenValue, err := jwt.ParseWithClaims(
		tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

	if claims, ok := tokenValue.Claims.(*CustomClaims); ok && tokenValue.Valid {
		return claims, nil
	}
	return nil, err
}

// AdminAuthenticationMiddleware checks admin field in jwt token
func AdminAuthenticationMiddleware(secret string) func(http.Handler) http.Handler {

	/*
		Check if stored JWT token has 'IsAdmin' field set to true
	*/

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			storedToken := r.Context().Value(contextJWTProperty)
			claims, err := parseJWTTokenClaims(storedToken.(*jwt.Token).Raw, secret)
			if err != nil || !claims.IsAdmin {
				common.WriteErrorToResponse(w, http.StatusUnauthorized,
					http.StatusText(http.StatusUnauthorized),
					"User has no admin role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuthenticationMiddleware checks if provided jwt token is valid and not expired
func AuthenticationMiddleware(secret string) func(http.Handler) http.Handler {

	/*
	   go-jwt-middleware does not provide support for chi framework, that's
	   why this wrapper was implemented.

	   at the time of middleware initialization - jwtmiddleware.JWTMiddleware
	   struct is initialized in wrapped context. authenticationMiddleware
	   returns function compatible with chi framework interface. this function
	   uses JWTMiddleware functions to validate JWT token.
	*/

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  contextJWTProperty,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if jwtMiddleware.CheckJWT(w, r) != nil {
				// 401 code and error message is set in http response by
				// jwtMiddleware, we have just to return response
				return
			}
			storedToken := r.Context().Value(contextJWTProperty)
			claims, err := parseJWTTokenClaims(storedToken.(*jwt.Token).Raw, secret)
			if err == nil {
				newRequest := r.WithContext(context.WithValue(r.Context(),
					common.OrganizationIDContext, claims.ProjectID))
				*r = *newRequest
			}
			next.ServeHTTP(w, r)
		})
	}
}
