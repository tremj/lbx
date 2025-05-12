package store

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/parser"
	redisclient "github.com/tremj/lbx/internal/redisClient"
	"gopkg.in/yaml.v3"
)

var (
	ErrAbortConfigReplace = fmt.Errorf("aborting config replacement")
)

func validateConfigName(cmd *cobra.Command, _ []string) (string, error) {
	configName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", fmt.Errorf("error getting name flag: %v", err)
	}

	redisClient := cmd.Context().Value("redisClient").(*redisclient.RedisClient)
	val, err := redisClient.Get(configName)
	if errors.Is(err, redis.Nil) || val != "" {
		return configName, err
	}

	return "", err
}

func setConfig(cmd *cobra.Command, configName string, config parser.Config) error {
	redisClient := cmd.Context().Value("redisClient").(*redisclient.RedisClient)
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
		switch strings.TrimSpace(res) {
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

func save(cmd *cobra.Command, args []string) error {
	config, err := parser.RetrieveAndValidate(cmd, args)
	if err != nil {
		return err
	}

	configName, err := validateConfigName(cmd, args)
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("error validating config name: %v", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Saving config...")
	if errors.Is(err, redis.Nil) {
		return setConfig(cmd, configName, config)
	}

	return replaceConfig(cmd, args, configName, config)
}

func Save(cmd *cobra.Command, args []string) {

	if err := save(cmd, args); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
		os.Exit(1)
	}
}

func deleteConfig(cmd *cobra.Command, args []string) error {
	configName, err := validateConfigName(cmd, args)
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("config name %s does not exist", configName)
	} else if err != nil {
		return err
	}
	redisClient := cmd.Context().Value("redisClient").(*redisclient.RedisClient)
	fmt.Fprintln(cmd.OutOrStdout(), "Deleting config...")
	return redisClient.Del(configName)
}

func Delete(cmd *cobra.Command, args []string) {
	if err := deleteConfig(cmd, args); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
		os.Exit(1)
	}
}
