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

package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudawan/cloudone_analysis/utility/database/elasticsearch"
	"github.com/cloudawan/cloudone_utility/audit"
	elasticsearchlib "github.com/cloudawan/cloudone_utility/database/elasticsearch"
)

func init() {
	createIndexTemplate()
}

func createIndexTemplate() error {

	tempateBody := `
	{
		"template": "` + indexAuditLogIndex + `",
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
					"Component": {
						"type": "string"
					},
					"Kind": {
						"type": "string"
					},
					"Path": {
						"type": "string"
					},
					"UserName": {
						"type": "string"
					},
					"RemoteAddress": {
						"type": "string",
						"index": "not_analyzed"
					},
					"RemoteHost": {
						"type": "string"
					},
					"CreatedTime": {
						"type": "date",
						"format": "dateOptionalTime"
					},
					"QueryParameterMap" : {
						"properties" : {
						}
					},
					"PathParameterMap" : {
						"properties" : {
						}
					},
					"RequestMethod": {
						"type": "string"
					},
					"RequestURI": {
						"type": "string",
						"index": "not_analyzed"
					},
					"RequestBody" : {
						"type": "string",
						"index": "not_analyzed"
					},
					"RequestHeader" : {
						"properties" : {
						}
					},
					"Description": {
						"type": "string",
						"index": "not_analyzed"
					}
				}
			}
		}
	}
	`

	connection := elasticsearch.ElasticSearchClient.GetConnection()
	request, err := connection.NewRequest("PUT", "/_template/template_"+indexAuditLogIndex, "")
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

func SaveAudit(auditLog *audit.AuditLog, refreshForSearch bool) error {
	id := fmt.Sprintf("%d_%d", auditLog.CreatedTime.Unix(), auditLog.CreatedTime.UnixNano())
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Index(indexAuditLogIndex, auditLog.UserName, id, nil, auditLog)
	if err != nil {
		log.Error(err)
		return err
	} else {
		if refreshForSearch {
			if _, err := connection.Refresh(indexAuditLogIndex); err != nil {
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

func searchAuditLogRawJson(index string, _type string, query interface{}) ([]byte, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	searchResult, err := connection.Search(index, _type, nil, query)
	if err != nil {
		return nil, err
	} else {
		return searchResult.RawJSON, nil
	}
}

func DeleteAuditLogIndex(index string) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.DeleteIndex(index)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetAuditLog(documentType string, id string) (*audit.AuditLog, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	baseResponse, err := connection.Get(indexAuditLogIndex, documentType, id, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		audit := &audit.AuditLog{}
		decoder := json.NewDecoder(bytes.NewReader(*baseResponse.Source))
		decoder.UseNumber()
		err := decoder.Decode(&audit)
		if err != nil {
			log.Error(err)
			return nil, err
		} else {
			return audit, nil
		}
	}
}
