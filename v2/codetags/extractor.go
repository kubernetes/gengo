/*
Copyright 2025 The Kubernetes Authors.

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

package codetags

import (
	"strings"
)

// Extract extracts comments for special metadata tags. Lines are filtered by a
// string prefix that matching lines must start with.
//
// Example: When called with prefix "+k8s:", lines:
//
//	Comment line without marker
//	+k8s:required
//	+listType=set
//	+k8s:format=k8s-long-name
//
// Then this function will return:
//
//	map[string][]string{
//		"required": {"required"},
//		"required": {"format=k8s-long-name"},
//	}
func Extract(prefix string, lines []string) map[string][]string {
	out := map[string][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		line = line[len(prefix):]

		nameEnd := len(line)
		if idx := strings.IndexAny(line, "(="); idx > 0 {
			nameEnd = idx
		}
		name := line[:nameEnd]
		out[name] = append(out[name], line)
	}
	return out
}

// ExtractAndParse combines Extract and Parse.
func ExtractAndParse(prefix string, lines []string) (map[string][]TypedTag, error) {
	out := map[string][]TypedTag{}
	for name, lines := range Extract(prefix, lines) {
		for _, line := range lines {
			tag, err := Parse(line)
			if err != nil {
				return nil, err
			}
			out[name] = append(out[name], tag)
		}
	}
	return out, nil
}
