package config

import (
	"errors"
	"net"

	"github.com/azyoskol/word_of_wisdom/internal/log"
	"github.com/spf13/viper"
)

const (
	ConfigPrefix = "SRV"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

func NewConfig() (*viper.Viper, error) {
	cfg := viper.New()

	cfg.SetEnvPrefix(ConfigPrefix)
	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	err := cfg.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			log.Warnw("Error when Fetching Configuration",
				log.M{
					"err": err,
				},
			)
			return nil, err
		}
	}

	return cfg, nil
}

type serverConfig struct {
	vpr *viper.Viper
}

func (c *serverConfig) GetAddress() string {
	endpoint := net.JoinHostPort(
		c.vpr.GetString("HOST"),
		c.vpr.GetString("PORT"),
	)
	return endpoint
}

func (c *serverConfig) GetNumberOfTimesRAppliesF() int64 {
	return c.vpr.GetInt64("NUMBER_OF_TIMES_R_APPLIES_F")
}
func (c *serverConfig) GetSizeOfEachValue() int64 {
	return c.vpr.GetInt64("SIZE_OF_EACH_VALUE")
}

func (c *serverConfig) prepare() (err error) {
	c.vpr, err = NewConfig()
	if err != nil {
		return err
	}

	return nil
}

func NewServerConfig() (*serverConfig, error) {
	cfg := &serverConfig{}

	err := cfg.prepare()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
