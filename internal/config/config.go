package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Hour
	defaultRefreshTokenTTL        = 24 * time.Hour
	defaultVerificationCodeLength = 8
)

type (
	Config struct {
		Environment string
		HTTP        HTTPConfig
		Auth        AuthConfig
	}

	HTTPConfig struct {
		Host              string        `mapstructure:"host"`
		Port              string        `mapstructure:"port"`
		ReadTimeout       time.Duration `mapstructure:"read_timeout"`
		WriteTimeOut      time.Duration `mapstructure:"write_timeout"`
		MaxHeaderMegabyte int           `mapstructure:"max_header_megabyte"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verification_code_length"`
	}

	JWTConfig struct {
		SecretKey    string
		AccessToken  TokenConfig
		RefreshToken TokenConfig
	}

	TokenConfig struct {
		PrivateKey string        `mapstructure:"private_key"`
		PublicKey  string        `mapstructure:"public_key"`
		MaxAge     time.Duration `mapstructure:"max_age"`
	}
)

// TODO: uncomment for use when writing the server
//const configDir = "configs"
//func New() (*Config, error) {
//	return Init(configDir, os.Getenv("APP_ENV"))
//}

func Init(configsDir string, env string) (*Config, error) {
	if err := parseConfig(configsDir, env); err != nil {
		return nil, err
	}

	fmt.Printf("testConfigFetch: %#v\n", viper.AllSettings())

	populateDefault()

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	fetchFromEnv(configsDir, &cfg)

	return &cfg, nil
}

func populateDefault() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.read_timeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.write_timeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.max_header_megabyte", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("auth.verification_code_length", defaultVerificationCodeLength)
	viper.SetDefault("auth.access_token.max_age", defaultAccessTokenTTL)
	viper.SetDefault("auth.refresh_token.max_age", defaultRefreshTokenTTL)
}

func parseConfig(folder string, env string) error {
	viper.AddConfigPath(folder)

	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	baseConf := viper.AllSettings()

	viper.SetConfigName(env)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	envConf := viper.AllSettings()

	deepMerge(baseConf, envConf)

	//return viper.MergeInConfig()
	return viper.MergeConfigMap(baseConf)
}

func unmarshal(config *Config) error {
	// http
	if err := viper.UnmarshalKey("http.port", &config.HTTP.Port); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.max_header_megabyte", &config.HTTP.MaxHeaderMegabyte); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.read_timeout", &config.HTTP.ReadTimeout); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http.write_timeout", &config.HTTP.WriteTimeOut); err != nil {
		return err
	}

	// auth
	if err := viper.UnmarshalKey("auth.verification_code_length", &config.Auth.VerificationCodeLength); err != nil {
		return err
	}

	// jwt access_token
	if err := viper.UnmarshalKey("jwt.access_token.private_key", &config.Auth.JWT.AccessToken.PrivateKey); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt.access_token.public_key", &config.Auth.JWT.AccessToken.PublicKey); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt.access_token.max_age", &config.Auth.JWT.AccessToken.MaxAge); err != nil {
		return err
	}

	// jwt refresh_token
	if err := viper.UnmarshalKey("jwt.refresh_token.private_key", &config.Auth.JWT.RefreshToken.PrivateKey); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt.refresh_token.public_key", &config.Auth.JWT.RefreshToken.PublicKey); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt.refresh_token.max_age", &config.Auth.JWT.RefreshToken.MaxAge); err != nil {
		return err
	}

	return nil
}

func fetchFromEnv(configsDir string, config *Config) {
	err := godotenv.Load(configsDir + "/.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	config.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	config.Auth.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	config.HTTP.Host = os.Getenv("HTTP_HOST")
	config.Environment = os.Getenv("APP_ENV")
}

func deepMerge(dst map[string]interface{}, src ...map[string]interface{}) {
	for _, srcMap := range src {
		for key, srcValue := range srcMap {
			dstValue, exists := dst[key]

			if !exists {
				dst[key] = srcValue
				continue
			}

			srcMap, srcMapOk := srcValue.(map[string]interface{})
			dstMap, dstMapOk := dstValue.(map[string]interface{})

			if srcMapOk && dstMapOk {
				deepMerge(dstMap, srcMap)
			} else {
				dst[key] = srcValue
			}
		}
	}
}
