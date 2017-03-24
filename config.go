package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// server modes
const (
	DEVELOP_MODE    = "develop"
	PRODUCTION_MODE = "production"
	TEST_MODE       = "test"
)

// config holds all configuration for the server. It pulls from three places (in order):
// 		1. environment variables
// 		2. config.[server_mode].json <- eg: config.test.json
// 		3. config.json
//
// env variables win, but can only set config who's json is ALL_CAPS
// it's totally fine to not have, say, config.develop.json defined, and just
// rely on a base config.json. But if you're in production mode & config.production.json
// exists, that will be read *instead* of config.json.
//
// configuration is read at startup and cannot be alterd without restarting the server.
type config struct {
	// configurable environment variable
	Environment string `json:"GOLANG_ENV"`

	// port to listen on, will be read from PORT env variable if present.
	Port string `json:"PORT"`

	// root url for service
	UrlRoot string `json:"URL_ROOT"`

	// url of postgres app db. required.
	PostgresDbUrl string `json:"POSTGRES_DB_URL"`

	// secret for encrypted session cookies. required.
	SessionSecret string `json:"SESSION_SECRET"`

	// Public Key to use for signing. Not yet used, but required soon.
	PublicKey string `json:"PUBLIC_KEY"`

	// TLS (HTTPS) enable support via LetsEncrypt, default false
	// should be true in production
	TLS bool `json:"TLS"`

	// Key to store user's cookie under
	UserCookieKey string `json:"USER_COOKIE_KEY"`

	// support CORS signing from a list of origins
	AllowedOrigins []string `json:"ALLOWED_ORIGINS"`

	// setting HTTP_AUTH_USERNAME & HTTP_AUTH_PASSWORD
	// will enable basic http auth for the server. This is a single
	// username & password that must be passed in with every request.
	// leaving these values blank will disable http auth
	// read from env variable: HTTP_AUTH_USERNAME
	HttpAuthUsername string `json:"HTTP_AUTH_USERNAME"`
	// read from env variable: HTTP_AUTH_PASSWORD
	HttpAuthPassword string `json:"HTTP_AUTH_PASSWORD"`

	// if true, requests that have X-Forwarded-Proto: http will be redirected
	// to their https variant
	ProxyForceHttps bool
	// CertbotResponse is only for doing manual SSL certificate generation
	// via LetsEncrypt.
	CertbotResponse string `json:"CERTBOT_RESPONSE"`
}

// initConfig pulls configuration from config.json
func initConfig(mode string) (cfg *config, err error) {
	cfg = &config{}

	if err := loadConfigFile(mode, cfg); err != nil {
		return cfg, err
	}

	// override config settings with env settings, passing in the current configuration
	// as the default. This has the effect of leaving the config.json value unchanged
	// if the env variable is empty
	cfg.Port = readEnvString("PORT", cfg.Port)
	cfg.UrlRoot = readEnvString("URL_ROOT", cfg.UrlRoot)
	cfg.PublicKey = readEnvString("PUBLIC_KEY", cfg.PublicKey)
	cfg.TLS = readEnvBool("TLS", cfg.TLS)
	cfg.UrlRoot = readEnvString("SESSION_SECRET", cfg.SessionSecret)
	cfg.UserCookieKey = readEnvString("USER_COOKIE_KEY", cfg.UserCookieKey)
	cfg.PostgresDbUrl = readEnvString("POSTGRES_DB_URL", cfg.PostgresDbUrl)
	cfg.AllowedOrigins = readEnvStringSlice("ALLOWED_ORIGINS", cfg.AllowedOrigins)
	cfg.HttpAuthUsername = readEnvString("HTTP_AUTH_USERNAME", cfg.HttpAuthUsername)
	cfg.HttpAuthPassword = readEnvString("HTTP_AUTH_PASSWORD", cfg.HttpAuthPassword)
	cfg.CertbotResponse = readEnvString("CERTBOT_RESPONSE", cfg.CertbotResponse)

	// make sure port is set
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.UserCookieKey == "" {
		// TODO - this should be set based on on url root value
		cfg.UserCookieKey = "identity.user"
	}

	err = requireConfigStrings(map[string]string{
		"PORT":            cfg.Port,
		"POSTGRES_DB_URL": cfg.PostgresDbUrl,
		"PUBLIC_KEY":      cfg.PublicKey,
	})

	return
}

// readEnvString reads key from the environment, returns def if empty
func readEnvString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

// readEnvBool read key form the env, converting to a boolean value. returns def if empty
func readEnvBool(key string, def bool) bool {
	if env := os.Getenv(key); env != "" {
		return env == "true" || env == "TRUE" || env == "t"
	}
	return def
}

// readEnvString reads a slice of strings from key environment var, returns def if empty
func readEnvStringSlice(key string, def []string) []string {
	if env := os.Getenv(key); env != "" {
		return strings.Split(env, ",")
	}
	return def
}

// requireConfigStrings panics if any of the passed in values aren't set
func requireConfigStrings(values map[string]string) error {
	for key, value := range values {
		if value == "" {
			return fmt.Errorf("%s env variable or config key must be set", key)
		}
	}

	return nil
}

// checks for config.[mode].json file to read configuration from if the file exists
// defaults to config.json, silently fails if no configuration file is present.
func loadConfigFile(mode string, cfg *config) (err error) {
	var data []byte

	fileName := fmt.Sprintf("config.%s.json", mode)
	if !fileExists(fileName) {
		fileName = "config.json"
		if !fileExists(fileName) {
			return nil
		}
	}

	logger.Printf("reading config file: %s", fileName)
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("error reading %s: %s", fileName, err)
		return
	}

	// unmarshal ("decode") config data into a config struct
	if err = json.Unmarshal(data, cfg); err != nil {
		err = fmt.Errorf("error parsing %s: %s", fileName, err)
		return
	}

	return
}

// Does this file exist?
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// outputs any notable settings to stdout
func printConfigInfo() {
	// TODO
}
