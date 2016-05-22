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
	"errors"
	"github.com/cloudawan/cloudone_analysis/control"
	"github.com/cloudawan/cloudone_utility/jsonparse"
	"github.com/cloudawan/cloudone_utility/logger"
	"strconv"
	"time"
)

func RecordHistoricalReplicationController(kubeApiServerEndPoint string, kubeApiServerToken string, namespace string, replicationControllerName string) (returnedReplicationControllerContainerRecordSlice []map[string]interface{}, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("RecordHistoricalReplicationController Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
			returnedReplicationControllerContainerRecordSlice = nil
		}
	}()

	podNameSlice, err := control.GetAllPodNameBelongToReplicationController(kubeApiServerEndPoint, kubeApiServerToken, namespace, replicationControllerName)
	if err != nil {
		log.Error("Fail to get all pod name belong to the replication controller with endpoint %s, token: %s, namespace: %s, replication controller name: %s", kubeApiServerEndPoint, kubeApiServerToken, namespace, replicationControllerName)
		return nil, err
	}

	errorBuffer := bytes.Buffer{}
	errorBuffer.WriteString("The following container has error: ")
	errorHappened := false

	replicationControllerContainerRecordSlice := make([]map[string]interface{}, 0)
	for _, podName := range podNameSlice {
		podContainerRecordSlice, err := RecordHistoricalPod(kubeApiServerEndPoint, kubeApiServerToken, namespace, replicationControllerName, podName)
		if err != nil {
			errorHappened = true
			log.Error("RecordHistoricalPod error %s", err)
			errorBuffer.WriteString("RecordHistoricalPod error " + err.Error())
		} else {
			replicationControllerContainerRecordSlice = append(replicationControllerContainerRecordSlice, podContainerRecordSlice...)
		}
	}

	if errorHappened {
		log.Error("Fail to get all container inofrmation with endpoint %s, token: %s, namespace: %s, error %s", kubeApiServerEndPoint, kubeApiServerToken, namespace, errorBuffer.String())
		return nil, errors.New(errorBuffer.String())
	} else {
		return replicationControllerContainerRecordSlice, nil
	}
}

func GetAllHistoricalReplicationControllerMetrics(namespace string,
	aggregationAmount int, from time.Time, to time.Time) (returnedJsonMap map[string]interface{}, returnedError error) {
	replicationControllerNameSlice, err := GetAllReplicationControllerNameInNameSpace(namespace)
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		namespaceJsonMap := make(map[string]interface{})
		for _, replicationControllerName := range replicationControllerNameSlice {
			replicationControllerJsonMap, err := GetHistoricalReplicationControllerMetrics(namespace,
				replicationControllerName, aggregationAmount, from, to)
			if err != nil {
				log.Error(err)
			} else {
				namespaceJsonMap[replicationControllerName] = replicationControllerJsonMap
			}
		}
		return namespaceJsonMap, nil
	}
}

func GetHistoricalReplicationControllerMetrics(namespace string,
	replicationControllerName string, aggregationAmount int, from time.Time,
	to time.Time) (returnedJsonMap map[string]interface{}, returnedError error) {
	byteSlice, err := searchHistoricalReplicationControllerMetrics(namespace,
		replicationControllerName, aggregationAmount, from, to)
	if err != nil {
		return nil, err
	} else {
		jsonMap := make(map[string]interface{})
		decoder := json.NewDecoder(bytes.NewReader(byteSlice))
		decoder.UseNumber()
		err := decoder.Decode(&jsonMap)
		if err != nil {
			return nil, err
		} else {
			timestampSlice := make([]string, 0)
			replicationControllerJsonMap := make(map[string]interface{})

			timeBucketSlice, _ := jsonMap["aggregations"].(map[string]interface{})["aggregation_time_interval"].(map[string]interface{})["buckets"].([]interface{})
			timeBucketAmount := len(timeBucketSlice)
			for timeIndex, timeBucket := range timeBucketSlice {
				timestamp, _ := timeBucket.(map[string]interface{})["key_as_string"].(string)
				timestampSlice = append(timestampSlice, timestamp)

				podBucketSlice, _ := timeBucket.(map[string]interface{})["aggregation_pod"].(map[string]interface{})["buckets"].([]interface{})
				for _, podBucket := range podBucketSlice {
					podName, _ := podBucket.(map[string]interface{})["key"].(string)
					containerBucketSlice, _ := podBucket.(map[string]interface{})["aggregation_container"].(map[string]interface{})["buckets"].([]interface{})

					podJsonMap, _ := replicationControllerJsonMap[podName].(map[string]interface{})
					if podJsonMap == nil {
						podJsonMap = make(map[string]interface{})
					}

					for _, containerBucket := range containerBucketSlice {
						containerName, _ := containerBucket.(map[string]interface{})["key"].(string)

						documentCount, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["doc_count"])
						minimumCpuUsageTotal, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_cpu_usage_total"].(map[string]interface{})["value"])
						averageMemoryUsage, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["average_memory_usage"].(map[string]interface{})["value"])
						minimumDiskioIoServiceBytesStatsTotal, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_diskio_io_service_bytes_stats_total"].(map[string]interface{})["value"])
						minimumDiskioIoServicedStatsTotal, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_diskio_io_serviced_stats_total"].(map[string]interface{})["value"])
						minimumNetworkRxPackets, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_network_rx_packets"].(map[string]interface{})["value"])
						minimumNetworkTxPackets, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_network_tx_packets"].(map[string]interface{})["value"])
						minimumNetworkRxBytes, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_network_rx_bytes"].(map[string]interface{})["value"])
						minimumNetworkTxBytes, _ := jsonparse.ConvertToInt64(containerBucket.(map[string]interface{})["minimum_network_tx_bytes"].(map[string]interface{})["value"])

						containerJsonMap, _ := podJsonMap[containerName].(map[string]interface{})
						if containerJsonMap == nil {
							containerJsonMap = make(map[string]interface{})
						}

						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "documentCountSlice", int64(documentCount))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumCpuUsageTotalSlice", int64(minimumCpuUsageTotal))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "averageMemoryUsageSlice", int64(averageMemoryUsage))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumDiskioIoServiceBytesStatsTotalSlice", int64(minimumDiskioIoServiceBytesStatsTotal))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumDiskioIoServicedStatsTotalSlice", int64(minimumDiskioIoServicedStatsTotal))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumNetworkRxPacketsSlice", int64(minimumNetworkRxPackets))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumNetworkTxPacketsSlice", int64(minimumNetworkTxPackets))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumNetworkRxBytesSlice", int64(minimumNetworkRxBytes))
						appendToSliceInJsonMap(timeBucketAmount, timeIndex, containerJsonMap, "minimumNetworkTxBytesSlice", int64(minimumNetworkTxBytes))

						podJsonMap[containerName] = containerJsonMap
					}

					replicationControllerJsonMap[podName] = podJsonMap
				}
			}

			// Interpolate the hole
			for podName, _ := range replicationControllerJsonMap {
				for containerName, _ := range replicationControllerJsonMap[podName].(map[string]interface{}) {
					for metricsName, _ := range replicationControllerJsonMap[podName].(map[string]interface{})[containerName].(map[string]interface{}) {
						fillTheNullDataWithInterpolationForInt64Slice(replicationControllerJsonMap[podName].(map[string]interface{})[containerName].(map[string]interface{})[metricsName].([]interface{}))
					}
				}
			}

			// Add timestamp
			replicationControllerJsonMap["timestamp"] = timestampSlice

			return replicationControllerJsonMap, nil
		}
	}
}

