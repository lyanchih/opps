package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var (
	cfg              *Config
	ErrConfigNotInit = errors.New("Config still not initialize yet")
)

func ParseConf(file string) (*Config, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	tmpCfg := &Config{}
	if err = json.Unmarshal(bs, tmpCfg); err != nil {
		return nil, err
	}

	cfg = tmpCfg
	return cfg, nil
}

func GetConf() (*Config, error) {
	if cfg == nil {
		return nil, ErrConfigNotInit
	}

	return cfg, nil
}
