package entity

// CustomClaims represents JWT claims structure compatible with serverless auth
// Based on the pattern from tc-golunch-api monolith
type CustomClaims struct {
	UserID   string                 `json:"user_id"`
	UserType string                 `json:"user_type"`
	Custom   map[string]interface{} `json:"custom,omitempty"`
	// Standard JWT fields
	ExpiresAt int64 `json:"exp,omitempty"`
	IssuedAt  int64 `json:"iat,omitempty"`
	NotBefore int64 `json:"nbf,omitempty"`
}