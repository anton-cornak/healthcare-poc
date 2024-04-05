package main

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name  string
		cfg   config
		valid bool
	}{
		{
			name: "ValidConfig",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:     "localhost",
					port:     "5432",
					user:     "user",
					password: "password",
					dbname:   "database",
					sslmode:  "disable",
				},
			},
			valid: true,
		},
		{
			name: "MissingPort",
			cfg: config{
				dbConn: dbConfig{
					host:     "localhost",
					port:     "5432",
					user:     "user",
					password: "password",
					dbname:   "database",
					sslmode:  "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnHost",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					port:     "5432",
					user:     "user",
					password: "password",
					dbname:   "database",
					sslmode:  "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnPort",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:     "localhost",
					user:     "user",
					password: "password",
					dbname:   "database",
					sslmode:  "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnUser",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:     "localhost",
					port:     "5432",
					password: "password",
					dbname:   "database",
					sslmode:  "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnPassword",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:    "localhost",
					port:    "5432",
					user:    "user",
					dbname:  "database",
					sslmode: "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnDbname",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:     "localhost",
					port:     "5432",
					user:     "user",
					password: "password",
					sslmode:  "disable",
				},
			},
			valid: false,
		},
		{
			name: "MissingDbConnSslmode",
			cfg: config{
				port: "8080",
				dbConn: dbConfig{
					host:     "localhost",
					port:     "5432",
					user:     "user",
					password: "password",
					dbname:   "database",
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.cfg)
			if tt.valid && err != nil {
				t.Errorf("Expected valid configuration, but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid configuration, but got no error")
			}
		})
	}
}
