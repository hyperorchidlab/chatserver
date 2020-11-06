package cmdcommon

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hyperorchidlab/chatserver/config"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var filenotfind = fmt.Errorf("File Not Found")

func Save2File(data []byte, filename string) error {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		f.Close()
		log.Fatal(err)
	}

	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func OpenAndReadAll(filename string) (data []byte, err error) {
	if !FileExists(filename) {
		return nil, filenotfind
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func CheckPortUsed(iptyp, ipaddr string, port uint16) bool {
	if strings.Contains(strings.ToLower(iptyp), "udp") {
		netaddr := &net.UDPAddr{IP: net.ParseIP(ipaddr), Port: int(port)}
		if c, err := net.ListenUDP(iptyp, netaddr); err != nil {
			return true
		} else {
			c.Close()
			return false
		}
	} else {
		netaddr := &net.TCPAddr{IP: net.ParseIP(ipaddr), Port: int(port)}
		if c, err := net.ListenTCP(iptyp, netaddr); err != nil {
			return true
		} else {
			c.Close()
			return false
		}
	}
}

func GetIPPort(addr string) (ip string, port int, err error) {
	arraddr := strings.Split(addr, ":")
	if len(arraddr) != 2 {
		return "", 0, errors.New("address error")
	}

	ip = arraddr[0]
	port, err = strconv.Atoi(arraddr[1])
	if err != nil {
		return "", 0, err
	}
	if port < 1024 || port > 65535 {
		return "", 0, errors.New("port error")
	}

	if _, err = net.ResolveIPAddr("ip4", ip); err != nil {
		return "", 0, err
	}

	return ip, port, nil
}

func IsProcessCanStarted() (bool, error) {

	cfg := config.PreLoad()

	if cfg == nil {
		return true, nil
	}

	ip, port, err := GetIPPort(cfg.CmdListenPort)
	if err != nil {

		return false, errors.New("Command line listen address error")
	}

	if CheckPortUsed("tcp", ip, uint16(port)) {

		return false, errors.New("Process have started")
	}

	return true, nil
}

func IsProcessStarted() (bool, error) {
	if !config.IsInitialized() {
		return false, errors.New("need to initialize config file first")
	}

	cfg := config.PreLoad()
	if cfg == nil {
		return false, errors.New("load config failed")
	}

	ip, port, err := GetIPPort(cfg.CmdListenPort)
	if err != nil {

		return false, errors.New("Command line listen address error")
	}

	if CheckPortUsed("tcp", ip, uint16(port)) {
		return true, nil
	}

	return false, errors.New("process is not started")
}

func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func GetNowMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}
