package api

import "fmt"

func buildEventAPIEndpoint(tag string, token string, isHTTPS bool) string {
	return buildAPIEndpoint("logs-01.loggly.com/inputs/%s/tag/%s/", tag, token, isHTTPS)
}

func buildBulkAPIEndpoint(tag string, token string, isHTTPS bool) string {
	return buildAPIEndpoint("logs-01.loggly.com/bulk/%s/tag/%s/", tag, token, isHTTPS)
}

func buildAPIEndpoint(baseResource string, tag string, token string, isHTTPS bool) string {
	protocol := "http"
	if isHTTPS {
		protocol = "https"
	}

	return fmt.Sprintf(protocol+"://"+baseResource, token, tag)
}
