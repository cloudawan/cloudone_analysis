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

package build

import (
	"bytes"
	"encoding/json"
	"github.com/cloudawan/cloudone_analysis/utility/database/elasticsearch"
	"github.com/cloudawan/cloudone_utility/build"
	elasticsearchlib "github.com/cloudawan/cloudone_utility/database/elasticsearch"
	"strings"
)

func init() {
	createIndexTemplate()
}

func createIndexTemplate() error {

	tempateBody := `
	{
		"template": "` + indexBuildLogIndexPrefix + `*",
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
					"ImageInformation": {
						"type": "string"
					},
					"Version": {
						"type": "string"
					},
					"VersionInfo" : {
						"properties" : {
						}
					},
					"CreatedTime": {
						"type": "date",
						"format": "dateOptionalTime"
					},
					"Content": {
						"type": "string"
					}
				}
			}
		}
	}
	`

	connection := elasticsearch.ElasticSearchClient.GetConnection()
	request, err := connection.NewRequest("PUT", "/_template/template_"+indexBuildLogIndexPrefix, "")
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

func checkFormatForElasticSearchData(buildLog *build.BuildLog) {
	if buildLog.VersionInfo != nil {
		for key, value := range buildLog.VersionInfo {
			if strings.Contains(key, ".") {
				newKey := strings.Replace(key, ".", "_", -1)
				buildLog.VersionInfo[newKey] = value
				delete(buildLog.VersionInfo, key)
			}
		}
	}
}

func getIndexName(imageInformation string) string {
	return indexBuildLogIndexPrefix + strings.ToLower(imageInformation)
}

func SaveBuildLog(buildLog *build.BuildLog, refreshForSearch bool) error {
	checkFormatForElasticSearchData(buildLog)
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Index(getIndexName(buildLog.ImageInformation), indexBuildLogType, buildLog.Version, nil, buildLog)
	if err != nil {
		log.Debug(buildLog)
		log.Error(err)
		return err
	} else {
		if refreshForSearch {
			if _, err := connection.Refresh(getIndexName(buildLog.ImageInformation)); err != nil {
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

func searchBuildLogRawJson(index string, _type string, query interface{}) ([]byte, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	searchResult, err := connection.Search(index, _type, nil, query)
	if err != nil {
		return nil, err
	} else {
		return searchResult.RawJSON, nil
	}
}

func DeleteBuildLogBelongingToImageInformation(imageInformation string) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.DeleteIndex(getIndexName(imageInformation))
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteBuildLog(imageInformation string, version string) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Delete(getIndexName(imageInformation), indexBuildLogType, version, nil)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetBuildLog(imageInformation string, version string) (*build.BuildLog, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	baseResponse, err := connection.Get(getIndexName(imageInformation), indexBuildLogType, version, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		buildLog := &build.BuildLog{}
		decoder := json.NewDecoder(bytes.NewReader(*baseResponse.Source))
		decoder.UseNumber()
		err := decoder.Decode(&buildLog)
		if err != nil {
			log.Error(err)
			return nil, err
		} else {
			return buildLog, nil
		}
	}
}
