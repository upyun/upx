package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const (
	LOGIN    = true
	NO_LOGIN = false
)

type Config struct {
	SessionId int        `json:"user_idx"`
	Sessions  []*Session `json:"users"`
}

func (c *Config) PopCurrent() {
	if c.SessionId == -1 {
		c.SessionId = 0
	}

	c.Sessions = append(c.Sessions[0:c.SessionId], c.Sessions[c.SessionId+1:]...)
	c.SessionId = 0
}

func (c *Config) Insert(sess *Session) {
	for idx, s := range c.Sessions {
		if s.Bucket == sess.Bucket && s.Operator == sess.Operator {
			c.Sessions[idx] = sess
			c.SessionId = idx
			return
		}
	}
	c.Sessions = append(c.Sessions, sess)
	c.SessionId = len(c.Sessions) - 1
}

var (
	confname string
	config   *Config
)

func makeAuthStr(bucket, operator, password string) (string, error) {
	sess := &Session{
		Bucket:   bucket,
		Operator: operator,
		Password: password,
		CWD:      "/",
	}
	if err := sess.Init(); err != nil {
		return "", err
	}

	s := []string{bucket, operator, password}

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return hashEncode(base64.StdEncoding.EncodeToString(b)), nil
}

func authStrToConfig(auth string) error {
	data, err := base64.StdEncoding.DecodeString(hashEncode(auth))
	if err != nil {
		return err
	}
	ss := []string{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	if len(ss) == 3 {
		session = &Session{
			Bucket:   ss[0],
			Operator: ss[1],
			Password: ss[2],
			CWD:      "/",
		}
		if err := session.Init(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid auth string")
	}
	return nil
}

func readConfigFromFile(login bool) {
	if confname == "" {
		confname = getConfigName()
	}

	b, err := ioutil.ReadFile(confname)
	if err != nil {
		os.RemoveAll(confname)
		if os.IsNotExist(err) && login == NO_LOGIN {
			return
		}
		PrintErrorAndExit("read config: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(hashEncode(string(b)))
	if err != nil {
		os.RemoveAll(confname)
		PrintErrorAndExit("read config: %v", err)
	}

	config = &Config{SessionId: -1}
	if err := json.Unmarshal(data, config); err != nil {
		os.RemoveAll(confname)
		PrintErrorAndExit("read config: %v", err)
	}

	if config.SessionId != -1 && config.SessionId < len(config.Sessions) {
		session = config.Sessions[config.SessionId]
		if login == LOGIN {
			if err := session.Init(); err != nil {
				config.PopCurrent()
				PrintErrorAndExit("Log in: %v", err)
			}
		}
	} else {
		if login == LOGIN {
			PrintErrorAndExit("Log in to UpYun first")
		}
	}
}

func saveConfigToFile() {
	if confname == "" {
		confname = getConfigName()
	}

	b, err := json.Marshal(config)
	if err != nil {
		PrintErrorAndExit("save config: %v", err)
	}
	s := hashEncode(base64.StdEncoding.EncodeToString(b))

	fd, err := os.Create(confname)
	if err != nil {
		PrintErrorAndExit("save config: %v", err)
	}
	defer fd.Close()
	_, err = fd.WriteString(s)

	if err != nil {
		PrintErrorAndExit("save config: %v", err)
	}
}

func getConfigName() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), ".upx.cfg")
	}
	return filepath.Join(os.Getenv("HOME"), ".upx.cfg")
}

func hashEncode(s string) string {
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		switch {
		case r[i] >= 'a' && r[i] <= 'z':
			r[i] += 'A' - 'a'
		case r[i] >= 'A' && r[i] <= 'Z':
			r[i] += 'a' - 'A'
		case r[i] >= '0' && r[i] <= '9':
			r[i] = '9' - r[i] + '0'
		}
	}
	return string(r)
}