func appendToSliceInJsonMap(timeBucketAmount int, timeIndex int, jsonMap map[string]interface{}, sliceName string, value int64) {
	slice, ok := jsonMap[sliceName].([]interface{})
	if ok == false {
		slice = make([]interface{}, timeBucketAmount)
	}
	slice[timeIndex] = value
	jsonMap[sliceName] = slice
}

func fillTheNullDataWithInterpolationForInt64Slice(dataSlice []interface{}) {
	for i := 0; i < len(dataSlice); i++ {
		if dataSlice[i] == nil {
			// Find left
			var left interface{} = nil
			for j := i; j >= 0; j-- {
				if dataSlice[j] != nil {
					left = dataSlice[j]
					break
				}
			}
			// Find right
			var right interface{} = nil
			for j := i; j < len(dataSlice); j++ {
				if dataSlice[j] != nil {
					right = dataSlice[j]
					break
				}
			}
			// If both exists, use interpolation. Otherwise single side.
			if left != nil && right != nil {
				leftValue, _ := left.(int64)
				rightValue, _ := right.(int64)
				dataSlice[i] = (leftValue + rightValue) / 2
			} else if left != nil {
				dataSlice[i] = left
			} else {
				dataSlice[i] = right
			}
		}
	}
}

func searchHistoricalReplicationControllerMetrics(
	namespace string, replicationControllerName string, aggregationAmount int,
	from time.Time, to time.Time) (returnedByteSlice []byte, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("SearchHistoricalReplicationControllerMetrics Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedByteSlice = nil
			returnedError = err.(error)
		}
	}()

	if from.After(to) {
		return nil, errors.New("From " + from.String() + " can't be after to " + to.String())
	}

	duration := int(to.Sub(from).Seconds())
	gte := from.UTC().Format(time.RFC3339Nano)
	lte := to.UTC().Format(time.RFC3339Nano)
	intervalInSecond := strconv.Itoa(int(duration / aggregationAmount))

	query := `
	{
		"query": {		
			"range" : {
				"stats.timestamp" : {
					"gte": "` + gte + `",
					"lte": "` + lte + `",
					"time_zone": "+00:00"
				}
			}
	    },
		"size": 0,
		"aggregations": {
			"aggregation_time_interval": { 
				"date_histogram": {
					"field": "stats.timestamp",
					"interval" : "` + intervalInSecond + `s"
				},
				"aggregations": {
					"aggregation_pod": {	
						"terms": {
							"field": "searchMetaData.podName"
						},
						"aggregations" : {
							"aggregation_container": {	
								"terms": {
									"field": "searchMetaData.containerName"
								},
								"aggregations" : {
									"minimum_cpu_usage_total" : { "min" : { "field" : "stats.cpu.usage.total" } },
									"average_memory_usage" : { "avg" : { "field" : "stats.memory.usage" } },
									"minimum_network_rx_bytes" : { "min" : { "field" : "stats.network.rx_bytes" } },
									"minimum_network_tx_bytes" : { "min" : { "field" : "stats.network.tx_bytes" } },
									"minimum_network_rx_packets" : { "min" : { "field" : "stats.network.rx_packets" } },
									"minimum_network_tx_packets" : { "min" : { "field" : "stats.network.tx_packets" } },
									"minimum_diskio_io_service_bytes_stats_total" : { "min" : { "field" : "stats.diskio.io_service_bytes.stats.Total" } },
									"minimum_diskio_io_serviced_stats_total" : { "min" : { "field" : "stats.diskio.io_serviced.stats.Total" } }
								}
							}
						}
					}
				}
			}
		}
	}
	`
	return SearchContainerRecordRawJson(getDocumentIndex(namespace), getDocumentType(replicationControllerName), query)
}
