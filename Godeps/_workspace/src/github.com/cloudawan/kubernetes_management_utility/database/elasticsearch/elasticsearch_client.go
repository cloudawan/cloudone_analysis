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
	"encoding/json"
	"github.com/cloudawan/kubernetes_management_analysis/Godeps/_workspace/src/code.google.com/p/log4go"
	"github.com/cloudawan/kubernetes_management_analysis/Godeps/_workspace/src/github.com/cloudawan/kubernetes_management/utility/logger"
	elastigo "github.com/cloudawan/kubernetes_management_analysis/Godeps/_workspace/src/github.com/mattbaird/elastigo/lib"
	"strconv"
)

var log log4go.Logger = logger.GetLogger("utility")

type ElasticSearchClient struct {
	connection *elastigo.Conn
	host       []string
	port       int
}

func CreateElasticSearchClient(host []string, port int) *ElasticSearchClient {
	elasticSearchClient := &ElasticSearchClient{nil, host, port}
	elasticSearchClient.GetConnection()
	return elasticSearchClient
}

func (elasticSearchClient *ElasticSearchClient) GetConnection() *elastigo.Conn {
	if elasticSearchClient.connection != nil {
		return elasticSearchClient.connection
	} else {
		c := elastigo.NewConn()
		c.SetHosts(elasticSearchClient.host)
		c.SetPort(strconv.Itoa(elasticSearchClient.port))

		elasticSearchClient.connection = c
		return elasticSearchClient.connection
	}
}

func (elasticSearchClient *ElasticSearchClient) CloseConnection() {
	elasticSearchClient.connection.Close()
	elasticSearchClient.connection = nil
}

func (elasticSearchClient *ElasticSearchClient) Reconnect() {
	elasticSearchClient.CloseConnection()
	elasticSearchClient.GetConnection()
}

// Bulk Process
type BulkProcessor struct {
	bulkIndexer *elastigo.BulkIndexer
}

// WARNING: Due to goelstic's hard code channel size, if the buffered document amount reach 100, it will get stuck
func (bulkProcessor *BulkProcessor) BufferIndex(
	index string, _type string, id string, jsonMap map[string]interface{}) error {
	return bulkProcessor.bulkIndexer.Index(index, _type, id, "", nil, jsonMap, false)
}

func (bulkProcessor *BulkProcessor) FlushAndStopBulkProcessor() {
	bulkProcessor.bulkIndexer.Flush()
	bulkProcessor.bulkIndexer.Stop()
}

func (elasticSearchClient *ElasticSearchClient) CreateBulkProcessor(maxConnection int) *BulkProcessor {
	connection := elasticSearchClient.GetConnection()
	bulkIndexer := connection.NewBulkIndexer(maxConnection)
	bulkIndexer.Start()
	return &BulkProcessor{bulkIndexer}
}

// Get all types for index
func (elasticSearchClient *ElasticSearchClient) GetAllTypeForIndex(index string) ([]string, error) {
	connection := elasticSearchClient.GetConnection()
	request, err := connection.NewRequest("GET", "/"+index, "")
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		_, bodyBytes, err := request.Do(nil)
		if err != nil {
			log.Error(err)
			return nil, err
		} else {
			responseJsonMap := make(map[string]interface{})
			err := json.Unmarshal(bodyBytes, &responseJsonMap)
			if err != nil {
				log.Error(err)
				return nil, err
			} else {
				typeSlice := make([]string, 0)

				mappingMap, _ := responseJsonMap[index].(map[string]interface{})["mappings"].(map[string]interface{})
				for key, _ := range mappingMap {
					if key != "_default_" {
						typeSlice = append(typeSlice, key)
					}
				}

				return typeSlice, nil
			}
		}
	}
}
