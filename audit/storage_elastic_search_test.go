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

/*
import (
	"fmt"
	"github.com/cloudawan/cloudone_utility/audit"
	"testing"
	"time"
)
*/
/*
func TestGetAuditLog(t *testing.T) {
	fmt.Println(GetAuditLog("admin", "1460177558_1460177558072337023"))
}
*/
/*
func TestDeleteAuditLogIndex(t *testing.T) {
	fmt.Println(DeleteAuditLogIndex(indexAuditLogIndex))
}
*/
/*
func TestSearchAuditLog(t *testing.T) {
	slice, err := SearchAuditLog("*", nil, nil, 10, 0)
	fmt.Println(err)
	fmt.Println(len(slice))
	fmt.Println(slice)
}
*/
/*
func TestSaveAuditLog(t *testing.T) {
	userName := "admin"
	queryParameterMap := make(map[string][]string)
	queryParameter := make([]string, 0)
	queryParameter = append(queryParameter, "false")
	queryParameterMap["acknowledge"] = queryParameter
	pathParameterMap := make(map[string]string)
	pathParameterMap["namespace"] = "default"
	requestMethod := "GET"
	requestURI := "https://192.168.0.17:8082/api/v1/historicalevents/default?acknowledge=false&size=100&offset=0"
	requestBody := ""
	requestHeader := make(map[string][]string)
	requestSlice := make([]string, 0)
	requestSlice = append(requestSlice, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	requestHeader["token"] = requestSlice

	auditLog := audit.CreateAuditLog(
		"cloudone",
		"/api/v1/historicalevents/{namespace}",
		userName,
		"127.0.0.1",
		queryParameterMap,
		pathParameterMap,
		requestMethod,
		requestURI,
		requestBody,
		requestHeader)

	fmt.Println(SaveAudit(auditLog, false))
}
*/
