package upyun

import (
	"fmt"
	"sort"
	"strings"
)

type RESTAuthConfig struct {
	Method    string
	Uri       string
	DateStr   string
	LengthStr string
}

type PurgeAuthConfig struct {
	PurgeList string
	DateStr   string
}

type UnifiedAuthConfig struct {
	Method     string
	Uri        string
	DateStr    string
	Policy     string
	ContentMD5 string
}

func (u *UpYun) MakeRESTAuth(config *RESTAuthConfig) string {
	sign := []string{
		config.Method,
		config.Uri,
		config.DateStr,
		config.LengthStr,
		u.Password,
	}
	return "UpYun " + u.Operator + ":" + md5Str(strings.Join(sign, "&"))
}

func (u *UpYun) MakePurgeAuth(config *PurgeAuthConfig) string {
	sign := []string{
		config.PurgeList,
		u.Bucket,
		config.DateStr,
		u.Password,
	}
	return "UpYun " + u.Bucket + ":" + u.Operator + ":" + md5Str(strings.Join(sign, "&"))
}

func (u *UpYun) MakeFormAuth(policy string) string {
	return md5Str(base64ToStr([]byte(policy)) + "&" + u.Secret)
}

func (u *UpYun) MakeProcessAuth(kwargs map[string]string) string {
	keys := []string{}
	for k, _ := range kwargs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	auth := ""
	for _, k := range keys {
		auth += k + kwargs[k]
	}
	return fmt.Sprintf("UpYun %s:%s", u.Operator, md5Str(u.Operator+auth+u.Password))
}

func (u *UpYun) MakeUnifiedAuth(config *UnifiedAuthConfig) string {
	sign := []string{
		config.Method,
		config.Uri,
		config.DateStr,
		config.Policy,
		config.ContentMD5,
	}
	signNoEmpty := []string{}
	for _, v := range sign {
		if v != "" {
			signNoEmpty = append(signNoEmpty, v)
		}
	}
	signStr := base64ToStr(hmacSha1(u.Password, []byte(strings.Join(signNoEmpty, "&"))))
	return "UpYun " + u.Operator + ":" + signStr
}
