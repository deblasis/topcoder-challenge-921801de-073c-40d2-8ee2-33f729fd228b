// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
//
package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/common/config"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtHandler struct {
	jwtConfig config.JWTConfig
	logger    log.Logger

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJwtHandler(logger log.Logger, cfg config.JWTConfig) *JwtHandler {

	var privK *rsa.PrivateKey
	var pubK *rsa.PublicKey

	if cfg.PrivKeyPath != "" {
		privBytes, err := ioutil.ReadFile(cfg.PrivKeyPath)
		if err != nil {
			level.Debug(logger).Log("err", err)
		}

		privK, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
		if err != nil {
			level.Debug(logger).Log("err", err)
		}
	}

	if cfg.PubKeyPath != "" {
		pubBytes, err := ioutil.ReadFile(cfg.PubKeyPath)
		if err != nil {
			level.Error(logger).Log("err", err)
		}

		pubK, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
		if err != nil {
			level.Error(logger).Log("err", err)
		}
	}

	return &JwtHandler{
		jwtConfig:  cfg,
		logger:     logger,
		privateKey: privK,
		publicKey:  pubK,
	}
}

type jwtTokenClaims struct {
	Token  string
	Claims STCClaims
}

func (h *JwtHandler) NewJWTToken(userId uuid.UUID, username, role, issuer string) (tokenClaims jwtTokenClaims, expiresAt int64, err error) {
	now := time.Now().Unix()
	expiresAt = now + int64(h.jwtConfig.TokenDuration)
	claims := STCClaims{
		UserId:   userId.String(),
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Audience:  "deblasis.SpaceTrafficControl",
			ExpiresAt: expiresAt,
			IssuedAt:  now,
			Issuer:    issuer,
		},
	}

	token, err := h.generateToken(h.jwtConfig, claims)

	return jwtTokenClaims{
		Token:  token,
		Claims: claims,
	}, expiresAt, err
}

func (h *JwtHandler) ExtractTokenFromHTTPRequest(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (h *JwtHandler) VerifyToken(tokenString string) (*jwt.Token, error) {
	//tokenString := h.ExtractToken(r)
	token, err := jwt.ParseWithClaims(tokenString, &STCClaims{}, func(token *jwt.Token) (interface{}, error) {
		return h.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (h *JwtHandler) ValidateToken(tokenString string, validatorFn func(*jwt.Token) error) error {
	token, err := h.VerifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(STCClaims); !ok && !token.Valid {
		return err
	}

	err = defaultValidator(token)
	if err != nil {
		return err
	}

	if validatorFn != nil {
		err = validatorFn(token)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *JwtHandler) ExtractClaims(tokenString string) (*STCClaims, error) {
	token, err := h.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(STCClaims)
	if ok && token.Valid {
		return &claims, nil
	}
	return nil, errors.New("cannot extract claims from token")

}

func (h *JwtHandler) generateToken(config config.JWTConfig, claims STCClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(h.privateKey)

	return tokenString, err
}

func fatal(err error) {
	if err != nil {
		level.Error(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))).Log("err", err)
		panic(err)
	}
}
