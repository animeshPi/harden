package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Result struct {
	ID         string
	PrettyName string
	Family     string
	VersionID  string
}

var debianDistros = map[string]bool{
	"debian": true, "ubuntu": true, "linuxmint": true,
	"pop": true, "kali": true, "raspbian": true,
}

var rhelDistros = map[string]bool{
	"rhel": true, "fedora": true, "centos": true,
	"rocky": true, "almalinux": true, "amzn": true,
}

func Detect() (Result, error) {

	switch runtime.GOOS {

	case "linux":
		return DetectFromFile("/etc/os-release")

	case "windows":
		return detectWindows(), nil

	default:
		fmt.Println("unsupported os: ", runtime.GOOS)
		return Result{}, errors.New("unsupported os")
	}
}

func detectWindows() Result {
	// runtime.GOOS only gives "windows", so best we can do
	// without syscall is to set family + ID.
	return Result{
		ID:     "windows",
		Family: "windows",
		// Windows version detection requires syscall; left minimal intentionally
		VersionID: runtime.GOARCH,
	}
}

func DetectFromFile(path string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	info := map[string]string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.Trim(strings.TrimSpace(line[idx+1:]), `"`)
			info[key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("scan failed: %w", err)
	}

	id := strings.ToLower(info["ID"])
	idLike := strings.ToLower(info["ID_LIKE"])
	versionID := info["VERSION_ID"]

	family := detectLinuxFamily(id, idLike)

	return Result{
		ID:         id,
		PrettyName: info["PRETTY_NAME"],
		Family:     family,
		VersionID:  versionID,
	}, nil
}

func detectLinuxFamily(id, idLike string) string {
	if debianDistros[id] {
		return "debian"
	}
	if rhelDistros[id] {
		return "rhel"
	}

	for _, like := range strings.Fields(idLike) {
		if debianDistros[like] {
			return "debian"
		}
		if rhelDistros[like] {
			return "rhel"
		}
	}

	return "unknown"
}
