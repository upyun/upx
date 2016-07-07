package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type userInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Bucket   string `json:"bucket"`
	CurDir   string `sjon:"curdir"`
}

type Config struct {
	Idx   int         `json:"user_idx"`
	Users []*userInfo `json:"users"`
}

func (c *Config) Load(fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	var b []byte
	if b, err = ioutil.ReadAll(fd); err == nil {
		if b, err = base64.StdEncoding.DecodeString(string(b)); err == nil {
			err = json.Unmarshal(b, c)
		}
	}

	return err
}

func (c *Config) Save(fname string) error {
	if len(c.Users) == 0 {
		os.Remove(fname)
		return nil
	}
	fd, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	var b []byte
	if b, err = json.Marshal(c); err == nil {
		s := base64.StdEncoding.EncodeToString(b)
		_, err = fd.WriteString(s)
	}

	return err
}

func (c *Config) GetCurUser() *userInfo {
	if c.Idx >= 0 && c.Idx < len(c.Users) {
		return c.Users[c.Idx]
	}
	return nil
}

func (c *Config) UpdateUserInfo(u *userInfo) {
	c.Idx = -1
	for k, v := range c.Users {
		if v.Bucket == u.Bucket {
			c.Idx = k
			break
		}
	}
	if c.Idx == -1 {
		c.Idx = len(c.Users)
		c.Users = append(c.Users, u)
	} else {
		c.Users[c.Idx] = u
	}
}

func (c *Config) SwitchBucket(bucket string) error {
	for k, v := range c.Users {
		if v.Bucket == bucket {
			c.Idx = k
			return nil
		}
	}
	return errors.New("no such bucket")
}

func (c *Config) RemoveBucket() error {
	if c.Idx >= 0 && c.Idx < len(c.Users) {
		c.Users = append(c.Users[0:c.Idx], c.Users[c.Idx+1:]...)
		c.Idx = 0
		return nil
	}
	return errors.New("no such bucket")
}
