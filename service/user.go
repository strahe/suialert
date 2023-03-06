package service

import (
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10"

	"github.com/strahe/suialert/model"
)

type UserService struct {
	db model.Storage
}

func NewUserService(db model.Storage) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindByDiscordID(id string) (*model.User, error) {
	var user model.User
	fmt.Println(s.db)
	err := s.db.AsORM().Model(&user).Where("discord_id =?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		fmt.Println(reflect.TypeOf(err))
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Create(user *model.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	_, err := s.db.AsORM().Model(user).Insert()
	return err
}
