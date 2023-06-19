package emailstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
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
	if s.ReadOneEmail(ctx, email) {
		return storageerrors.ErrEmailExists
	}

	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	byteSlice := []byte(fmt.Sprint(email, "\n"))
	_, err = file.Write(byteSlice)
	if err != nil {
		return err
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
		if strings.Contains(scanner.Text(), string(email)) {
			return true
		}
	}

	return false
}

func (s *Storage) ReadAllEmails(_ context.Context) ([]btcpricelb.Email, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return []btcpricelb.Email{}, err
	}

	var emailsList []btcpricelb.Email
	if err = json.Unmarshal(data, &emailsList); err != nil {
		return []btcpricelb.Email{}, err
	}

	return emailsList, err
}
