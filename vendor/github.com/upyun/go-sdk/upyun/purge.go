package upyun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	URL "net/url"
	"strings"
	"time"
)

// TODO
func (up *UpYun) Purge(urls []string) (fails []string, err error) {
	purge := "http://purge.upyun.com/purge/"
	date := makeRFC1123Date(time.Now())
	purgeList := strings.Join(urls, "\n")

	headers := map[string]string{
		"Date": date,
		"Authorization": up.MakePurgeAuth(&PurgeAuthConfig{
			PurgeList: purgeList,
			DateStr:   date,
		}),
		"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
	}

	form := make(URL.Values)
	form.Add("purge", purgeList)

	body := strings.NewReader(form.Encode())
	resp, err := up.doHTTPRequest("POST", purge, headers, body)
	if err != nil {
		return fails, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fails, err
	}

	if resp.StatusCode/100 == 2 {
		result := map[string]interface{}{}
		if err := json.Unmarshal(content, &result); err != nil {
			return fails, err
		}
		if it, ok := result["invalid_domain_of_url"]; ok {
			if urls, ok := it.([]interface{}); ok {
				for _, url := range urls {
					fails = append(fails, url.(string))
				}
			}
		}
		return fails, nil
	}

	return nil, fmt.Errorf("purge %d %s", resp.StatusCode, string(content))
}
