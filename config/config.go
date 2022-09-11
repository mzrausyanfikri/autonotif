package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Base       Base             `yaml:"base" mapstructure:"base"`
	ChainList  map[string]Chain `yaml:"chain" mapstructure:"chain"`
	Target     Target           `yaml:"target" mapstructure:"target"`
	Repository Repository       `yaml:"repository" mapstructure:"repository"`
}

func (cfg *Config) Validate() error {
	if cfg.Repository.Postgresql.Enable && cfg.Repository.Localfile.Enable {
		return errors.New("only one repository option can be enabled")
	}
	return nil
}

type Base struct {
	SchedulerPeriod time.Duration `yaml:"schedulerPeriod" mapstructure:"schedulerPeriod"`
}

type Chain struct {
	Enable                bool                  `yaml:"enable" mapstructure:"enable"`
	ChainAPI              ChainAPI              `yaml:"api" mapstructure:"api"`
	ChainMessageStructure ChainMessageStructure `yaml:"messageStructure" mapstructure:"messageStructure"`
}

type ChainAPI struct {
	Endpoint string        `yaml:"endpoint" mapstructure:"endpoint"`
	Nodepool []string      `yaml:"nodepool" mapstructure:"nodepool"`
	Retry    int           `yaml:"retry" mapstructure:"retry"`
	Timeout  time.Duration `yaml:"timeout" mapstructure:"timeout"`
}

type ChainMessageStructure struct {
	Name       ChainMessageStructureAttr `yaml:"name" mapstructure:"name"`
	ProposalID ChainMessageStructureAttr `yaml:"proposalId" mapstructure:"proposalId"`
	Title      ChainMessageStructureAttr `yaml:"title" mapstructure:"title"`
	Status     ChainMessageStructureAttr `yaml:"status" mapstructure:"status"`
	Type       ChainMessageStructureAttr `yaml:"type" mapstructure:"type"`
	StartTime  ChainMessageStructureAttr `yaml:"startTime" mapstructure:"startTime"`
	EndTime    ChainMessageStructureAttr `yaml:"endTime" mapstructure:"endTime"`
	ViewLink   ChainMessageStructureAttr `yaml:"viewLink" mapstructure:"viewLink"`
}

type ChainMessageStructureAttr struct {
	Const    string `yaml:"const" mapstructure:"const"`
	JSONPath string `yaml:"jsonPath" mapstructure:"jsonPath"`
}

type Target struct {
	Telegram Telegram `yaml:"telegram" mapstructure:"telegram"`
}

type Telegram struct {
	Enable    bool   `yaml:"enable" mapstructure:"enable"`
	Token     string `yaml:"token" mapstructure:"token"`
	ChannelID int64  `yaml:"channelId" mapstructure:"channelId"`
}

type Repository struct {
	Postgresql Postgresql `yaml:"postgresql" mapstructure:"postgresql"`
	Localfile  Localfile  `yaml:"localfile" mapstructure:"localfile"`
}

type Postgresql struct {
	Enable   bool   `yaml:"enable" mapstructure:"enable"`
	Address  string `yaml:"address" mapstructure:"address"`
	Database string `yaml:"database" mapstructure:"database"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
}

type Localfile struct {
	Enable bool   `yaml:"enable" mapstructure:"enable"`
	Dir    string `yaml:"dir" mapstructure:"dir"`
}

// fullpath example: ./config/config.yaml
func NewConfig(fullpath string) (*Config, error) {
	var cfg = Config{}

	configPath, configName, configType := parseConfigPath(fullpath)
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read viper config: %s", err)
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %s", err)
	}

	return &cfg, nil
}

func (p *Postgresql) FormatURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", p.Username, p.Password, p.Address, p.Database)
}

func parseConfigPath(fullpath string) (string, string, string) {
	fullSlice := strings.Split(fullpath, "/")

	filename := fullSlice[len(fullSlice)-1]
	fileSlice := strings.Split(filename, ".")

	configPath := strings.Join(fullSlice[:len(fullSlice)-1], "/")
	configName := fileSlice[0]
	configType := fileSlice[1]
	return configPath, configName, configType
}
