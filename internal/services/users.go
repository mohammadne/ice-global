package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/storage"
)

type Users interface {
	RetrieveUserOptional(ctx context.Context, cookie string) (*entities.User, error)
	RetrieveUserRequired(ctx context.Context, cookie string) (*entities.User, error)
}

func NewUsers(usersStorage storage.Users) Users {
	return &users{usersStorage: usersStorage}
}

type users struct {
	usersStorage storage.Users
}

func (u *users) RetrieveUserOptional(ctx context.Context, cookie string) (result *entities.User, err error) {
	{ // validation
		if cookie == "" {
			return nil, errors.New("the cookie should be provided")
		}
	}

	storageUser, err := u.usersStorage.RetrieveUserByCookie(ctx, cookie)
	if err != nil {
		if err == storage.ErrorUserNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	return &entities.User{
		Id:     storageUser.Id,
		Cookie: storageUser.Cookie,
	}, nil
}

func (u *users) RetrieveUserRequired(ctx context.Context, cookie string) (result *entities.User, err error) {
	{ // validation
		if cookie == "" {
			return nil, errors.New("the cookie should be provided")
		}
	}

	storageUser, err := u.usersStorage.RetrieveUserByCookie(ctx, cookie)
	if err != nil {
		if err == storage.ErrorUserNotFound {
			storageUser := &storage.User{
				Cookie:    cookie,
				CreatedAt: time.Now(),
			}

			id, err := u.usersStorage.CreateUser(ctx, storageUser)
			if err != nil {
				return nil, fmt.Errorf("error creating user: %v", err)
			}
			return &entities.User{
				Id:     id,
				Cookie: cookie,
			}, nil
		}
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	return &entities.User{
		Id:     storageUser.Id,
		Cookie: storageUser.Cookie,
	}, nil
}
