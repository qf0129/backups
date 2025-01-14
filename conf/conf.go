package conf

import (
	"encoding/json"
	"errors"
	"os"
)

var Conf *config

type config struct {
	Qiniu       qiniuConfig
	Paths       []string
	RotateByDay bool
	SkipFailed  bool
}

type qiniuConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	BucketDir string
}

func LoadConfig(configFile string) error {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &Conf); err != nil {
		return err
	}
	if len(Conf.Paths) == 0 {
		return errors.New("no paths to backup")
	}
	return nil
}
