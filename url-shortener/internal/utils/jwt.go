package utils // package declaration for the module

import ( // start import block
	"errors" // import package
	"time" // import package

	"github.com/golang-jwt/jwt/v5" // import package
	"github.com/google/uuid" // import package
) // end import block or block scope

type JWTManager struct { // declare struct type
	secret []byte // execute statement
	expiry time.Duration // execute statement
} // end block

type AccessTokenClaims struct { // declare struct type
	UserID uuid.UUID `json:"user_id"` // execute statement
	Email  string    `json:"email"` // execute statement
	jwt.RegisteredClaims // execute statement
} // end block

func NewJWTManager(secret string, expiry time.Duration) *JWTManager { // declare function
	return &JWTManager{secret: []byte(secret), expiry: expiry} // return statement
} // end block

func (m *JWTManager) Generate(userID uuid.UUID, email string) (string, error) { // declare method
	now := time.Now() // declare and initialize variable
	claims := AccessTokenClaims{ // declare and initialize variable
		UserID: userID, // execute statement
		Email:  email, // execute statement
		RegisteredClaims: jwt.RegisteredClaims{ // execute statement
			Subject:   userID.String(), // execute statement
			IssuedAt:  jwt.NewNumericDate(now), // execute statement
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)), // execute statement
		}, // close block
	} // end block

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // declare and initialize variable
	return token.SignedString(m.secret) // return statement
} // end block

func (m *JWTManager) Parse(tokenString string) (*AccessTokenClaims, error) { // declare method
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) { // declare and initialize variable
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // if condition
			return nil, errors.New("unexpected signing method") // return statement
		} // end block
		return m.secret, nil // return statement
	}) // close block
	if err != nil { // if condition
		return nil, err // return statement
	} // end block

	claims, ok := token.Claims.(*AccessTokenClaims) // declare and initialize variable
	if !ok || !token.Valid { // if condition
		return nil, errors.New("invalid token") // return statement
	} // end block

	return claims, nil // return statement
} // end block
