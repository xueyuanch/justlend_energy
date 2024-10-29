package config

import (
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnv looks up the given key from the environment, returning its value if
// it exists, and otherwise returning the given fallback value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvBool looks up the given key from the environment, returning true if
// it exists and value is equals to `TRUE`, and otherwise returning the
// given fallback value.
func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok &&
		strings.ToUpper(value) == "TRUE" {
		return true
	}
	return fallback
}

// GetEnvInt looks up the given key from the environment and expects an integer,
// returning the integer value if it exists, and otherwise returning the given
// fallback value.
//
// If the environment variable has a value, but it can't be parsed as an integer,
// GetEnvInt terminates the program.
func GetEnvInt(key string, fallback int) int {
	if s, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("getEnvInt: bad value %q for %s: %v", s, key, err)
		}
		return v
	}
	return fallback
}

func GetEnvInt64(key string, fallback int64) int64 {
	if s, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalf("getEnvInt64: bad value %q for %s: %v", s, key, err)
		}
		return v
	}
	return fallback
}

func GetEnvInt32(key string, fallback int32) int32 {
	if s, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Fatalf("getEnvInt32: bad value %q for %s: %v", s, key, err)
		}
		return int32(v)
	}
	return fallback
}

// GetEnvDuration wraps GetEnvInt and returns the result as type of duration.
func GetEnvDuration(key string, fallback int) time.Duration {
	return time.Duration(GetEnvInt(key, fallback))
}

// GetEnvFloat64 looks up the given key from the environment and expects a
// float64, returning the float64 value if it exists, and otherwise returning
// the given fallback value.
//
// If the environment variable has a value, but it can't be parsed as an float,
// GetEnvFloat64 terminates the program.
func GetEnvFloat64(key string, fallback float64) float64 {
	if s, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatalf("getEnvFloat64: bad value %q for %s: %v", s, key, err)
		}
		return v
	}
	return fallback
}

// GetEnvHexBytes looks up the given key from the environment and expects the
// bytes of a hexadecimal string, returning the hexadecimal string's bytes if
// it exists, and otherwise returning the given fallback value's bytes.
//
// If the environment variable or fallback has a value but it can't be parsed by
// hex.DecodeString() GetEnvHexBytes terminates the program.
func GetEnvHexBytes(key string, fallback string) []byte {
	if s, ok := os.LookupEnv(key); ok {
		fallback = s
	}
	v, err := hex.DecodeString(fallback)
	if err != nil {
		log.Fatalf("getEnvHexBytes: bad value %q for %s: %v", fallback, key, err)
	}
	return v
}

// Config holds shared configuration values used in instantiating
// our server components.
type Config struct {

	// - Addr used for trxEnergy hosting.
	Addr string

	// The duration for which the server gracefully wait for existing
	// connections to finish - e.g. 15s or 1m
	GracefulTimeout time.Duration

	// Domain maps to system.
	Domain string

	// Keys used for secure cookie encryption.
	SCHashKey, SCBlockKey []byte
}

const (
	// Defines default value for SCHashKey & SCBlockKey.
	defaultSCHashKey  = "00EC379CC076D7779011961363D1F831"
	defaultSCBlockKey = "8CDB4C835C4741B01710E91617EC7EA5"

	// StatementTimeout is the value of the Postgres statement_timeout parameter.
	// Statements that run longer than this are terminated.
	StatementTimeout = 5 * time.Minute
)

// Resolve resolves all configuration values provided by the config package. It
// must be called before any configuration values are used.
func Resolve() *Config {
	return &Config{
		// Resolve host information.
		Addr:            GetEnv("MS_ADDR", ":8085"),
		Domain:          GetEnv("DOMAIN", ""),
		GracefulTimeout: time.Duration(GetEnvInt("GRACEFUL_TIMEOUT", 15)) * time.Second,
		// Resolve http cookie hash & block keys.
		SCHashKey:  GetEnvHexBytes("SESSION_HASH_KEY", defaultSCHashKey),
		SCBlockKey: GetEnvHexBytes("SESSION_BLOCK_KEY", defaultSCBlockKey),
	}
}

func IsReleaseMode() bool { return os.Getenv("RUNNING_MODE") == "Release" }
