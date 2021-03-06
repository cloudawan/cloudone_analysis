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

package healthcheck

import (
	"github.com/cloudawan/cloudone_analysis/utility/database/elasticsearch"
)

func saveTest(index string, documentType string, id string, jsonMap map[string]interface{}) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Index(index, documentType, id, nil, jsonMap)
	if err != nil {
		log.Error(err)
		return err
	} else {
		return nil
	}
}
