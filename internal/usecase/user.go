package usecase

import (
	"context"
	"errors"

	"example.com/user-management/internal/domain/entity"
	"example.com/user-management/internal/domain/repository"
)

type UserUseCase interface {
	Create(ctx context.Context, name, email, password string) (*entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, id, name, email string) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entity.User, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (u *userUseCase) Create(ctx context.Context, name, email, password string) (*entity.User, error) {
	// Check if user with same email exists
	existingUser, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user, err := entity.NewUser(name, email, password)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Get(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepository.Get(ctx, id)
}

func (u *userUseCase) Update(ctx context.Context, id, name, email string) (*entity.User, error) {
	// Get existing user
	user, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update user
	if err := user.Update(name, email); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Delete(ctx context.Context, id string) error {
	return u.userRepository.Delete(ctx, id)
}

func (u *userUseCase) List(ctx context.Context) ([]*entity.User, error) {
	return u.userRepository.List(ctx)
}
