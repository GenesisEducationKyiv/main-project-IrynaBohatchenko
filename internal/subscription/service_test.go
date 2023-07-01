//go:build integration
// +build integration

package subscription

import (
	"bufio"
	"context"
	"github.com/btc-price/internal/storage"
	"github.com/btc-price/internal/storageerrors"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type SubscribeSuite struct {
	suite.Suite
	FilePath string
	Service  *Service
	Ctx      context.Context
}

func (s *SubscribeSuite) SetupSuite() {
	s.FilePath = "./emails_test.txt"
	s.Service = NewService(storage.NewStorage(s.FilePath))
	s.Ctx = context.Background()
}

func (s *SubscribeSuite) SetupTest() {
	_, err := os.Create(s.FilePath)
	if err != nil {
		s.FailNowf("", "create file %s", err)
	}
}

func (s *SubscribeSuite) TearDownTest() {
	os.Remove(s.FilePath)
}

func (s *SubscribeSuite) TestSuccessfulSubscription() {
	testCases := []struct {
		name  string
		email string
	}{
		{
			name:  "successfully subscribed",
			email: "test_email@gmail.com",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			if err := s.Service.HandleSubscribe(s.Ctx, tt.email); err != nil {
				s.FailNowf("", "error subscribe %s", err)
			}

			file, err := os.Open(s.FilePath)
			if err != nil {
				s.FailNowf("", "error open file %s", err)
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				s.True(strings.EqualFold(scanner.Text(), tt.email))
			}
		})
	}
}

func (s *SubscribeSuite) TestRepeatedEmail() {
	type args struct {
		email         string
		emailRepeated string
	}

	testCases := []struct {
		name string
		args args
		want error
	}{
		{
			name: "write repeated email",
			args: args{
				email:         "test_email@gmail.com",
				emailRepeated: "test_email@gmail.com",
			},
			want: storageerrors.ErrEmailExists,
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			if err := s.Service.HandleSubscribe(s.Ctx, tt.args.email); err != nil {
				s.FailNowf("", "error subscribe %s", err)
			}

			s.ErrorIs(s.Service.HandleSubscribe(s.Ctx, tt.args.emailRepeated), tt.want)
		})
	}
}

func TestService_HandleSubscribe(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}
