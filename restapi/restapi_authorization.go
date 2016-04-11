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
	"errors"
	"github.com/cloudawan/cloudone_analysis/utility/configuration"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"github.com/emicklei/go-restful"
	"strconv"
	"time"
)

const (
	componentName      = "cloudone_analysis"
	cacheCheckInterval = time.Minute
	cacheTTL           = cacheCheckInterval * 60
)

func init() {
	periodicallyCleanCache()
}

var closed bool = false

func Close() {
	closed = true
}

func periodicallyCleanCache() {
	go func() {
		for {
			if closed {
				break
			}

			rbac.CheckCacheTimeout()

			time.Sleep(cacheCheckInterval)
		}
	}()
}

func getCache(token string) (*rbac.User, error) {
	// Get from cache first
	user := rbac.GetCache(token)
	if user == nil {
		// Not exist. Ask the authorization server.
		cloudoneProtocol, ok := configuration.LocalConfiguration.GetString("cloudoneProtocol")
		if ok == false {
			log.Error("Unable to get configuration cloudoneProtocol")
			return nil, errors.New("Unable to get configuration cloudoneProtocol")
		}

		cloudoneHost, ok := configuration.LocalConfiguration.GetString("cloudoneHost")
		if ok == false {
			log.Error("Unable to get configuration cloudoneHost")
			return nil, errors.New("Unable to get configuration cloudoneHost")
		}

		cloudonePort, ok := configuration.LocalConfiguration.GetInt("cloudonePort")
		if ok == false {
			log.Error("Unable to get configuration cloudonePort")
			return nil, errors.New("Unable to get configuration cloudonePort")
		}

		url := cloudoneProtocol + "://" + cloudoneHost + ":" + strconv.Itoa(cloudonePort) +
			"/api/v1/authorizations/tokens/" + token + "/components/" + componentName
		user := &rbac.User{}
		_, err := restclient.RequestGetWithStructure(url, &user, nil)
		if err != nil {
			log.Debug(err)
			return nil, err
		} else {
			// Set Cache
			rbac.SetCache(token, user, cacheTTL)
			log.Info("Cache user %s", user.Name)

			return user, nil
		}
	} else {
		return user, nil
	}
}

func authorize(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	token := req.Request.Header.Get("token")

	// Get cache. If not exsiting, retrieving from authorization server.
	user, err := getCache(token)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Unable to authorize due to error"
		jsonMap["ErrorMessage"] = err.Error()
		resp.WriteHeaderAndJson(500, jsonMap, "{}")
		return
	}

	// Verify
	if user != nil {
		authorized := false
		if user.HasPermission(componentName, req.Request.Method, req.SelectedRoutePath()) {
			// Resource check
			namespace := req.PathParameter("namespace")
			namespacePass := false
			if namespace != "" {
				if user.HasResource(componentName, "/namespaces/"+namespace) {
					namespacePass = true
				}
			} else {
				namespacePass = true
			}
			if namespacePass {
				authorized = true
			}
		}

		if authorized {
			chain.ProcessFilter(req, resp)
		} else {
			jsonMap := make(map[string]interface{})
			jsonMap["Error"] = "Not Authorized"
			jsonMap["Format"] = "Put correct token in the header token"
			resp.WriteHeaderAndJson(401, jsonMap, "{}")
		}
	} else {
		// Cache doesn't exist
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Token doesn't exist"
		jsonMap["ErrorMessage"] = "Token is incorrect or expired. Please get token with username and password again."
		resp.WriteHeaderAndJson(401, jsonMap, "{}")
	}
}
