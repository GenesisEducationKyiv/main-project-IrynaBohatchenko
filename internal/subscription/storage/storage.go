package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/btc-price/internal/models"
	"os"
	"strings"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{
		path: path,
	}
}

func (s *Storage) AddUser(_ context.Context, user *models.User) error {
	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("open file: %s", err)
	}
	defer file.Close()

	byteSlice := []byte(fmt.Sprint(user, "\n"))
	_, err = file.Write(byteSlice)
	if err != nil {
		return fmt.Errorf("write to file: %s", err)
	}

	return nil
}

func (s *Storage) GetUser(_ context.Context, user *models.User) bool {
	file, err := os.Open(s.path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.EqualFold(scanner.Text(), string(user.Email)) {
			return true
		}
	}

	return false
}

func (s *Storage) GetUsersList(_ context.Context) ([]*models.User, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return []*models.User{}, fmt.Errorf("read from file: %s", err)
	}

	var usersList []*models.User
	if err = json.Unmarshal(data, &usersList); err != nil {
		return []*models.User{}, fmt.Errorf("unmarshal file: %s", err)
	}

	return usersList, nil
}
