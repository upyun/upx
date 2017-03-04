package upyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CommitTasksConfig struct {
	AppName   string
	Accept    string
	Source    string
	NotifyUrl string
	Tasks     []interface{}
}

func (up *UpYun) CommitTasks(config *CommitTasksConfig) (taskIds []string, err error) {
	b, err := json.Marshal(config.Tasks)
	if err != nil {
		return nil, err
	}

	kwargs := map[string]string{
		"app_name":   config.AppName,
		"tasks":      base64ToStr(b),
		"notify_url": config.NotifyUrl,

		// for naga
		"source": config.Source,
	}
	if config.Accept != "" {
		kwargs["accept"] = config.Accept
	}

	err = up.doProcessRequest("POST", "/pretreatment/", kwargs, &taskIds)
	return
}

func (up *UpYun) GetProgress(taskIds []string) (result map[string]int, err error) {
	kwargs := map[string]string{
		"task_ids": strings.Join(taskIds, ","),
	}
	v := map[string]map[string]int{}
	err = up.doProcessRequest("GET", "/status/", kwargs, &v)
	if err != nil {
		return
	}

	if r, ok := v["tasks"]; ok {
		return r, err
	}
	return nil, fmt.Errorf("no tasks")
}

func (up *UpYun) GetResult(taskIds []string) (result map[string]interface{}, err error) {
	kwargs := map[string]string{
		"task_ids": strings.Join(taskIds, ","),
	}
	v := map[string]map[string]interface{}{}
	err = up.doProcessRequest("GET", "/result/", kwargs, &v)
	if err != nil {
		return
	}

	if r, ok := v["tasks"]; ok {
		return r, err
	}
	return nil, fmt.Errorf("no tasks")
}

func (up *UpYun) doProcessRequest(method, uri string,
	kwargs map[string]string, v interface{}) error {
	if _, ok := kwargs["service"]; !ok {
		kwargs["service"] = up.Bucket
	}

	if method == "GET" {
		uri = addQueryToUri(uri, kwargs)
	}

	headers := make(map[string]string)
	headers["Date"] = makeRFC1123Date(time.Now())
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	if up.deprecated {
		headers["Authorization"] = up.MakeProcessAuth(kwargs)
	} else {
		headers["Authorization"] = up.MakeUnifiedAuth(&UnifiedAuthConfig{
			Method:  method,
			Uri:     uri,
			DateStr: headers["Date"],
		})
	}

	var resp *http.Response
	var err error
	endpoint := up.doGetEndpoint("p0.api.upyun.com")
	rawurl := fmt.Sprintf("http://%s%s", endpoint, uri)
	switch method {
	case "GET":
		resp, err = up.doHTTPRequest(method, rawurl, headers, nil)
	case "POST":
		payload := encodeQueryToPayload(kwargs)
		resp, err = up.doHTTPRequest(method, rawurl, headers, bytes.NewBufferString(payload))
	default:
		return fmt.Errorf("Unknown method")
	}

	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("%d %s", resp.StatusCode, string(b))
	}

	return json.Unmarshal(b, v)
}
