package domain

import "context"

type UserRepository interface {
	Register(ctx context.Context, user *User) error
}
