package urlparser

import (
	"net/url"
	"regexp"
	"strings"
)

func Format(rawUrl string, format string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	matches := regexp.MustCompile(`(?i){[a-z0-9:+-]+}`).FindAllString(format, -1)
	if matches == nil {
		return "", nil
	}

	requiredComponents := make(map[string]interface{})
	for _, match := range matches {
		requiredComponents[match] = nil
	}

	replacements := make([]string, 0, len(requiredComponents)*2)
	for placeholder := range requiredComponents {
		component, err := getComponent(parsedUrl, placeholder[1:len(placeholder)-1])
		if err != nil {
			return "", err
		}

		replacements = append(replacements, placeholder, component)
	}

	return strings.NewReplacer(replacements...).Replace(format), nil
}
