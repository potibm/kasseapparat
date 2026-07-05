package config

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisURL_URLObject(t *testing.T) {
	tests := []struct {
		name     string
		url      *RedisURL
		wantNil  bool
		wantHost string
		wantPath string
	}{
		{
			name:    "nil receiver",
			url:     nil,
			wantNil: true,
		},
		{
			name:     "valid redis URL",
			url:      ptr(RedisURL("redis://localhost:6379/0")),
			wantNil:  false,
			wantHost: "localhost:6379",
			wantPath: "/0",
		},
		{
			name:     "valid rediss URL with password",
			url:      ptr(RedisURL("rediss://:secret@redis.example.com:6380/5")),
			wantNil:  false,
			wantHost: "redis.example.com:6380",
			wantPath: "/5",
		},
		{
			name:    "invalid URL",
			url:     ptr(RedisURL("://invalid")),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.url.URLObject()
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantHost, got.Host)
				assert.Equal(t, tt.wantPath, got.Path)
			}
		})
	}
}

func TestRedisURL_RedisOptions(t *testing.T) {
	tests := []struct {
		name         string
		url          RedisURL
		wantNil      bool
		wantAddr     string
		wantPassword string
		wantDB       int
		wantTLS      bool
	}{
		{
			name:    "empty string",
			url:     RedisURL(""),
			wantNil: true,
		},
		{
			name:     "basic redis URL",
			url:      RedisURL("redis://localhost:6379"),
			wantAddr: "localhost:6379",
			wantDB:   0,
			wantTLS:  false,
		},
		{
			name:     "redis URL with db",
			url:      RedisURL("redis://localhost:6379/3"),
			wantAddr: "localhost:6379",
			wantDB:   3,
		},
		{
			name:         "redis URL with password",
			url:          RedisURL("redis://:mypass@redis.example.com:6380/7"),
			wantAddr:     "redis.example.com:6380",
			wantPassword: "mypass",
			wantDB:       7,
		},
		{
			name:         "redis URL with username and password",
			url:          RedisURL("redis://user:pass@localhost:6379/2"),
			wantAddr:     "localhost:6379",
			wantPassword: "pass",
			wantDB:       2,
		},
		{
			name:    "invalid db now correctly fails",
			url:     RedisURL("redis://localhost:6379/notadb"),
			wantNil: true,
		},
		{
			name:     "rediss URL sets TLS config",
			url:      RedisURL("rediss://localhost:6379"),
			wantAddr: "localhost:6379",
			wantDB:   0,
			wantTLS:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.url.RedisOptions()
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantAddr, got.Addr)
				assert.Equal(t, tt.wantPassword, got.Password)
				assert.Equal(t, tt.wantDB, got.DB)

				if tt.wantTLS {
					assert.NotNil(t, got.TLSConfig, "expected TLSConfig to be set for rediss://")
				} else {
					assert.Nil(t, got.TLSConfig, "expected TLSConfig to be nil for redis://")
				}
			}
		})
	}
}

func TestRedisURL_RedisOptions_ReturnsRedisOptionsType(t *testing.T) {
	url := RedisURL("redis://localhost:6379/1")
	opts := url.RedisOptions()

	assert.IsType(t, &redis.Options{}, opts)
}

func TestRedisURL_Redacted(t *testing.T) {
	tests := []struct {
		name string
		url  RedisURL
		want RedisURL
	}{
		{
			name: "redacts password",
			url:  RedisURL("redis://:secret@localhost:6379/0"),
			want: RedisURL("redis://:%2A%2A%2AREDACTED%2A%2A%2A@localhost:6379/0"),
		},
		{
			name: "no password",
			url:  RedisURL("redis://localhost:6379/0"),
			want: RedisURL("redis://localhost:6379/0"),
		},
		{
			name: "empty string",
			url:  RedisURL(""),
			want: RedisURL(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.url.Redacted())
		})
	}
}

func TestRedisURL_Validate(t *testing.T) {
	tests := []struct {
		name    string
		url     RedisURL
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid redis URL",
			url:     RedisURL("redis://localhost:6379/0"),
			wantErr: false,
		},
		{
			name:    "valid rediss URL",
			url:     RedisURL("rediss://redis.example.com:6380/1"),
			wantErr: false,
		},
		{
			name:    "invalid URL",
			url:     RedisURL("://invalid"),
			wantErr: true,
			errMsg:  "missing protocol scheme",
		},
		{
			name:    "http scheme",
			url:     RedisURL("http://localhost:6379/0"),
			wantErr: true,
			errMsg:  "invalid URL scheme: http",
		},
		{
			name:    "missing host IS NOW VALID",
			url:     RedisURL("redis:///0"),
			wantErr: false,
		},
		{
			name:    "host without port IS NOW VALID",
			url:     RedisURL("redis://localhost/0"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.url.Validate()
			if tt.wantErr {
				assert.Error(t, err)

				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
