package repository

import (
	"context"
	"fmt"

	"kasiraiai/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UmkmRepo interface {
	Create(ctx context.Context, u *model.Umkm) error
	FindByPhone(ctx context.Context, phone string) (*model.Umkm, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Umkm, error)
}

type umkmRepo struct {
	pool *pgxpool.Pool
}

func NewUmkmRepo(pool *pgxpool.Pool) UmkmRepo {
	return &umkmRepo{pool: pool}
}

func (r *umkmRepo) Create(ctx context.Context, u *model.Umkm) error {
	query := `
		INSERT INTO umkms (name, business_name, phone_number, email, address, business_type, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		u.Name,
		u.BusinessName,
		u.PhoneNumber,
		u.Email,
		u.Address,
		u.BusinessType,
		u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("umkm_repo.Create: %w", err)
	}
	return nil
}

func (r *umkmRepo) FindByPhone(ctx context.Context, phone string) (*model.Umkm, error) {
	query := `
		SELECT id, name, business_name, phone_number, email, address, business_type,
		       password_hash, is_active, created_at, updated_at, deleted_at
		FROM umkms
		WHERE phone_number = $1 AND deleted_at IS NULL`

	u := &model.Umkm{}
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&u.ID, &u.Name, &u.BusinessName, &u.PhoneNumber, &u.Email,
		&u.Address, &u.BusinessType, &u.PasswordHash, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("umkm_repo.FindByPhone: %w", err)
	}
	return u, nil
}

func (r *umkmRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Umkm, error) {
	query := `
		SELECT id, name, business_name, phone_number, email, address, business_type,
		       password_hash, is_active, created_at, updated_at, deleted_at
		FROM umkms
		WHERE id = $1 AND deleted_at IS NULL`

	u := &model.Umkm{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.BusinessName, &u.PhoneNumber, &u.Email,
		&u.Address, &u.BusinessType, &u.PasswordHash, &u.IsActive,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("umkm_repo.FindByID: %w", err)
	}
	return u, nil
}
