package service_test

import (
	"context"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo/postgres"
	"social-media-backend/internal/service"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo    postgres.UserRepository
		sessionRepo postgres.SessionsRepository
		// Named input parameters for target function.
		params  service.RegisterParams
		want    domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewAuthService(tt.userRepo, tt.sessionRepo)
			got, gotErr := s.Register(context.Background(), tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Register() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Register() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo    postgres.UserRepository
		sessionRepo postgres.SessionsRepository
		// Named input parameters for target function.
		params  service.LoginParams
		want    domain.AuthTokens
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewAuthService(tt.userRepo, tt.sessionRepo)
			got, gotErr := s.Login(context.Background(), tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userRepo    postgres.UserRepository
		sessionRepo postgres.SessionsRepository
		// Named input parameters for target function.
		params  service.RefreshParams
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewAuthService(tt.userRepo, tt.sessionRepo)
			got, gotErr := s.RefreshAccessToken(context.Background(), tt.params)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RefreshAccessToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RefreshAccessToken() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("RefreshAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
