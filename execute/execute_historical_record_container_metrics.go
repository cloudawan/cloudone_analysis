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

package execute

import (
	"github.com/cloudawan/cloudone_analysis/control"
	"github.com/cloudawan/cloudone_analysis/monitor"
	"github.com/cloudawan/cloudone_analysis/utility/configuration"
	"github.com/cloudawan/cloudone_utility/logger"
	"time"
)

func loopHistoricalRecordContainerMetrics(ticker *time.Ticker, checkingInterval time.Duration) {
	for {
		select {
		case <-ticker.C:
			// Historical record
			periodicalRunHistoricalRecordContainerMetrics()
		case <-quitChannel:
			ticker.Stop()
			log.Info("Loop historical record container metrics quit")
			return
		}
	}
}

func periodicalRunHistoricalRecordContainerMetrics() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("periodicalRunHistoricalRecordEvent Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
		}
	}()

	kubeapiHost, ok := configuration.LocalConfiguration.GetString("kubeapiHost")
	if ok == false {
		log.Error("Fail to get configuration kubeapiHost")
		return
	}
	kubeapiPort, ok := configuration.LocalConfiguration.GetInt("kubeapiPort")
	if ok == false {
		log.Error("Fail to get configuration kubeapiPort")
		return
	}
	namespaceNameSlice, err := control.GetAllNamespaceName(kubeapiHost, kubeapiPort)
	if err != nil {
		log.Error(err)
		return
	}
	for _, namespaceName := range namespaceNameSlice {
		replicationControllerNameSlice, err := control.GetAllReplicationControllerName(kubeapiHost, kubeapiPort, namespaceName)
		if err != nil {
			log.Error(err)
		} else {
			for _, replicationControllerName := range replicationControllerNameSlice {
				err := monitor.RecordHistoricalReplicationController(kubeapiHost, kubeapiPort, namespaceName, replicationControllerName)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
