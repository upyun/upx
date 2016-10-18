package main

import (
	"encoding/base64"
	"encoding/json"
)

func genAuth(bucket, username, password string) string {
	v := map[string]interface{}{
		"bucket":   bucket,
		"username": username,
		"password": password,
	}

	if b, err := json.Marshal(v); err == nil {
		return base64.StdEncoding.EncodeToString(b)
	}
	return ""
}
