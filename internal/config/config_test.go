package config

import "testing"

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "defaults when environment is empty",
			env:  map[string]string{},
			want: Config{
				Port:              "8080",
				AppEnv:            "development",
				GraphQLPlayground: true,
				LogLevel:          "info",
				RepoDriver:        "memory",
				DatabaseURL:       "",
			},
		},
		{
			name: "overrides from environment",
			env: map[string]string{
				"PORT":               "9090",
				"APP_ENV":            "production",
				"GRAPHQL_PLAYGROUND": "false",
				"LOG_LEVEL":          "debug",
				"REPO_DRIVER":        "postgres",
				"DATABASE_URL":       "postgres://user:pass@localhost:5432/products",
			},
			want: Config{
				Port:              "9090",
				AppEnv:            "production",
				GraphQLPlayground: false,
				LogLevel:          "debug",
				RepoDriver:        "postgres",
				DatabaseURL:       "postgres://user:pass@localhost:5432/products",
			},
		},
		{
			name: "invalid boolean falls back to default",
			env:  map[string]string{"GRAPHQL_PLAYGROUND": "notabool"},
			want: Config{
				Port:              "8080",
				AppEnv:            "development",
				GraphQLPlayground: true,
				LogLevel:          "info",
				RepoDriver:        "memory",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string { return tt.env[key] }
			got := Load(getenv)
			if got != tt.want {
				t.Fatalf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
