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
	"encoding/json"
	"errors"
	"github.com/cloudawan/cloudone_utility/build"
	"github.com/cloudawan/cloudone_utility/logger"
	"strconv"
	"time"
)

func SearchBuildLog(imageInformation string, from *time.Time, to *time.Time, size int,
	offset int) (returnedBuildLogSlice []build.BuildLog, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("SearchBuildLog Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedBuildLogSlice = nil
			returnedError = err.(error)
		}
	}()

	if from != nil && to != nil && from.After(*to) {
		return nil, errors.New("From " + from.String() + " can't be after to " + to.String())
	}

	var queryField string
	if from == nil && to != nil {
		lte := to.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"CreatedTime" : {
					"lte": "` + lte + `",
					"time_zone": "+0:00"
				}
			}
	    },`
	} else if from != nil && to == nil {
		gte := from.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"CreatedTime" : {
					"gte": "` + gte + `",
					"time_zone": "+0:00"
				}
			}
	    },`
	} else if from != nil && to != nil {
		lte := to.UTC().Format(time.RFC3339Nano)
		gte := from.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"CreatedTime" : {
					"lte": "` + lte + `",
					"gte": "` + gte + `",
					"time_zone": "+0:00"
				}
			}
	    }`
	} else {
		queryField = ``
	}

	query := `
	{
		"query": {
			"filtered": {
				` + queryField + `
			}
		},
		"sort" : [
	 		{ 
				"CreatedTime" : "desc"
			}
    		],
		"size": ` + strconv.Itoa(size) + `,
		"from": ` + strconv.Itoa(offset) + `
	}
	`

	byteSlice, err := searchBuildLogRawJson(getIndexName(imageInformation), indexBuildLogType, query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(byteSlice, &jsonMap); err != nil {
		log.Error(err)
		return nil, err
	}

	resultSlice, ok := jsonMap["hits"].(map[string]interface{})["hits"].([]interface{})
	if ok {
		buildLogSlice := make([]build.BuildLog, 0)
		for _, result := range resultSlice {
			resultJsonMap, _ := result.(map[string]interface{})
			sourceJsonMap := resultJsonMap["_source"].(map[string]interface{})

			imageInformation, _ := sourceJsonMap["ImageInformation"].(string)
			version, _ := sourceJsonMap["Version"].(string)

			versionInfoJsonMap, _ := sourceJsonMap["VersionInfo"].(map[string]interface{})
			versionInfoMap := make(map[string]string)
			for key, value := range versionInfoJsonMap {
				versionInfoMap[key], _ = value.(string)
			}

			createdTimeText, _ := sourceJsonMap["CreatedTime"].(string)
			createdTime, _ := time.Parse(time.RFC3339Nano, createdTimeText)
			content, _ := sourceJsonMap["Content"].(string)

			buildLog := build.BuildLog{
				imageInformation,
				version,
				versionInfoMap,
				createdTime,
				content,
			}
			buildLogSlice = append(buildLogSlice, buildLog)
		}
		return buildLogSlice, nil
	} else {
		log.Error("Fail to get with byteSlice %s", string(byteSlice))
		return nil, errors.New("Fail to get with byteSlice " + string(byteSlice))
	}
}
