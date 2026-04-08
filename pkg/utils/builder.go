package utils

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func BuildCanonicalParams(data map[string]interface{}) string {

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		val := data[k]

		strVal := fmt.Sprintf("%v", val)

		encodedVal := url.QueryEscape(strVal)
		parts = append(parts, fmt.Sprintf("%s=%s", k, encodedVal))
	}
	return strings.Join(parts, "&")
}
