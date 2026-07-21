// Package config centraliza la configuracion de la aplicacion leida desde el entorno.
package config

import "strconv"

// Config agrupa los parametros de ejecucion de la aplicacion.
type Config struct {
	Port              string
	AppEnv            string
	GraphQLPlayground bool
	LogLevel          string
	RepoDriver        string
	DatabaseURL       string
}

// Load construye la configuracion a partir de la funcion de busqueda indicada,
// aplicando valores por defecto. Recibe getenv para ser pura y testeable.
func Load(getenv func(string) string) Config {
	return Config{
		Port:              orDefault(getenv("PORT"), "8080"),
		AppEnv:            orDefault(getenv("APP_ENV"), "production"),
		GraphQLPlayground: parseBool(getenv("GRAPHQL_PLAYGROUND"), true),
		LogLevel:          orDefault(getenv("LOG_LEVEL"), "info"),
		RepoDriver:        orDefault(getenv("REPO_DRIVER"), "memory"),
		DatabaseURL:       getenv("DATABASE_URL"),
	}
}

// IsDevelopment indica si la app corre en modo desarrollo. Es un match exacto contra
// "development" (fail-safe): cualquier otro valor, o el default "production", deshabilita
// el playground y la introspection.
func (c Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// orDefault devuelve value si no esta vacio, o fallback en caso contrario.
func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// parseBool interpreta value como booleano, devolviendo fallback si esta vacio o es invalido.
func parseBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
