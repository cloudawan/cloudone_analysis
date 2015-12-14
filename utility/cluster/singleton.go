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

package cluster

import (
	"errors"
	"github.com/cloudawan/cloudone_analysis/utility/configuration"
	"github.com/cloudawan/cloudone_utility/network"
	"time"
)

var selfIPv4 string
var lockTimeout time.Duration
var wattingAfterBeingCandidate time.Duration

const (
	notFoundErrorMessage = "record not found"
)

func init() {
	var err error
	selfIPv4, err = network.GetFirstNonLoopbackLocalIPv4()
	if err != nil {
		log.Critical(err)
		panic(err)
	}

	singletonLockTimeoutInMilliSecond, ok := configuration.LocalConfiguration.GetInt("singletonLockTimeoutInMilliSecond")
	if ok == false {
		log.Critical("Can't load singletonLockTimeoutInMilliSecond")
		panic("Can't load singletonLockTimeoutInMilliSecond")
	}
	lockTimeout = time.Millisecond * time.Duration(singletonLockTimeoutInMilliSecond)

	singletonLockWaitingAfterBeingCandidateInMilliSecond, ok := configuration.LocalConfiguration.GetInt("singletonLockWaitingAfterBeingCandidateInMilliSecond")
	if ok == false {
		log.Critical("Can't load singletonLockWaitingAfterBeingCandidateInMilliSecond")
		panic("Can't load singletonLockWaitingAfterBeingCandidateInMilliSecond")
	}
	wattingAfterBeingCandidate = time.Millisecond * time.Duration(singletonLockWaitingAfterBeingCandidateInMilliSecond)
}

func IsSelectedAsSingleton(target string) (bool, error) {
	jsonMap, err := loadClusterSingletonLock(indexClusterSingletonLock, typeCloudoneAnalysis, target)
	if err != nil {
		if err.Error() == notFoundErrorMessage {
			// First time so set self as candidate
			if err := setSelfAsCandidate(target, nil); err != nil {
				log.Error(err)
				return false, err
			} else {
				return false, nil
			}
		} else {
			log.Error(err)
			return false, err
		}
	} else {
		lastTimeStampText, ok := jsonMap["lastTimeStamp"].(string)
		if ok == false {
			log.Error("Fail to get field lastTimeStamp")
			return false, errors.New("Fail to get field lastTimeStamp")
		}
		lastTimeStamp, err := time.Parse(time.RFC3339Nano, lastTimeStampText)
		if err != nil {
			log.Error(err)
			return false, err
		}

		now := time.Now()
		if now.Sub(lastTimeStamp) > lockTimeout {
			// Timeout so clear the current candidate and set to self
			if err := setSelfAsCandidate(target, nil); err != nil {
				log.Error(err)
				return false, err
			} else {
				return false, nil
			}
		} else {
			id, ok := jsonMap["id"].(string)
			if ok == false {
				log.Error("Fail to get field id")
				return false, errors.New("Fail to get field id")
			}
			if id == selfIPv4 {
				firstTimeStampText, ok := jsonMap["firstTimeStamp"].(string)
				if ok == false {
					log.Error("Fail to get field firstTimeStamp")
					return false, errors.New("Fail to get field firstTimeStamp")
				}
				firstTimeStamp, err := time.Parse(time.RFC3339Nano, firstTimeStampText)
				if err != nil {
					log.Error(err)
					return false, err
				}

				// Update timestamp
				if err := setSelfAsCandidate(target, &firstTimeStamp); err != nil {
					log.Error(err)
					return false, err
				}

				if lastTimeStamp.Sub(firstTimeStamp) > wattingAfterBeingCandidate {
					// Wait for enough time, so self is selected as the active one
					return true, nil
				} else {
					// Keep waiting as the candidate
					return false, nil
				}
			} else {
				// The candidate is someone else and not timeout yet
				return false, nil
			}
		}
	}
}

func setSelfAsCandidate(target string, firstTimeStamp *time.Time) error {
	lastTimeStamp := time.Now()
	if firstTimeStamp == nil {
		firstTimeStamp = &lastTimeStamp
	}
	jsonMap := make(map[string]interface{})
	jsonMap["id"] = selfIPv4
	jsonMap["firstTimeStamp"] = firstTimeStamp.Format(time.RFC3339Nano)
	jsonMap["lastTimeStamp"] = lastTimeStamp.Format(time.RFC3339Nano)
	if err := saveClusterSingletonLock(indexClusterSingletonLock, typeCloudoneAnalysis, target, jsonMap, false); err != nil {
		log.Error(err)
		return err
	} else {
		return nil
	}
}
