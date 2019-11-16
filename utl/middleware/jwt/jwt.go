package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"net/http"
	"strings"
	"time"
)

// New generates new JWT service necessery for auth middleware
func New(secret, algo string, d int) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	return &Service{
		key:      []byte(secret),
		algo:     signingMethod,
		duration: time.Duration(d) * time.Minute,
	}
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	duration time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// MWFunc makes JWT implement the Middleware interface.
func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseToken(c)
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)

			id := int64(claims["id"].(float64))
			name := claims["n"].(string)
			email := claims["e"].(string)
			role := claims["r"].(string)

			c.Set("id", id)
			c.Set("name", name)
			c.Set("email", email)
			c.Set("role", role)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (j *Service) ParseToken(c echo.Context) (*jwt.Token, error) {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, entity.ErrGeneric
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, entity.ErrGeneric
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if j.algo != token.Method {
			return nil, entity.ErrGeneric
		}
		return j.key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (j *Service) GenerateToken(u *model.Users) (string, string, error) {
	expire := time.Now().Add(j.duration)

	token := jwt.NewWithClaims((j.algo), jwt.MapClaims{
		"id":  u.Id,
		"e":   u.Email,
		"r":   u.AccessLevel,
		"n":   u.Name,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.key)

	return tokenString, expire.Format(time.RFC3339), err
}
