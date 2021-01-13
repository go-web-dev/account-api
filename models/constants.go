package models

const (
	// ContentType represents the Content-Type header key
	ContentType = "Content-Type"
	// ApplicationJSONType represents the application/json header value
	ApplicationJSONType = "application/json"

	GrantTypePassword          = "password"
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypeClientID          = "client_id"
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRtCookie          = "rt_cookie"
)
