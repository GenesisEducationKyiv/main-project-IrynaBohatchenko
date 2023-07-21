//go:build !integration

package usecase

import (
	"github.com/btc-price/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_validateEmail(t *testing.T) {
	t.Parallel()

	type args struct {
		email models.Email
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid email",
			args: args{email: "test_email@gmail.com"},
			want: true,
		},
		{
			name: "invalid email",
			args: args{email: "test_email@"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubscriptionService{}
			assert.Equal(t, tt.want, s.validateEmail(tt.args.email))
		})
	}
}
