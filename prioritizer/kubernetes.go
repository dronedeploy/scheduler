package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	apiHost       = "127.0.0.1:8001"
	patchEndpoint = "/api/v1/namespaces/default/pods/%s"
	podsEndpoint  = "/api/v1/pods"
)

func getUnscheduledPods() ([]Pod, error) {
	var podList PodList
	unscheduledPods := make([]Pod, 0)

	v := url.Values{}
	v.Set("fieldSelector", "spec.nodeName=")

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     apiHost,
			Path:     podsEndpoint,
			RawQuery: v.Encode(),
			Scheme:   "http",
		},
	}
	request.Header.Set("Accept", "application/json, */*")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return unscheduledPods, err
	}
	err = json.NewDecoder(resp.Body).Decode(&podList)
	if err != nil {
		return unscheduledPods, err
	}

	for _, pod := range podList.Items {
		if pod.Metadata.Annotations["scheduler.alpha.kubernetes.io/name"] == schedulerName {
			unscheduledPods = append(unscheduledPods, pod)
		}
	}

	return unscheduledPods, nil
}

func patchPriority(pod *Pod, priority string) error {
	patch := Pod{
		Metadata: Metadata{
			Annotations: map[string]string{"k8s_priority": priority},
			Labels:      map[string]string{"priority": priority},
		},
	}

	var b []byte
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(patch)
	if err != nil {
		return err
	}

	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPatch,
		URL: &url.URL{
			Host:   apiHost,
			Path:   fmt.Sprintf(patchEndpoint, pod.Metadata.Name),
			Scheme: "http",
		},
	}
	request.Header.Set("Content-Type", "application/merge-patch+json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Patch: Unexpected HTTP status code" + resp.Status)
	}
	return nil
}
