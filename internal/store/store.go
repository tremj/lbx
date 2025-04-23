package store

import (
	"bufio"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/parser"
	redisClient "github.com/tremj/lbx/internal/redisClient"
	"gopkg.in/yaml.v3"
)

var (
	ErrAbortConfigReplace = fmt.Errorf("aborting config replacement...")
)

func validateConfigName(cmd *cobra.Command, _ []string) (string, error) {
	configName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", fmt.Errorf("error getting name flag: %v", err)
	}

	redisClient := cmd.Context().Value("redisClient").(*redisClient.RedisClient)
	val, err := redisClient.Get(configName)
	if err == redis.Nil || val != "" {
		return configName, nil
	}

	return "", err
}

func setConfig(cmd *cobra.Command, configName string, config parser.Config) error {
	redisClient := cmd.Context().Value("redisClient").(*redisClient.RedisClient)
	cfg, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	if err = redisClient.Set(configName, string(cfg)); err != nil {
		return fmt.Errorf("error setting config: %v", err)
	}

	return nil
}

func replaceConfig(cmd *cobra.Command, _ []string, configName string, config parser.Config) error {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 3; i++ {
		fmt.Fprint(cmd.OutOrStdout(), "Config already exists... Do you want to replace? [Y/n] ")
		res, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		switch res {
		case "Y":
			return setConfig(cmd, configName, config)
		case "n":
			return ErrAbortConfigReplace
		default:
			continue
		}
	}

	return nil
}

func save(cmd *cobra.Command, args []string, config parser.Config) error {
	configName, err := validateConfigName(cmd, args)
	if err != nil {
		return fmt.Errorf("error validating config name: %v", err)
	}

	if configName != "" {
		return replaceConfig(cmd, args, configName, config)
	}

	return setConfig(cmd, configName, config)
}

func Save(cmd *cobra.Command, args []string) {
	lbConfig, err := parser.RetrieveAndValidate(cmd, args)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
	}
	if err := save(cmd, args, lbConfig); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
		os.Exit(1)
	}
}
