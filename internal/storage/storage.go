package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/btc-price/internal/storageerrors"

	"github.com/btc-price/pkg/btcpricelb"
)

type Storage struct {
	path string
}

func NewStorage(path string) *Storage {
	return &Storage{
		path: path,
	}
}

func (s *Storage) AddEmail(ctx context.Context, email btcpricelb.Email) error {
	if !s.validateEmail(email) {
		return storageerrors.ErrInvalidEmail
	}

	if s.ReadOneEmail(ctx, email) {
		return storageerrors.ErrEmailExists
	}

	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("open file: %s", err)
	}
	defer file.Close()

	byteSlice := []byte(fmt.Sprint(email, "\n"))
	_, err = file.Write(byteSlice)
	if err != nil {
		return fmt.Errorf("write to file: %s", err)
	}

	return nil
}

func (s *Storage) ReadOneEmail(_ context.Context, email btcpricelb.Email) bool {
	file, err := os.Open(s.path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.EqualFold(scanner.Text(), string(email)) {
			return true
		}
	}

	return false
}

func (s *Storage) ReadAllEmails(_ context.Context) ([]btcpricelb.Email, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return []btcpricelb.Email{}, fmt.Errorf("read from file: %s", err)
	}

	var emailsList []btcpricelb.Email
	if err = json.Unmarshal(data, &emailsList); err != nil {
		return []btcpricelb.Email{}, fmt.Errorf("unmarshal file: %s", err)
	}

	return emailsList, nil
}

func (s *Storage) validateEmail(email btcpricelb.Email) bool {
	_, err := mail.ParseAddress(string(email))
	return err == nil
}
