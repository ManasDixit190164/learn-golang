package middleware // package declaration for the module

import ( // start import block
	"context" // import package
	"net/http" // import package
	"strings" // import package

	"github.com/google/uuid" // import package
	"github.com/manasdixit/url-shortener/internal/utils" // import package
	"github.com/manasdixit/url-shortener/pkg/response" // import package
) // end import block or block scope

type contextKey string // declare custom type

const userIDKey contextKey = "userID" // declare constant

func Auth(jwtManager *utils.JWTManager) func(http.Handler) http.Handler { // declare function
	return func(next http.Handler) http.Handler { // return statement
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // return statement
			authHeader := r.Header.Get("Authorization") // declare and initialize variable
			if authHeader == "" { // if condition
				response.Error(w, http.StatusUnauthorized, "authorization header is required") // send an error response
				return // return statement
			} // end block

			parts := strings.SplitN(authHeader, " ", 2) // declare and initialize variable
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") { // if condition
				response.Error(w, http.StatusUnauthorized, "invalid authorization header") // send an error response
				return // return statement
			} // end block

			claims, err := jwtManager.Parse(parts[1]) // verify JWT token
			if err != nil { // if condition
				response.Error(w, http.StatusUnauthorized, "invalid or expired token") // send an error response
				return // return statement
			} // end block

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID) // declare and initialize variable
			next.ServeHTTP(w, r.WithContext(ctx)) // execute statement
		}) // close block
	} // end block
} // end block

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) { // declare function
	userID, ok := ctx.Value(userIDKey).(uuid.UUID) // declare and initialize variable
	return userID, ok // return statement
} // end block
