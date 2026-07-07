package domain // package declaration for the module

type SignupRequest struct { // declare struct type
	Name     string `json:"name"` // execute statement
	Email    string `json:"email"` // execute statement
	Password string `json:"password"` // execute statement
} // end block

type LoginRequest struct { // declare struct type
	Email    string `json:"email"` // execute statement
	Password string `json:"password"` // execute statement
} // end block

type RefreshTokenRequest struct { // declare struct type
	RefreshToken string `json:"refresh_token"` // execute statement
} // end block

type LogoutRequest struct { // declare struct type
	RefreshToken string `json:"refresh_token"` // execute statement
} // end block

type AuthResponse struct { // declare struct type
	User         UserResponse `json:"user"` // execute statement
	AccessToken  string       `json:"access_token"` // execute statement
	RefreshToken string       `json:"refresh_token"` // execute statement
} // end block
