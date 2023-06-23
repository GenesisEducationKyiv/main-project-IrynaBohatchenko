// integration test of file storage
package btcpriceservice

import (
	"bufio"
	"context"
	"github.com/btc-price/internal/emailstorage"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestService_HandleSubscribe(t *testing.T) {
	filePath := "./emails_test.txt"
	email := "test_email@gmail.com"
	ctx := context.Background()

	t.Run("successfully writen", func(t *testing.T) {
		file, err := os.Create(filePath)
		defer os.Remove(filePath)
		if err != nil {
			t.Fatalf("create file %s", err)
		}

		srv := NewService(nil, emailstorage.NewStorage(filePath), nil)

		if err := srv.HandleSubscribe(ctx, email); err != nil {
			t.Fatalf("error subscribe %s", err)
		}

		file, err = os.Open(filePath)
		if err != nil {
			t.Fatalf("error open file %s", err)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			assert.True(t, strings.Contains(scanner.Text(), email))
		}
	})

	t.Run("write repeated email", func(t *testing.T) {
		_, err := os.Create(filePath)
		defer os.Remove(filePath)
		if err != nil {
			t.Fatalf("create file %s", err)
		}

		srv := NewService(nil, emailstorage.NewStorage(filePath), nil)

		if err := srv.HandleSubscribe(ctx, email); err != nil {
			t.Fatalf("error subscribe %s", err)
		}

		assert.Error(t, srv.HandleSubscribe(ctx, email))
	})
}
