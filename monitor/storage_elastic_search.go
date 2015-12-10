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

package monitor

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
	// ElasticSearch doesn't allow to use character '.' in the field name so it should be replaced with '_'
	tempateBody := `
	{
		"template": "` + indexContainerMetricsIndexPrefix + `*",
		"mappings": {
			"_default_": {
				"_all": {
					"enabled":true
				},
				"dynamic_templates":[
					{
						"string_fields":{
							"match":"*",
							"match_mapping_type":"string",
							"mapping":{
								"type":"string",
								"index":"not_analyzed",
								"omit_norms":true
							}
						}
					}
				],
				"properties":{
					"aliases":{
						"type":"string",
						"index":"not_analyzed"
					},
					"name":{
						"type":"string",
						"index":"not_analyzed"
					},
					"namespace":{
						"type":"string",
						"index":"not_analyzed"
					},
					"searchMetaData":{
						"properties":{
							"containerName":{
								"type":"string",
								"index":"not_analyzed"
							},
							"namespace":{
								"type":"string",
								"index":"not_analyzed"
							},
							"podName":{
								"type":"string",
								"index":"not_analyzed"
							},
							"replicationControllerName":{
								"type":"string",
								"index":"not_analyzed"
							}
						}
					},
					"spec":{
						"properties":{
							"labels": {
								"properties":{
									"io_kubernetes_pod_name": {
										"type":"string",
										"index":"not_analyzed"
									},
									"io_kubernetes_pod_terminationGracePeriod": {
										"type":"string",
										"index":"not_analyzed"
									}
								}
							},
							"cpu":{
								"properties":{
									"limit":{
										"type":"long"
									},
									"mask":{
										"type":"string",
										"index":"not_analyzed"
									},
									"max_limit":{
										"type":"long"
									}
								}
							},
							"creation_time":{
								"type":"date",
								"format":"dateOptionalTime"
							},
							"has_cpu":{
								"type":"boolean"
							},
							"has_diskio":{
								"type":"boolean"
							},
							"has_filesystem":{
								"type":"boolean"
							},
							"has_memory":{
								"type":"boolean"
							},
							"has_network":{
								"type":"boolean"
							},
							"memory":{
								"properties":{
									"limit":{
										"type":"double"
									},
									"swap_limit":{
										"type":"double"
									}
								}
							},
							"has_custom_metrics": {
								"type":"boolean"
							},
							"image": {
								"type":"string",
								"index":"not_analyzed"
							}
						}
					},
					"stats":{
						"properties":{
							"cpu":{
								"properties":{
									"load_average":{
										"type":"long"
									},
									"usage":{
										"properties":{
											"per_cpu_usage":{
												"type":"long"
											},
											"system":{
												"type":"long"
											},
											"total":{
												"type":"long"
											},
											"user":{
												"type":"long"
											}
										}
									}
								}
							},
							"diskio":{
								"properties":{
									"io_service_bytes":{
										"properties":{
											"major":{
												"type":"long"
											},
											"minor":{
												"type":"long"
											},
											"stats":{
												"properties":{
													"Async":{
														"type":"long"
													},
													"Read":{
														"type":"long"
													},
													"Sync":{
														"type":"long"
													},
													"Total":{
														"type":"long"
													},
													"Write":{
														"type":"long"
													}
												}
											}
										}
									},
									"io_serviced":{
										"properties":{
											"major":{
												"type":"long"
											},
											"minor":{
												"type":"long"
											},
											"stats":{
												"properties":{
													"Async":{
														"type":"long"
													},
													"Read":{
														"type":"long"
													},
													"Sync":{
														"type":"long"
													},
													"Total":{
														"type":"long"
													},
													"Write":{
														"type":"long"
													}
												}
											}
										}
									}
								}
							},
							"filesystem":{
								"properties":{
									"available":{
										"type":"long"
									},
									"capacity":{
										"type":"long"
									},
									"device":{
										"type":"string",
										"index":"not_analyzed"
									},
									"io_in_progress":{
										"type":"long"
									},
									"io_time":{
										"type":"long"
									},
									"read_time":{
										"type":"long"
									},
									"reads_completed":{
										"type":"long"
									},
									"reads_merged":{
										"type":"long"
									},
									"sectors_read":{
										"type":"long"
									},
									"sectors_written":{
										"type":"long"
									},
									"usage":{
										"type":"long"
									},
									"weighted_io_time":{
										"type":"long"
									},
									"write_time":{
										"type":"long"
									},
									"writes_completed":{
										"type":"long"
									},
									"writes_merged":{
										"type":"long"
									}
								}
							},
							"memory":{
								"properties":{
									"container_data":{
										"properties":{
											"pgfault":{
												"type":"long"
											},
											"pgmajfault":{
												"type":"long"
											}
										}
									},
									"hierarchical_data":{
										"properties":{
											"pgfault":{
												"type":"long"
											},
											"pgmajfault":{
												"type":"long"
											}
										}
									},
									"failcnt":{
										"type":"long"
									},
									"usage":{
										"type":"long"
									},
									"working_set":{
										"type":"long"
									}
								}
							},
							"network":{
								"properties":{
									"name":{
										"type":"string",
										"index":"not_analyzed"
									},
									"rx_bytes":{
										"type":"long"
									},
									"rx_dropped":{
										"type":"long"
									},
									"rx_errors":{
										"type":"long"
									},
									"rx_packets":{
										"type":"long"
									},
									"tx_bytes":{
										"type":"long"
									},
									"tx_dropped":{
										"type":"long"
									},
									"tx_errors":{
										"type":"long"
									},
									"tx_packets":{
										"type":"long"
									},
									"tcp":{
										"properties":{
											"Established":{
												"type":"long"
											},
											"SynSent":{
												"type":"long"
											},
											"SynRecv":{
												"type":"long"
											},
											"FinWait1":{
												"type":"long"
											},
											"FinWait2":{
												"type":"long"
											},
											"TimeWait":{
												"type":"long"
											},
											"Close":{
												"type":"long"
											},
											"CloseWait":{
												"type":"long"
											},
											"LastAck":{
												"type":"long"
											},
											"Listen":{
												"type":"long"
											},
											"Closing":{
												"type":"long"
											}
										}
									},
									"tcp6":{
										"properties":{
											"Established":{
												"type":"long"
											},
											"SynSent":{
												"type":"long"
											},
											"SynRecv":{
												"type":"long"
											},
											"FinWait1":{
												"type":"long"
											},
											"FinWait2":{
												"type":"long"
											},
											"TimeWait":{
												"type":"long"
											},
											"Close":{
												"type":"long"
											},
											"CloseWait":{
												"type":"long"
											},
											"LastAck":{
												"type":"long"
											},
											"Listen":{
												"type":"long"
											},
											"Closing":{
												"type":"long"
											}
										}
									}
								}
							},
							"task_stats":{
								"properties":{
									"nr_io_wait":{
										"type":"long"
									},
									"nr_running":{
										"type":"long"
									},
									"nr_sleeping":{
										"type":"long"
									},
									"nr_stopped":{
										"type":"long"
									},
									"nr_uninterruptible":{
										"type":"long"
									}
								}
							},
							"timestamp":{
								"type":"date",
								"format":"dateOptionalTime"
							}
						}
					}
				}
			}
		}
	}
	`

	connection := elasticsearch.ElasticSearchClient.GetConnection()
	request, err := connection.NewRequest("PUT", "/_template/template_"+indexContainerMetricsIndexPrefix, "")
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

