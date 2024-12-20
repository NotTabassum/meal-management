package serializer

type LoginReq struct {
	UserID string `json:"user_id"`
}

type JWT struct {
	UserID            string `json:"user_id"`
	AccessToken       string `json:"access_token"`
	AccessTokenExpiry int64  `json:"access_token_expiry"`
}
