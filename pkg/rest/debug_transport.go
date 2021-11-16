package rest

import (
	"bytes"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
)

type debugTransport struct {
	delegatedRoundTripper http.RoundTripper
}

func (d *debugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(request.Body)
		klog.V(10).Infof("Request Body: %s", string(bodyBytes))
		request.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}
	response, err := d.delegatedRoundTripper.RoundTrip(request)
	if err == nil && response.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		bodyStr := string(bodyBytes)
		klog.V(10).Infof("Response Body: %s", bodyStr)
		response.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}
	return response, err
}
