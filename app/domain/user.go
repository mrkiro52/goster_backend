package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid; primaryKey"`
	Email        string    `gorm:"uniqueIndex; not null"`
	Role         string    `gorm:"type:varchar(20);not null; default:user"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("Failed to generate UUID: %w", err)
		}
		u.ID = id
	}

	if u.Role != RoleAdmin && u.Role != RoleUser {
		return errors.New("invalid role: must be 'admin' or 'user'")
	}
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

func (u *User) SetPassword(plainPassword string) error {
	if plainPassword == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}
