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

package restapi

import (
	"bytes"
	"github.com/cloudawan/cloudone_analysis/audit"
	utilityaudit "github.com/cloudawan/cloudone_utility/audit"
	"github.com/emicklei/go-restful"
	"io/ioutil"
)

func auditLog(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token := req.Request.Header.Get("token")
	requestURI := req.Request.URL.RequestURI()
	method := req.Request.Method
	path := req.SelectedRoutePath()
	queryParameterMap := req.Request.URL.Query()
	pathParameterMap := req.PathParameters()
	remoteAddress := req.Request.RemoteAddr

	requestBody, _ := ioutil.ReadAll(req.Request.Body)
	// Write data back for the later use
	req.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	go func() {
		sendAuditLog(token, requestURI, method, path, string(requestBody), queryParameterMap, pathParameterMap, remoteAddress)
	}()

	chain.ProcessFilter(req, resp)
}

func sendAuditLog(token string, requestURI string, method string, path string, requestBody string, queryParameterMap map[string][]string, pathParameterMap map[string]string, remoteAddress string) {
	// Get cache. If not exsiting, retrieving from authorization server.
	user, err := getCache(token)
	userName := ""
	if err != nil {
		log.Error(err)
		userName = "error_to_get_user"
	}
	if user != nil {
		userName = user.Name
	}

	// Header is not used since the header has no useful information for now
	auditLog := utilityaudit.CreateAuditLog(componentName, path, userName, remoteAddress, queryParameterMap, pathParameterMap, method, requestURI, requestBody, nil)

	err = audit.SaveAudit(auditLog, false)
	if err != nil {
		log.Error("Fail to send audit log with error %s", err)
	}
}
