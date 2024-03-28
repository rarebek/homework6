package tokens

import (
	"EXAM3/api-gateway/pkg/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTHandler struct {
	Sub       string
	Exp       string
	Role      string
	SignInKey string
	Log       logger.Logger
	Token     string
	Timeout   int
}

type CustomClaims struct {
	*jwt.Token
	Sub  string `json:"sub"`
	Exp  string `json:"exp"`
	Role string `json:"role"`
}

func (jwtHandler JWTHandler) GenerateAuthJWT() (access string, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
		rtClaims     jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Second * time.Duration(jwtHandler.Timeout)).Unix()
	claims["sub"] = jwtHandler.Sub
	claims["role"] = jwtHandler.Role
	access, err = accessToken.SignedString([]byte(jwtHandler.SignInKey))
	if err != nil {
		jwtHandler.Log.Error("cannot generate access token", logger.Error(err))
		return
	}

	rtClaims = refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = jwtHandler.Sub
	rtClaims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtHandler.Timeout)).Unix()
	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SignInKey))
	if err != nil {
		jwtHandler.Log.Error("cannot generate refresh token", logger.Error(err))
		return
	}

	return access, refresh, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SignInKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		jwtHandler.Log.Error("invalid jwt token")
		return nil, err
	}

	return claims, nil
}

func ExtractClaim(tokenStr string, signinKey []byte) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return signinKey, nil
	})

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
