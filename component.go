package urlparser

import (
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	Scheme      = "scheme"
	Authority   = "authority"
	Auth        = "auth"
	User        = "user"
	Password    = "password"
	HostPort    = "hostPort"
	Host        = "host"
	Tld         = "tld"
	Port        = "port"
	Path        = "path"
	Query       = "query"
	Fragment    = "fragment"
	BasePath    = "basePath"
	File        = "file"
	Ext         = "ext"
	RelativeUrl = "relativeUrl"

	boundedHostPrefix     = Host + ":"
	boundedHostPortPrefix = HostPort + ":"
	boundedPathPrefix     = Path + ":"
	singleQueryPrefix     = Query + ":"
	singleFragmentPrefix  = Fragment + ":"

	halfMaxInt = int(^(uint(0))>>1) / 2
)

func Component(rawUrl string, component string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	return getComponent(parsedUrl, component)
}

func PartialHost(from, count int) string {
	return boundedHostPrefix + strconv.Itoa(from) + ":" + strconv.Itoa(count)
}

func PartialHostFrom(from int) string {
	return boundedHostPrefix + strconv.Itoa(from)
}

func PartialHostPort(from, count int) string {
	return boundedHostPortPrefix + strconv.Itoa(from) + ":" + strconv.Itoa(count)
}

func PartialHostPortFrom(from int) string {
	return boundedHostPortPrefix + strconv.Itoa(from)
}

func PartialPath(from, count int) string {
	return boundedPathPrefix + strconv.Itoa(from) + ":" + strconv.Itoa(count)
}

func PartialPathFrom(from int) string {
	return boundedPathPrefix + strconv.Itoa(from)
}

func SingleQuery(name string) string {
	return singleQueryPrefix + name
}

func SingleFragment(name string) string {
	return singleFragmentPrefix + name
}

func getComponent(parsedUrl *url.URL, component string) (string, error) {
	var result string

	switch true {
	case component == Scheme:
		result = parsedUrl.Scheme
	case component == Authority:
		if parsedUrl.User == nil {
			result = parsedUrl.Host
		} else {
			result = parsedUrl.User.String() + "@" + parsedUrl.Host
		}
	case component == Auth:
		result = parsedUrl.User.String()
	case component == User:
		result = parsedUrl.User.Username()
	case component == Password:
		result, _ = parsedUrl.User.Password()
	case component == HostPort:
		result = parsedUrl.Host
	case component == Host:
		result = parsedUrl.Hostname()
	case component == Tld:
		host := parsedUrl.Hostname()
		if host == "" {
			break
		}

		tld, _ := publicsuffix.PublicSuffix(host)
		_, err := strconv.Atoi(tld)
		if err == nil {
			break // ipv4
		}

		if strings.Contains(host, ":") {
			break // ipv6
		}

		result = tld
	case component == Port:
		result = parsedUrl.Port()
	case component == Path:
		result = parsedUrl.Path
	case component == Query:
		result = parsedUrl.RawQuery
	case component == Fragment:
		result = parsedUrl.Fragment
	case component == BasePath:
		if strings.Contains(parsedUrl.Path, "/") {
			result = parsedUrl.Path[0 : strings.LastIndex(parsedUrl.Path, "/")+1]
		} else {
			result = "/"
		}
	case component == File:
		if strings.Contains(parsedUrl.Path, "/") {
			result = parsedUrl.Path[strings.LastIndex(parsedUrl.Path, "/")+1 : len(parsedUrl.Path)]
		} else {
			result = parsedUrl.Path
		}

		if !strings.Contains(result, ".") {
			result = ""
		}
	case component == Ext:
		result, _ = getComponent(parsedUrl, "file")
		if result != "" {
			result = result[strings.Index(result, ".")+1:]
		}
	case component == RelativeUrl:
		if parsedUrl.Path != "" {
			result += parsedUrl.Path
		}

		if parsedUrl.RawQuery != "" {
			result += "?" + parsedUrl.RawQuery
		}

		if parsedUrl.Fragment != "" {
			result += "#" + parsedUrl.Fragment
		}
	case strings.HasPrefix(component, boundedHostPrefix):
		if parsedUrl.Host == "" {
			break
		}

		host := parsedUrl.Hostname()

		tld, _ := publicsuffix.PublicSuffix(host)
		if tld != "" {
			host = host[0:len(host)-len(tld)] + "~~~"
		}

		hostParts, err := getSliceWithinBounds(component, strings.Split(host, "."))
		if err != nil {
			return "", err
		}

		if len(hostParts) == 0 {
			break
		}

		if hostParts[len(hostParts)-1] == "~~~" {
			hostParts[len(hostParts)-1] = tld
		}

		result = strings.Join(hostParts, ".")
	case strings.HasPrefix(component, boundedHostPortPrefix):
		result, _ = getComponent(parsedUrl, "host"+component[8:])
		port := parsedUrl.Port()
		if port != "" {
			result += ":" + port
		}
	case strings.HasPrefix(component, boundedPathPrefix):
		if parsedUrl.Path == "" {
			break
		}

		splitted := strings.Split(parsedUrl.Path, "/")
		slice, err := getSliceWithinBounds(component, splitted[1:])
		if err != nil {
			return "", err
		}

		if len(slice) == 0 {
			break
		}

		result = "/" + strings.Join(slice, "/")
	case strings.HasPrefix(component, singleQueryPrefix):
		result = parsedUrl.Query().Get(component[6:])
	case strings.HasPrefix(component, singleFragmentPrefix):
		if strings.Contains(parsedUrl.Fragment, "=") {
			parsedUrl, err := url.Parse("?" + parsedUrl.Fragment)
			if err != nil {
				return "", nil
			}

			result = parsedUrl.Query().Get(component[9:])
		}
	default:
		return "", newInvalidComponentErr(component)
	}

	return result, nil
}

func getSliceWithinBounds(component string, orig []string) ([]string, error) {
	var err error
	var start, count int
	var fromEnd = false

	splitted := strings.SplitN(component, ":", 3)
	if splitted[1][0] == '-' {
		fromEnd = true
		splitted[1] = splitted[1][1:]
	}

	start, err = strconv.Atoi(splitted[1])
	if err != nil || start < 0 {
		return nil, newInvalidComponentErr(component)
	}

	if len(splitted) == 3 {
		count, err = strconv.Atoi(splitted[2])
		if err != nil || count < 0 {
			return nil, newInvalidComponentErr(component)
		}
	} else {
		count = halfMaxInt
	}

	if start > len(orig) {
		return []string{}, nil
	}

	if fromEnd {
		return orig[max(0, len(orig)-count) : len(orig)-start], nil
	}

	return orig[start:min(len(orig), start+count)], nil

}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}
