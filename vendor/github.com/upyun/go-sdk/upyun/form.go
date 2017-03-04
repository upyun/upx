package upyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FormUploadConfig struct {
	LocalPath      string
	SaveKey        string
	ExpireAfterSec int64
	NotifyUrl      string
	Apps           []map[string]interface{}
	Options        map[string]interface{}
}

type FormUploadResp struct {
	Code      int      `json:"code"`
	Msg       string   `json:"message"`
	Url       string   `json:"url"`
	Timestamp int64    `json:"time"`
	ImgWidth  int      `json:"image-width"`
	ImgHeight int      `json:"image-height"`
	ImgFrames int      `json:"image-frames"`
	ImgType   string   `json:"image-type"`
	Sign      string   `json:"sign"`
	Taskids   []string `json:"task_ids"`
}

func (config *FormUploadConfig) Format() {
	if config.Options == nil {
		config.Options = make(map[string]interface{})
	}
	if config.SaveKey != "" {
		config.Options["save-key"] = config.SaveKey
	}
	if config.NotifyUrl != "" {
		config.Options["notify-url"] = config.NotifyUrl
	}
	if config.ExpireAfterSec > 0 {
		config.Options["expiration"] = time.Now().Unix() + config.ExpireAfterSec
	}
	if len(config.Apps) > 0 {
		config.Options["apps"] = config.Apps
	}
}

func (up *UpYun) FormUpload(config *FormUploadConfig) (*FormUploadResp, error) {
	config.Format()
	config.Options["bucket"] = up.Bucket

	args, err := json.Marshal(config.Options)
	if err != nil {
		return nil, err
	}
	policy := base64ToStr(args)

	formValues := make(map[string]string)
	formValues["policy"] = policy
	formValues["file"] = config.LocalPath

	if up.deprecated {
		formValues["signature"] = up.MakeFormAuth(policy)
	} else {
		sign := &UnifiedAuthConfig{
			Method: "POST",
			Uri:    "/" + up.Bucket,
			Policy: policy,
		}
		if v, ok := config.Options["date"]; ok {
			sign.DateStr = v.(string)
		}
		if v, ok := config.Options["content-md5"]; ok {
			sign.ContentMD5 = v.(string)
		}
		formValues["authorization"] = up.MakeUnifiedAuth(sign)
	}

	endpoint := up.doGetEndpoint("v0.api.upyun.com")
	url := fmt.Sprintf("http://%s/%s", endpoint, up.Bucket)
	resp, err := up.doFormRequest(url, formValues)
	if err != nil {
		return nil, err
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%s", string(b))
	}

	var r FormUploadResp
	err = json.Unmarshal(b, &r)
	return &r, err
}

func (up *UpYun) doFormRequest(url string, formValues map[string]string) (*http.Response, error) {
	formBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(formBody)
	defer formWriter.Close()

	for k, v := range formValues {
		if k != "file" {
			formWriter.WriteField(k, v)
		}
	}

	boundary := formWriter.Boundary()
	bdBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	fpath := formValues["file"]
	fd, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fInfo, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	_, err = formWriter.CreateFormFile("file", filepath.Base(fpath))
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":   "multipart/form-data; boundary=" + boundary,
		"Content-Length": fmt.Sprint(formBody.Len() + int(fInfo.Size()) + bdBuf.Len()),
	}

	body := io.MultiReader(formBody, fd, bdBuf)
	return up.doHTTPRequest("POST", url, headers, body)
}
