package authservice

import (
	"context"
	"taskmanager/repo"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServicer interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
}

type AuthService struct {
	userRepo  *repo.UserRepository
	jwtSecret string
	jwtIssuer string
	tokenTTL  time.Duration
}

type TokenManager struct {
	secretKey []byte
	issuer    string
}

func NewAuthService(repo *repo.UserRepository, secret, issuer string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: secret,
		jwtIssuer: issuer,
		tokenTTL:  time.Hour * 24,
	}
}
func NewTokenManager(secret string, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}
func (as *AuthService) Login(ctx context.Context, login, password string) (string, error) {

	user, err := as.userRepo.GetUserByLogin(login, ctx)
	if err != nil {
		return "", err
	}
	err = UnhashPassword(user.Password, password)
	if err != nil {
		return "", err
	}

	manager := NewTokenManager(as.jwtSecret, as.jwtIssuer)
	token, err := manager.GenerateToken(user.ID, as.tokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) Register(ctx context.Context, login, password string) error {
	hashedPassword := HashPassword(password)
	err := as.userRepo.CreateUser(login, hashedPassword, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *TokenManager) GenerateToken(userID int, duration time.Duration) (string, error) {
	// 1. Создаем Claims (полезную нагрузку)
	claims := jwt.MapClaims{
		"sub": userID,                          // Subject (кто это)
		"iss": m.issuer,                        // Issuer (кто выпустил)
		"exp": time.Now().Add(duration).Unix(), // Expiration (когда протухнет)
		"iat": time.Now().Unix(),               // Issued At (когда создан)
	}

	// 2. Выбираем алгоритм подписи (HS256 — стандарт)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Подписываем токен нашим секретным ключом
	return token.SignedString(m.secretKey)
}
func UnhashPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) string {
	Hashpas, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(Hashpas)
}
