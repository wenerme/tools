package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/wenerme/tools/pkg/apki"
)

func buildIndexer() (*apki.IndexerServer, error) {
	conf := apki.IndexerConf{}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	if port := os.Getenv("PORT"); port != "" {
		conf.Web.Addr = fmt.Sprintf("0.0.0.0:%v", port)
	}

	svr, err := apki.NewServer(&conf)
	if err != nil {
		return nil, err
	}
	return svr, nil
}
