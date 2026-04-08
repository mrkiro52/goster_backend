// service/user/service.go
package user

import (
	"context"
	"errors"

	"goster/domain"
	"goster/infrastructure/gorm"
)

type Service struct {
	userRepo  *gorm.UserRepository
	txFactory gorm.TransactionContextFactory
}

func NewService(userRepo *gorm.UserRepository, txFactory gorm.TransactionContextFactory) *Service {
	return &Service{
		userRepo:  userRepo,
		txFactory: txFactory,
	}
}

type RegisterRequest struct {
	Email    string
	Password string
	Role     string
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*domain.User, error) {
	// Создаем контекст с транзакцией
	txCtx := s.txFactory(ctx)
	defer txCtx.Rollback() // откатим если что-то пойдет не так

	// Проверяем существование
	exists, err := s.userRepo.ExistsByEmail(txCtx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Определяем роль
	role := req.Role
	if role != domain.RoleAdmin && role != domain.RoleUser {
		role = domain.RoleUser
	}

	// Создаем пользователя
	user := &domain.User{
		Email: req.Email,
		Role:  role,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(txCtx, user); err != nil {
		return nil, err
	}

	// Коммитим транзакцию
	if err := txCtx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.User, error) {
	// Логин без транзакции (только чтение)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
