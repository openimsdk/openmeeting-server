// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/openimsdk/openmeeting-server/pkg/common/cmd"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/system/program"
	"os"
	"path/filepath"
	"time"
)

const maxRetry = 180

func CheckEtcd(ctx context.Context, config *config.Etcd) error {
	return etcd.Check(ctx, config.Address, "/check_openim_component",
		true,
		etcd.WithDialTimeout(10*time.Second),
		etcd.WithMaxCallSendMsgSize(20*1024*1024),
		etcd.WithUsernameAndPassword(config.Username, config.Password))
}

func CheckRedis(ctx context.Context, config *config.Redis) error {
	return redisutil.Check(ctx, config.Build())
}

func CheckMongo(ctx context.Context, config *config.Mongo) error {
	return mongoutil.Check(ctx, config.Build())
}

func initConfig(configDir string) (*config.Mongo, *config.Redis, *config.Discovery, error) {
	var (
		mongoConfig = &config.Mongo{}
		redisConfig = &config.Redis{}

		discovery = &config.Discovery{}
	)
	err := config.LoadConfig(filepath.Join(configDir, cmd.MongodbConfigFileName), cmd.ConfigEnvPrefixMap[cmd.MongodbConfigFileName], mongoConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.RedisConfigFileName), cmd.ConfigEnvPrefixMap[cmd.RedisConfigFileName], redisConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.DiscoveryConfigFilename), cmd.ConfigEnvPrefixMap[cmd.DiscoveryConfigFilename], discovery)
	if err != nil {
		return nil, nil, nil, err
	}

	return mongoConfig, redisConfig, discovery, nil
}

func main() {
	var index int
	var configDir string
	flag.IntVar(&index, "i", 0, "Index number")
	defaultConfigDir := filepath.Join("..", "..", "..", "..", "..", "config")
	flag.StringVar(&configDir, "c", defaultConfigDir, "Configuration dir")
	flag.Parse()

	fmt.Printf("%s Index: %d, Config Path: %s\n", filepath.Base(os.Args[0]), index, configDir)

	mongoConfig, redisConfig, discoveryConfig, err := initConfig(configDir)
	if err != nil {
		program.ExitWithError(err)
	}

	ctx := context.Background()
	err = performChecks(ctx, mongoConfig, redisConfig, discoveryConfig, maxRetry)
	if err != nil {
		// Assume program.ExitWithError logs the error and exits.
		// Replace with your error handling logic as necessary.
		program.ExitWithError(err)
	}
}

func performChecks(ctx context.Context, mongoConfig *config.Mongo, redisConfig *config.Redis, discovery *config.Discovery, maxRetry int) error {
	checksDone := make(map[string]bool)

	checks := map[string]func(ctx context.Context) error{
		"Mongo": func(ctx context.Context) error {
			return CheckMongo(ctx, mongoConfig)
		},
		"Redis": func(ctx context.Context) error {
			return CheckRedis(ctx, redisConfig)
		},
	}
	if discovery.Enable == "etcd" {
		checks["Etcd"] = func(ctx context.Context) error {
			return CheckEtcd(ctx, &discovery.Etcd)
		}
	}

	for i := 0; i < maxRetry; i++ {
		allSuccess := true
		for name, check := range checks {
			if !checksDone[name] {
				if err := check(ctx); err != nil {
					fmt.Printf("%s check failed: %v\n", name, err)
					allSuccess = false
				} else {
					fmt.Printf("%s check succeeded.\n", name)
					checksDone[name] = true
				}
			}
		}

		if allSuccess {
			fmt.Println("All components checks passed successfully.")
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("not all components checks passed successfully after %d attempts", maxRetry)
}
