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

package elasticsearch

import (
	"github.com/cloudawan/kubernetes_management_analysis/utility/configuration"
	"github.com/cloudawan/kubernetes_management_analysis/utility/logger"
	"github.com/cloudawan/kubernetes_management_utility/database/elasticsearch"
)

var log = logger.GetLogManager().GetLogger("utility")

var ElasticSearchClient *elasticsearch.ElasticSearchClient

func init() {
	elasticsearchHost, ok := configuration.LocalConfiguration.GetStringSlice("elasticsearchHost")
	if ok == false {
		log.Critical("Can't load elasticsearchHost")
		panic("Can't load elasticsearchHost")
	}

	elasticsearchPort, ok := configuration.LocalConfiguration.GetInt("elasticsearchPort")
	if ok == false {
		log.Critical("Can't load elasticsearchPort")
		panic("Can't load elasticsearchPort")
	}

	ElasticSearchClient = elasticsearch.CreateElasticSearchClient(elasticsearchHost, elasticsearchPort)
}
