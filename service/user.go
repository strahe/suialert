package service

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"

	"github.com/strahe/suialert/model"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindByDiscordID(id string) (*model.User, error) {
	var user model.User
	err := s.db.Where(&model.User{DiscordID: &id}, "DiscordID").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (s *UserService) Create(user *model.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	return s.db.Create(user).Error
}

func (s *UserService) FindOrCreateByDiscordUser(du *discordgo.User) (*model.User, error) {
	if du == nil {
		return nil, fmt.Errorf("discord user is nil")
	}
	u := model.User{
		DiscordID:   &du.ID,
		Name:        fmt.Sprintf("%s#%s", du.Username, du.Discriminator),
		DiscordInfo: du,
	}
	return &u, s.db.Where(&model.User{DiscordID: &du.ID}, "DiscordID").FirstOrCreate(&u).Error
}