func saveContainerRecord(index string, documentType string, id string, jsonMap map[string]interface{}) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.Index(index, documentType, id, nil, jsonMap)
	if err != nil {
		log.Error(err)
		return err
	} else {
		return nil
	}
}

// Bulk Process
const (
	maxConnection = 5
)

func createBulkProcessor() *elasticsearchlib.BulkProcessor {
	return elasticsearch.ElasticSearchClient.CreateBulkProcessor(maxConnection)
}

func SearchContainerRecordRawJson(index string, _type string, query interface{}) ([]byte, error) {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	searchResult, err := connection.Search(index, _type, nil, query)
	if err != nil {
		return nil, err
	} else {
		return searchResult.RawJSON, nil
	}
}

func DeleteContainerRecordIndex(index string) error {
	connection := elasticsearch.ElasticSearchClient.GetConnection()
	_, err := connection.DeleteIndex(index)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetAllReplicationControllerNameInNameSpace(namespace string) ([]string, error) {
	documentTypeSlice, err := elasticsearch.ElasticSearchClient.GetAllTypeForIndex(getDocumentIndex(namespace))
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		replicationControllerNameSlice := make([]string, 0)
		for _, documentType := range documentTypeSlice {
			replicationControllerName := getReplicationControllerNameFromDocumentType(documentType)
			replicationControllerNameSlice = append(replicationControllerNameSlice, replicationControllerName)
		}
		return replicationControllerNameSlice, nil
	}
}

func GetContainerRecord(index string, documentType string, id string) (map[string]interface{}, error) {
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
