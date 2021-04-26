/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"github.com/google/uuid"
)

// Storage is the general storage interface that is used to abstract the
// direct data access.
type Storage interface {
	// GetTimer returns the timer specified by the given uuid. If no timer
	// is found, types.ErrNotFound is returned
	GetTimer(uuid uuid.UUID) (Timer, error)
	// GetLastTimer returns either the last timer of all timers if running
	// is true or the last non-running timer if running is false. If no
	// timer is found types.ErrNotFound is returned or wrapped in the
	// returned error.
	GetLastTimer(running bool) (Timer, error)
	// GetTimers returns all Timers that match the filter. If no timer is
	// found the returned error wraps types.ErrNotFound
	GetTimers(filter Filter) (Timers, error)
	// StoreTimer writes the given timer to the configured data source.
	StoreTimer(timer Timer) error
	// UpdateTimer only allows the stop time to be updated
	// any other changes will be discarded.
	UpdateTimer(timer Timer) error
}
