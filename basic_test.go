package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoginAndLogout(t *testing.T) {
	b, err := Upx("login", BUCKET_1, USERNAME, PASSWORD)
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("Welcome to %s, %s!\n", BUCKET_1, USERNAME))

	b, err = Upx("login", BUCKET_2, USERNAME, PASSWORD)
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("Welcome to %s, %s!\n", BUCKET_2, USERNAME))

	b, err = Upx("logout")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("Goodbye %s/%s ~~\n", USERNAME, BUCKET_2))

	b, err = Upx("logout")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("Goodbye %s/%s ~~\n", USERNAME, BUCKET_1))
}

func TestGetInfo(t *testing.T) {
	SetUp()
	defer TearDown()
	pwd, _ := Upx("pwd")
	b, err := Upx("info")
	Nil(t, err)
	s := []string{
		"ServiceName:   " + BUCKET_1,
		"Operator:      " + USERNAME,
		"CurrentDir:    " + strings.TrimRight(string(pwd), "\n"),
		"Usage:         ",
	}
	Equal(t, strings.HasPrefix(string(b), strings.Join(s, "\n")), true)
}

func TestSessionsAndSwitch(t *testing.T) {
	SetUp()
	defer TearDown()
	b, err := Upx("sessions")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("> %s\n", BUCKET_1))

	Upx("login", BUCKET_2, USERNAME, PASSWORD)
	b, err = Upx("sessions")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("  %s\n> %s\n", BUCKET_1, BUCKET_2))

	Upx("switch", BUCKET_1)
	b, err = Upx("sessions")
	Nil(t, err)
	Equal(t, string(b), fmt.Sprintf("> %s\n  %s\n", BUCKET_1, BUCKET_2))

	pwd, _ := Upx("pwd")
	b, err = Upx("info")
	Nil(t, err)
	s := []string{
		"ServiceName:   " + BUCKET_1,
		"Operator:      " + USERNAME,
		"CurrentDir:    " + strings.TrimRight(string(pwd), "\n"),
		"Usage:         ",
	}
	Equal(t, strings.HasPrefix(string(b), strings.Join(s, "\n")), true)
}

//TODO
func TestAuth(t *testing.T) {
}
