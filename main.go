package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ncruces/zenity"
)

// The directory containing the static files.
const staticDir = "build"

func main() {
	http.HandleFunc("/", handleStaticFiles)
	log.Println("Starting server on port 8602")
	go func() {
		if err := http.ListenAndServe(":8602", nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	if !isChromeInstalled() {
		showInstallChromeDialog()
	} else {
		if err := openBrowser("http://localhost:8602"); err != nil {
			log.Println("Failed to open browser:", err)
		}
	}

	select {}
}

func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
}

func isChromeInstalled() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("where", "chrome.exe")
		if err := cmd.Run(); err != nil {
			return false
		}
		return true
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("ls", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome")
		if err := cmd.Run(); err != nil {
			return false
		}
		return true
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("which", "google-chrome")
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("which", "chromium-browser")
			if err := cmd.Run(); err != nil {
				return false
			}
		}
		return true
	}
	return false
}

func detectLang() string {
	var lang string
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("powershell", "Get-Culture | select -exp Name").Output()
		if err != nil {
			log.Fatal(err)
		}
		lang = strings.TrimSpace(string(out))
	default:
		lang = os.Getenv("LANG")
	}
	return strings.ToLower(lang[:2])
}

func showInstallChromeDialog() {
	lang := detectLang()
	messages := map[string]map[string]string{
		"en": {
			"info":     "Chrome not detected. Do you want to download it?",
			"cancel":   "Cancel",
			"download": "Download",
			"title":    "Download Chrome",
		},
		"zh": {
			"info":     "未检测到chrome，是否去下载？",
			"cancel":   "取消",
			"download": "下载",
			"title":    "下载谷歌浏览器",
		},
	}
	fmt.Print(lang)
	if _, ok := messages[lang]; !ok {
		lang = "en"
	}
	err := zenity.Question(messages[lang]["info"], zenity.Title(messages[lang]["title"]), zenity.OKLabel(messages[lang]["download"]), zenity.CancelLabel(messages[lang]["cancel"]), zenity.Icon("./icon/ScratchDesktop.png"))
	if err == nil {
		var cmd string
		switch runtime.GOOS {
		case "darwin":
			cmd = "open"
		case "windows":
			cmd = "start"
		default:
			cmd = "xdg-open"
		}
		err := exec.Command(cmd, "https://www.google.com/chrome/").Start()
		if err != nil {
			return
		}
	}
	os.Exit(0)
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
