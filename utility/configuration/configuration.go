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
	"github.com/cloudawan/kubernetes_management_analysis/utility/logger"
	"github.com/cloudawan/kubernetes_management_utility/configuration"
)

var log = logger.GetLog("utility")

var configurationContent = `
{
	"certificate": "/etc/kubernetes_management_analysis/development_cert.pem",
	"key": "/etc/kubernetes_management_analysis/development_key.pem",
	"elasticsearchHost": ["127.0.0.1"],
	"elasticsearchPort": 9200,
	"kubeapiHost": "127.0.0.1",
	"kubeapiPort": 8080
}
`

var LocalConfiguration *configuration.Configuration

func init() {
	var err error
	LocalConfiguration, err = configuration.CreateConfiguration("kubernetes_management_analysis", configurationContent)
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}
