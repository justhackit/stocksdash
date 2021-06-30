package service

import (
	"errors"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
	configutils "github.com/justhackit/stocksdash/config"
)

type Stocksdash interface {
	ValidateAccessToken(token string) (string, error)
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}

type StockdashSvc struct {
	logger  hclog.Logger
	configs *configutils.Configurations
}

func NewStockdashSvc(l hclog.Logger, c *configutils.Configurations) *StockdashSvc {
	return &StockdashSvc{l, c}
}

func (stock *StockdashSvc) ValidateAccessToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			stock.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("unexpected signing method in auth token")
		}
		//verifyBytes, err := ioutil.ReadFile(auth.configs.AccessTokenConf.AccessTokenPublicKeyPath)
		verifyBytes, err := ioutil.ReadFile("/Users/AEdapa/Personal/justhackit_github/stocksdash/access-public.pem")
		if err != nil {
			stock.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			stock.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		stock.logger.Error("unable to parse claims", "error", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, nil
}
