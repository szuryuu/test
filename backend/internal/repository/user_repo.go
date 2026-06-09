package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepo — reserved for future multi-user scope.
// MVP uses UmkmRepo for all user/auth operations.
type UserRepo interface{}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepo {
	return &userRepo{pool: pool}
}
