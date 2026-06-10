package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"kasiraiai/backend/config"
	"kasiraiai/backend/internal/model"
	"kasiraiai/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrPhoneAlreadyExists adalah sentinel error untuk nomor WhatsApp yang sudah terdaftar.
var ErrPhoneAlreadyExists = errors.New("nomor WhatsApp sudah terdaftar")

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*model.Umkm, string, error)
	Login(ctx context.Context, phone, password string) (*model.Umkm, string, error)
	ValidateToken(tokenStr string) (uuid.UUID, error)
}

type RegisterRequest struct {
	Name         string `json:"name" binding:"required"`
	BusinessName string `json:"business_name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password" binding:"required,min=6"`
	BusinessType string `json:"business_type,omitempty"`
}

type authService struct {
	umkmRepo repository.UmkmRepo
	cfg      *config.Config
}

func NewAuthService(umkmRepo repository.UmkmRepo, cfg *config.Config) AuthService {
	return &authService{umkmRepo: umkmRepo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*model.Umkm, string, error) {
	existing, err := s.umkmRepo.FindByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, "", fmt.Errorf("auth_service.Register: check phone: %w", err)
	}
	if existing != nil {
		return nil, "", ErrPhoneAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, "", fmt.Errorf("auth_service.Register: hash password: %w", err)
	}

	u := &model.Umkm{
		Name:         req.Name,
		BusinessName: req.BusinessName,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(hash),
		IsActive:     true,
	}
	if req.Email != "" {
		u.Email = &req.Email
	}
	if req.BusinessType != "" {
		u.BusinessType = &req.BusinessType
	}

	if err := s.umkmRepo.Create(ctx, u); err != nil {
		return nil, "", fmt.Errorf("auth_service.Register: create: %w", err)
	}

	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	slog.Info("UMKM terdaftar", "umkm_id", u.ID, "phone", u.PhoneNumber)
	return u, token, nil
}

func (s *authService) Login(ctx context.Context, phone, password string) (*model.Umkm, string, error) {
	u, err := s.umkmRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, "", fmt.Errorf("auth_service.Login: find: %w", err)
	}
	if u == nil {
		return nil, "", fmt.Errorf("nomor WhatsApp atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("nomor WhatsApp atau password salah")
	}

	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	slog.Info("UMKM login", "umkm_id", u.ID)
	return u, token, nil
}

func (s *authService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak sesuai")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("token tidak valid")
	}

	idStr, ok := claims["umkm_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("klaim umkm_id tidak ditemukan")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("umkm_id tidak valid")
	}

	return id, nil
}

func (s *authService) generateToken(umkmID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"umkm_id": umkmID.String(),
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("auth_service.generateToken: %w", err)
	}
	return tokenStr, nil
}
