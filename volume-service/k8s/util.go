package k8s

import (
	"fmt"
	"regexp"
	"strings"
)

func makeValidObjectName(prefix string, name string) string {
	name = strings.ToLower(name)
	// change other patterns with hyphen(-)
	replacedName := regexp.MustCompile(`[^a-z0-9\-]+`).ReplaceAllString(name, "-")
	// trim leading or trailing dashes
	replacedName = strings.TrimSuffix(strings.TrimPrefix(replacedName, "-"), "-")
	return fmt.Sprintf("%s-%s", prefix, replacedName)
}
