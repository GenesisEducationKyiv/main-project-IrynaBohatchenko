// unit test of validator
package emailstorage

import (
	"github.com/btc-price/pkg/btcpricelb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_validateEmail(t *testing.T) {
	type fields struct {
		path string
	}
	type args struct {
		email btcpricelb.Email
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "valid email",
			fields: fields{path: "./emails_test.txt"},
			args:   args{email: "test_email@gmail.com"},
			want:   true,
		},
		{
			name:   "invalid email",
			fields: fields{path: "./emails_test.txt"},
			args:   args{email: "test_email@"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				path: tt.fields.path,
			}
			assert.Equal(t, tt.want, s.validateEmail(tt.args.email))
		})
	}
}
