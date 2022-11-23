package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/viper"
)

type DB struct {
	DSN     string `mapstructure:"dsn"`
	Timeout string `mapstructure:"timeout"`
}

type Config struct {
	DB *DB `mapstructure:"indexStorageDB"`
}

const (
	FileName = "application.yaml"
)

var (
	ErrSSMFetch = errors.New("config: not exist Parameter at ssm")
)

func MustNew() *Config {
	var (
		cfg = &Config{}
	)

	viper.SetConfigType("yaml")

	yaml, err := getYamlFromSSM(fmt.Sprintf("/github/drakejin/crawler/%s", FileName))
	if err != nil {
		panic(err)
	}
	if err = viper.ReadConfig(strings.NewReader(yaml)); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func getYamlFromSSM(ssmPath string) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	res, err := ssm.New(sess).GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(ssmPath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	p := res.Parameter
	if p == nil {
		return "", ErrSSMFetch
	}

	return *p.Value, nil
}
