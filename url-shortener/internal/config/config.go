package config // package declaration for the module

import ( // start import block
	"errors" // import package
	"os" // import package
	"strconv" // import package
	"time" // import package

	"github.com/joho/godotenv" // import package
) // end import block or block scope

type Config struct { // declare struct type
	AppEnv             string // execute statement
	Port               string // execute statement
	BaseURL            string // execute statement
	DatabaseURL        string // execute statement
	JWTAccessSecret    string // execute statement
	AccessTokenExpiry  time.Duration // execute statement
	RefreshTokenExpiry time.Duration // execute statement
	ShortCodeLength    int // execute statement
} // end block

func Load() (Config, error) { // declare function
	_ = godotenv.Load() // assign value

	accessExpiry, err := time.ParseDuration(getEnv("ACCESS_TOKEN_EXPIRY", "15m")) // declare and initialize variable
	if err != nil { // if condition
		return Config{}, err // return statement
	} // end block

	refreshExpiry, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h")) // declare and initialize variable
	if err != nil { // if condition
		return Config{}, err // return statement
	} // end block

	shortCodeLength, err := strconv.Atoi(getEnv("SHORT_CODE_LENGTH", "7")) // declare and initialize variable
	if err != nil { // if condition
		return Config{}, err // return statement
	} // end block

	cfg := Config{ // declare and initialize variable
		AppEnv:             getEnv("APP_ENV", "development"), // execute statement
		Port:               getEnv("PORT", "8080"), // execute statement
		BaseURL:            getEnv("BASE_URL", "http: // execute statement
		DatabaseURL:        os.Getenv("DATABASE_URL"), // execute statement
		JWTAccessSecret:    os.Getenv("JWT_ACCESS_SECRET"), // execute statement
		AccessTokenExpiry:  accessExpiry, // execute statement
		RefreshTokenExpiry: refreshExpiry, // execute statement
		ShortCodeLength:    shortCodeLength, // execute statement
	} // end block

	if cfg.DatabaseURL == "" { // if condition
		return Config{}, errors.New("DATABASE_URL is required") // return statement
	} // end block

	if cfg.JWTAccessSecret == "" || cfg.JWTAccessSecret == "change-this-access-secret" { // if condition
		return Config{}, errors.New("JWT_ACCESS_SECRET must be set to a strong secret") // return statement
	} // end block

	return cfg, nil // return statement
} // end block

func getEnv(key, fallback string) string { // declare function
	value := os.Getenv(key) // declare and initialize variable
	if value == "" { // if condition
		return fallback // return statement
	} // end block
	return value // return statement
} // end block
