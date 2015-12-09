// Copyright 2015 CloudAwan LLC
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

package configuration

import (
	"errors"
	analysisLogger "github.com/cloudawan/cloudone_analysis/utility/logger"
	"github.com/cloudawan/cloudone_utility/configuration"
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"strings"
	"time"
)

var log = analysisLogger.GetLogManager().GetLogger("utility")

var configurationContent = `
{
	"certificate": "/etc/cloudone_analysis/development_cert.pem",
	"key": "/etc/cloudone_analysis/development_key.pem",
	"elasticsearchHost": ["127.0.0.1"],
	"elasticsearchPort": 9200,
	"kubeapiHostAndPort": ["127.0.0.1:8080"],
	"kubeapiHealthCheckTimeoutInMilliSecond": 1000
}
`

var LocalConfiguration *configuration.Configuration
var KubeapiHealthCheckTimeoutInMilliSecond = 1000

func init() {
	var err error
	LocalConfiguration, err = configuration.CreateConfiguration("cloudone_analysis", configurationContent)
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func GetAvailableKubeapiHostAndPort() (returnedHost string, returnedPort int, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedHost = ""
			returnedPort = 0
			returnedError = err.(error)
			log.Error("GetAvailableKubeapiHostAndPort Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
		}
	}()

	kubeapiHostAndPortSlice, ok := LocalConfiguration.GetStringSlice("kubeapiHostAndPort")
	if ok == false {
		log.Error("Fail to get configuration kubeapiHostAndPort")
		return "", 0, errors.New("Fail to get configuration kubeapiHostAndPort")
	}

	kubeapiHealthCheckTimeoutInMilliSecond, ok := LocalConfiguration.GetInt("kubeapiHealthCheckTimeoutInMilliSecond")
	if ok == false {
		kubeapiHealthCheckTimeoutInMilliSecond = KubeapiHealthCheckTimeoutInMilliSecond
	}

	for _, kubeapiHostAndPort := range kubeapiHostAndPortSlice {
		result, err := restclient.HealthCheck("http://"+kubeapiHostAndPort,
			time.Duration(kubeapiHealthCheckTimeoutInMilliSecond)*time.Millisecond)

		if result {
			splitSlice := strings.Split(kubeapiHostAndPort, ":")
			host := splitSlice[0]
			port, err := strconv.Atoi(splitSlice[1])
			if err != nil {
				log.Error(err)
				return "", 0, err
			}
			return host, port, nil
		} else {
			if err != nil {
				log.Error(err)
			}
		}
	}

	log.Error("No available host and port")
	return "", 0, errors.New("No available host and port")
}
