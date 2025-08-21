package httpclient

import (
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	host       string
	path       string
	queryParam map[string]string
	pathParam  map[string]string
}

func BuildURL(host, path string) URL {
	return URL{
		host:       host,
		path:       path,
		queryParam: make(map[string]string),
		pathParam:  make(map[string]string),
	}
}

func (u URL) WithQueryParam(key, value string) URL {
	u.queryParam[key] = value

	return u
}

func (u URL) WithPathParam(placeHolder, value string) URL {
	u.pathParam[placeHolder] = value

	return u
}

func (u URL) Build() (string, error) {
	path := u.path
	for key, value := range u.pathParam {
		pathParam := "{" + key + "}"
		path = strings.ReplaceAll(path, pathParam, url.PathEscape(value))
	}

	baseURL, err := url.Parse(u.host + path)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %w", err)
	}

	query := url.Values{}
	for key, value := range u.queryParam {
		query.Add(key, value)
	}

	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil

}
