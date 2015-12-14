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

package cluster

import (
	"bytes"
	"encoding/json"
	"github.com/cloudawan/cloudone_analysis/utility/database/elasticsearch"
	elasticsearchlib "github.com/cloudawan/cloudone_utility/database/elasticsearch"
)

func init() {
	createIndexTemplate()
}

func createIndexTemplate() error {

	tempateBody := `
	{
		"template": "` + indexClusterSingletonLock + `",
		"mappings": {
			"_default_": {
				"_all": {
					"enabled": true
				},
				"dynamic_templates": [
					{
						"string_fields": {
							"match": "*",
							"match_mapping_type": "string",
							"mapping": {
								"type": "string",
								"index": "not_analyzed",
								"omit_norms": true
							}
						}
					}
				],
				"properties": {
					"id": {
						"type": "string",
						"index": "not_analyzed"
					},
					"firstTimeStamp": {
						"type": "date",
						"format": "dateOptionalTime"
					},
					"lastTimeStamp": {
						"type": "date",
						"format": "dateOptionalTime"
					}
				}
			}
		}
	}
	`

	connection := elasticsearch.ElasticSearchClient.GetConnection()
	request, err := connection.NewRequest("PUT", "/_template/template_"+indexClusterSingletonLock, "")
	if err != nil {
		log.Error(err)
		return err
	}
	request.SetBodyString(tempateBody)
	statusCode, bodyBytes, err := request.Do(nil)
	if err != nil {
		log.Error(err)
		log.Error("statusCode %d", statusCode)
		log.Error(string(bodyBytes))
		return err
	} else {
		//log.Info("statusCode %d", statusCode)
		//log.Info(string(bodyBytes))
	}

	return nil
}

func saveClusterSingletonLock(index string, documentType string, id string, jsonMap map[string]interface{}, refreshForSearch bool) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Index(index, documentType, id, nil, jsonMap)
	if err != nil {
		if err.Error() != notFoundErrorMessage {
			log.Error(err)
		}
		return err
	} else {
		if refreshForSearch {
			if _, err := connection.Refresh(index); err != nil {
				log.Error(err)
				return err
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
}

// Bulk Process
const (
	maxConnection = 5
)

func createBulkProcessor() *elasticsearchlib.BulkProcessor {
	return elasticsearch.ElasticSearchClient.CreateBulkProcessor(maxConnection)
}

func loadClusterSingletonLock(index string, documentType string, id string) (map[string]interface{}, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	baseResponse, err := connection.Get(index, documentType, id, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		jsonMap := make(map[string]interface{})
		decoder := json.NewDecoder(bytes.NewReader(*baseResponse.Source))
		decoder.UseNumber()
		err := decoder.Decode(&jsonMap)
		if err != nil {
			log.Error(err)
			return nil, err
		} else {
			return jsonMap, nil
		}
	}
}
