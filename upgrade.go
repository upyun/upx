package upx

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	VERSION_URL         = "https://raw.githubusercontent.com/upyun/upx/master/VERSION"
	DOWNLOAD_URL_PREFIX = "http://collection.b0.upaiyun.com/softwares/upx/upx-"
)

func GetCurrentBinPath() string {
	p, _ := os.Executable()
	return p
}

func GetLatestVersion() (string, error) {
	url := VERSION_URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GetVersion: %v", err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("Get %s: %d", url, resp.StatusCode)
	}

	return strings.TrimSpace(string(content)), nil
}

func DownloadBin(version, binPath string) error {
	fd, err := os.Create(binPath)
	if err != nil {
		return fmt.Errorf("Create %s: %v", binPath, err)
	}
	defer fd.Close()

	url := DOWNLOAD_URL_PREFIX + fmt.Sprintf("%s-%s-%s", runtime.GOOS, runtime.GOARCH, version)
	if runtime.GOOS == "windows" {
		url += ".exe"
	}
	fmt.Print(url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Download %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		content, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Download %s %d: %s", url, resp.StatusCode, string(content))
	}

	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return fmt.Errorf("Download %s: copy %v", url, err)
	}
	return nil
}

func ChmodAndRename(src, dst string) error {
	err := os.Chmod(src, 0755)
	if err != nil {
		return fmt.Errorf("chmod %s: %v", src, err)
	}
	err = os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("rename %s %s: %v", src, dst, err)
	}
	return nil
}

func Upgrade() {
	lv, err := GetLatestVersion()

	if err != nil {
		PrintErrorAndExit("Find Latest Version: %v", err)
	}

	Print("Find Latest Version: %s", lv)
	Print("Current Version: %s", VERSION)

	if lv == VERSION {
		return
	}

	binPath := GetCurrentBinPath()
	tmpBinPath := binPath + ".upgrade"
	err = DownloadBin(lv, tmpBinPath)
	if err != nil {
		PrintErrorAndExit("Download Binary %s: %v", VERSION, err)
	}
	Print("Download Binary %s: OK", VERSION)

	err = ChmodAndRename(tmpBinPath, binPath)
	if err != nil {
		PrintErrorAndExit("Chmod %s: %v", binPath, err)
	}
	PrintErrorAndExit("Chmod %s: OK", binPath)

	return
}
