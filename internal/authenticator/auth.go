package authenticator

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/satori/go.uuid"

	"github.com/Alena-Kurushkina/shortener/internal/logger"
	"github.com/Alena-Kurushkina/shortener/internal/sherr"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const tokenExp = time.Hour * 3

// TODO перенести в env
const secretKey = "secretkey"

// BuildJWTString создаёт токен и возвращает его в виде строки.
func buildJWTString(id uuid.UUID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: id,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func getUserID(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired || v.Errors == jwt.ValidationErrorSignatureInvalid {
			return uuid.Nil, sherr.ErrTokenInvalid
		}
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, sherr.ErrTokenInvalid
	}
	if claims.UserID == uuid.Nil {
		return uuid.Nil, sherr.ErrNoUserIDInToken
	}
	logger.Log.Infof("User token is valid")
	return claims.UserID, nil
}

func setNewTokenInCookie(w http.ResponseWriter, userID uuid.UUID) error {
	jwt, err := buildJWTString(userID)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: "token", Value: jwt, MaxAge: 0})
	return nil
}

// AuthMiddleware realises middleware for user authentication
func AuthMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie){
				logger.Log.Infof("No cookie in request, method %s", r.Method)

				if r.Method != http.MethodPost {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				userID := uuid.NewV4()

				err := setNewTokenInCookie(w, userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				q := r.URL.Query()
				q.Add("userUUID", userID.String())
				r.URL.RawQuery = q.Encode()

				logger.Log.Infof("New user was registered with id %s", userID)

				h.ServeHTTP(w, r)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		userID, err := getUserID(cookie.Value)

		if err != nil {
			switch {
			case errors.Is(err, sherr.ErrNoUserIDInToken):
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			case errors.Is(err, sherr.ErrTokenInvalid):
				logger.Log.Infof("Token invalid")

				userID = uuid.NewV4()

				errt := setNewTokenInCookie(w, userID)
				if errt != nil {
					http.Error(w, errt.Error(), http.StatusInternalServerError)
					return
				}

				logger.Log.Infof("New user was registered with id %s", userID)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		q := r.URL.Query()
		q.Add("userUUID", userID.String())
		r.URL.RawQuery = q.Encode()

		logger.Log.Infof("Got user id %s from token", userID)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(logFn)
}
