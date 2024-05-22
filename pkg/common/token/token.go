package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/tools/errs"
	"time"
)

type claims struct {
	UserID string
	jwt.RegisteredClaims
}

type Token struct {
	Expires time.Duration
	Secret  string
}

func New(expire int, secret string) *Token {
	return &Token{
		Expires: time.Duration(expire) * time.Hour * 24,
		Secret:  secret,
	}
}

func (t *Token) secret() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(t.Secret), nil
	}
}

func (t *Token) buildClaims(userID string) claims {
	now := time.Now()
	return claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(t.Expires)),    // Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                   // Issuing time
			NotBefore: jwt.NewNumericDate(now.Add(-time.Minute)), // Begin Effective time
		},
	}
}

func (t *Token) getToken(str string) (string, error) {
	token, err := jwt.ParseWithClaims(str, &claims{}, t.secret())
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", errs.ErrTokenMalformed.Wrap()
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errs.ErrTokenExpired.Wrap()
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", errs.ErrTokenNotValidYet.Wrap()
			} else {
				return "", errs.ErrTokenUnknown.Wrap()
			}
		} else {
			return "", errs.ErrTokenNotValidYet.Wrap()
		}
	} else {
		claims, ok := token.Claims.(*claims)
		if ok && token.Valid {
			return claims.UserID, nil
		}
		return "", errs.ErrTokenNotValidYet.Wrap()
	}
}

func (t *Token) CreateToken(UserID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.buildClaims(UserID))
	str, err := token.SignedString([]byte(t.Secret))
	if err != nil {
		return "", errs.Wrap(err)
	}
	return str, nil
}

func (t *Token) GetToken(token string) (string, error) {
	userID, err := t.getToken(token)
	if err != nil {
		return "", err
	}
	return userID, nil
}
