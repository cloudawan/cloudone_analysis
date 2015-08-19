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

/*
import (
	"fmt"
	"strconv"
	"testing"
	"time"
)
*/
/*
func TestGetContainerRecordIndex(t *testing.T) {
	timestamp := time.Now().Add(-0 * time.Minute)

	fmt.Println(GetContainerRecordIndex(
		getDocumentIndex("default"),
		getDocumentType("kube-dns-v6"),
		getDocumentID("kube-dns-v6-4ui44", "etcd", timestamp)))
}
*/
/*
func TestDeleteContainerRecordIndex(t *testing.T) {
	fmt.Println(DeleteContainerRecordIndex(getDocumentIndex("nucleus")))
}
*/
/*
func TestCreateIndexTemplate(t *testing.T) {
	fmt.Println(createIndexTemplate())
}
*/
/*
func TestLoadContainerRecordRawJson(t *testing.T) {
	current := time.Now()

	nodeAmount := 60
	from := current.Add(-2 * time.Minute)
	to := current.Add(-1 * time.Minute)

	duration := int(to.Sub(from).Seconds())

	gte := from.UTC().Format(time.RFC3339Nano)
	lte := to.UTC().Format(time.RFC3339Nano)

	intervalInSecond := strconv.Itoa(int(duration / nodeAmount))

	fmt.Println(intervalInSecond)
	fmt.Println(from)
	fmt.Println(to)

	query := `
	{
		"query": {
			"range" : {
				"stats.timestamp" : {
					"gte": "` + gte + `",
					"lte": "` + lte + `",
					"time_zone": "+0:00"
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
									"minimum_diskio_io_serviced_stats_total" : { "sum" : { "field" : "stats.diskio.io_serviced.stats.Total" } }
								}
							}
						}
					}
				}
			}
		}
	}
	`
	result, err := SearchContainerRecordRawJson(getDocumentIndex("default"), getDocumentType("cassandra"), query)
	fmt.Println(err)
	fmt.Println(string(result))
}
*/
