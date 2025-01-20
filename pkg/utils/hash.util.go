package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type HashUtils interface {
	HashPassword(password string) (string, error)
	CheckHash(hashedPassword string, password string) bool
	GenerateJWT(id string) (*AuthJWT, error)
}

type hashUtils struct{}

/*
Funcao que criptografa a senha do usuário antes de ser salva no banco de dados
*/
func (u *hashUtils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

/*
Funcao que verifica a senha criptografada salva no banco de dados com a senha enviada pelo client
*/
func (u *hashUtils) CheckHash(hashedPassword string, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type TokenClaims struct {
	ID string
	jwt.StandardClaims
}

type JWT struct {
	Hash      string
	ExpiresIn time.Time
}

type AuthJWT struct {
	Token   JWT
	Refresh JWT
}

/*
Funcao que cria um token JWT para autorizacao de acesso às APIs protegidas
*/
func (u *hashUtils) GenerateJWT(id string) (*AuthJWT, error) {

	expireToken := time.Now().Add(time.Hour * 24)

	expireRefreshToken := time.Now().Add(time.Hour * 24 * 90)

	claims := TokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken.Unix(),
			Issuer:    "viniblima-auth",
			Subject:   "auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errToken := token.SignedString([]byte(os.Getenv("PASSWORD_SECRET")))

	if err := errToken; err != nil {
		return nil, err
	}
	//Refresh
	claimsRefresh := TokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireRefreshToken.Unix(),
			Issuer:    "viniblima-auth",
			Subject:   "refresh",
		},
	}
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	tokenStringRefresh, errRefresh := tokenRefresh.SignedString([]byte(os.Getenv("PASSWORD_SECRET")))

	if err := errRefresh; err != nil {
		return nil, err
	}

	return &AuthJWT{
		Token: JWT{
			tokenString,
			expireToken,
		},
		Refresh: JWT{
			tokenStringRefresh,
			expireRefreshToken,
		},
	}, nil
}

func NewHashUtils() HashUtils {
	return &hashUtils{}
}
