package loginspost

import (
	"github.com/bitterbattles/api/pkg/jwt"
	"github.com/bitterbattles/api/pkg/lambda/api"
	"github.com/bitterbattles/api/pkg/time"
)

const accessExpiresIn = 3600      // 1 hour
const refreshExpiresIn = 15768000 // 6 months

// CreateResponse creates a login response
func CreateResponse(userID string, accessTokenSecret string, refreshTokenSecret string) (*Response, error) {
	accessToken, err := createToken(userID, accessExpiresIn, accessTokenSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := createToken(userID, refreshExpiresIn, refreshTokenSecret)
	if err != nil {
		return nil, err
	}
	response := &Response{
		AccessToken:      accessToken,
		AccessExpiresIn:  accessExpiresIn,
		RefreshToken:     refreshToken,
		RefreshExpiresIn: refreshExpiresIn,
	}
	return response, nil
}

func createToken(userID string, expiresIn int, secret string) (string, error) {
	now := time.NowUnix()
	authContext := &api.AuthContext{
		UserID:    userID,
		CreatedOn: now,
		ExpiresOn: now + accessExpiresIn,
	}
	return jwt.NewHS256(authContext, secret)
}
