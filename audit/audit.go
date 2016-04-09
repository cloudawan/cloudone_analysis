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
	"encoding/json"
	"errors"
	"github.com/cloudawan/cloudone_utility/audit"
	"github.com/cloudawan/cloudone_utility/logger"
	"strconv"
	"time"
)

func SearchAuditLog(userName string, from *time.Time, to *time.Time, size int,
	offset int) (returnedAuditLogSlice []audit.AuditLog, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("SearchAuditLog Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedAuditLogSlice = nil
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

	byteSlice, err := searchAuditLogRawJson(indexAuditLogIndex, userName, query)
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
		auditLogSlice := make([]audit.AuditLog, 0)
		for _, result := range resultSlice {
			resultJsonMap, _ := result.(map[string]interface{})
			sourceJsonMap := resultJsonMap["_source"].(map[string]interface{})

			component, _ := sourceJsonMap["Component"].(string)
			kind, _ := sourceJsonMap["Kind"].(string)
			path, _ := sourceJsonMap["Path"].(string)
			userName, _ := sourceJsonMap["UserName"].(string)
			remoteAddress, _ := sourceJsonMap["RemoteAddress"].(string)
			remoteHost, _ := sourceJsonMap["RemoteHost"].(string)
			createdTimeText, _ := sourceJsonMap["CreatedTime"].(string)
			createdTime, _ := time.Parse(time.RFC3339Nano, createdTimeText)
			queryParameterJsonMap, _ := sourceJsonMap["QueryParameterMap"].(map[string]interface{})
			queryParameterMap := make(map[string][]string)
			for key, value := range queryParameterJsonMap {
				queryParameterSlice := make([]string, 0)
				queryParameterJsonSlice, _ := value.([]interface{})
				for _, queryParameterInterface := range queryParameterJsonSlice {
					queryParameter, _ := queryParameterInterface.(string)
					queryParameterSlice = append(queryParameterSlice, queryParameter)
				}
				queryParameterMap[key] = queryParameterSlice
			}
			pathParameterJsonMap, _ := sourceJsonMap["PathParameterMap"].(map[string]interface{})
			pathParameterMap := make(map[string]string)
			for key, value := range pathParameterJsonMap {
				pathParameterMap[key], _ = value.(string)
			}
			requestMethod, _ := sourceJsonMap["RequestMethod"].(string)
			requestURI, _ := sourceJsonMap["RequestURI"].(string)
			requestBody, _ := sourceJsonMap["RequestBody"].(string)
			requestHeaderJsonMap, _ := sourceJsonMap["RequestHeader"].(map[string]interface{})
			requestHeader := make(map[string][]string)
			for key, value := range requestHeaderJsonMap {
				requestHeaderSlice := make([]string, 0)
				requestHeaderJsonSlice, _ := value.([]interface{})
				for _, requestHeaderInterface := range requestHeaderJsonSlice {
					requestHeaderValue, _ := requestHeaderInterface.(string)
					requestHeaderSlice = append(requestHeaderSlice, requestHeaderValue)
				}
				requestHeader[key] = requestHeaderSlice
			}
			description, _ := sourceJsonMap["Description"].(string)

			auditLog := audit.AuditLog{
				component,
				kind,
				path,
				userName,
				remoteAddress,
				remoteHost,
				createdTime,
				queryParameterMap,
				pathParameterMap,
				requestMethod,
				requestURI,
				requestBody,
				requestHeader,
				description,
			}
			auditLogSlice = append(auditLogSlice, auditLog)
		}
		return auditLogSlice, nil
	} else {
		return nil, errors.New("Fail to get with byteSlice " + string(byteSlice))
	}
}
