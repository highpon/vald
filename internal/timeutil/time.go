//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package timeutil

import "time"

// ParseTime parses string to time.Duration.
func Parse(t string) (time.Duration, error) {
	if t == "" {
		return 0, nil
	}
	dur, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}
	return dur, nil
}

// ParseWithDefault parses string to time.Duration and returns d when the parse failed.
func ParseWithDefault(t string, d time.Duration) time.Duration {
	if t == "" {
		return d
	}

	parsed, err := Parse(t)
	if err != nil {
		return d
	}

	return parsed
}

type DurationString string

func (d DurationString) Duration() (time.Duration, error) {
	return Parse(string(d))
}

func (d DurationString) DurationWithDefault(def time.Duration) time.Duration {
	return ParseWithDefault(string(d), def)
}
