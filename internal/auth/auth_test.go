package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if hash == password {
		t.Fatal("hash should not be equal to the password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "correct password",
			input:    password,
			expected: true,
		},
		{
			name:     "incorrect password",
			input:    wrongPassword,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.input, hash)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if match != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, match)
			}
		})
	}
}

func TestHashPassword_UniqueHashes(t *testing.T) {
	password := "samepassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("expected different hashes for same password due to salting")
	}
}

func TestJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		returnedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate token: %v", err)
		}

		if returnedID != userID {
			t.Errorf("expected %v, got %v", userID, returnedID)
		}
	})

	t.Run("invalid secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, "wrongsecret")
		if err == nil {
			t.Error("expected error for invalid secret, got nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := ValidateJWT("not.a.valid.token", secret)
		if err == nil {
			t.Error("expected error for malformed token, got nil")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		expected    string
		expectError bool
	}{
		{
			name:        "valid bearer token",
			header:      "Bearer abc123",
			expected:    "abc123",
			expectError: false,
		},
		{
			name:        "missing header",
			header:      "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid format - no bearer",
			header:      "abc123",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid format - only bearer",
			header:      "Bearer",
			expected:    "",
			expectError: true,
		},
		{
			name:        "extra spaces",
			header:      "Bearer    xyz789",
			expected:    "xyz789",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, token)
			}
		})
	}
}
