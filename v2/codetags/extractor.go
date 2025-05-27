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

// Extract extracts comments for special metadata tags. The marker argument
// and group should be unique enough to identify the tags needed.
//
// This function looks for input of the form:
//
//	+<group>:<tag>
//
// This function returns comment lines, with the marker and group prefixes
// removed, grouped by group/name identifier keys. If group is nil, all groups
// are returned. A group of "" represents comment lines without a group prefix.
//
// Example: When called with "+" for 'marker', and "k8s" for group for these comment lines:
//
//	Comment line without marker
//	+k8s:required
//	+listType=set
//	+k8s:format=k8s-long-name
//
// Then this function will return:
//
//	map[TagIdentifier][]string{
//		{Group: "k8s", Name: "required"}: {"required"},
//		{Group: "k8s", Name: "required"}: {"format=k8s-long-name"},
//	}
func Extract(marker string, group *string, lines []string) map[TagIdentifier][]string {
	out := map[TagIdentifier][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, marker) {
			continue
		}
		line = line[len(marker):]

		nameEnd := len(line)
		if idx := strings.IndexAny(line, "(="); idx > 0 {
			nameEnd = idx
		}
		ident := TagIdentifierFromString(line[:nameEnd])

		if group != nil && ident.Group != *group {
			continue
		}
		if ident.Group != "" {
			line = line[len(ident.Group)+1:]
		}

		out[ident] = append(out[ident], line)
	}
	return out
}

// ExtractAndParse combines Extract and Parse.
func ExtractAndParse(marker string, group *string, lines []string) (map[TagIdentifier][]TypedTag, error) {
	out := map[TagIdentifier][]TypedTag{}
	for ident, lines := range Extract(marker, group, lines) {
		for _, line := range lines {
			tag, err := Parse(line)
			if err != nil {
				return nil, err
			}
			out[ident] = append(out[ident], tag)
		}
	}
	return out, nil
}
