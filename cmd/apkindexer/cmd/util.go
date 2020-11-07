package cmd

import (
	"github.com/spf13/viper"
	"github.com/wenerme/tools/pkg/apki"
)

func buildIndexer() (*apki.IndexerServer, error) {
	conf := apki.IndexerConf{}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	svr, err := apki.NewServer(&conf)
	if err != nil {
		return nil, err
	}
	return svr, nil
}
