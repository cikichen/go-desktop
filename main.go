package mc_desktop

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/ncruces/zenity"
)

// The directory containing the static files.
const staticDir = "build"

func main() {
	http.HandleFunc("/", handleStaticFiles)
	log.Println("Starting server on port 8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	if !isChromeInstalled() {
		showInstallChromeDialog()
	} else {
		if err := openBrowser("http://localhost:8080"); err != nil {
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
			showInstallChromeDialog()
			return false
		}
		return true
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("which", "google-chrome")
		if err := cmd.Run(); err != nil {
			showInstallChromeDialog()
			return false
		}
		return true
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("which", "google-chrome")
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("which", "chromium-browser")
			if err := cmd.Run(); err != nil {
				showInstallChromeDialog()
				return false
			}
		}
		return true
	}
	return false
}

func showInstallChromeDialog() {
	title := "Chrome not found"
	message := "Chrome browser is required to run this application. Please install Chrome and try again."
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		err := zenity.Warning(message)
		if err != nil {
			log.Println("Failed to show dialog:", err)
			return
		}
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("zenity", "--warning", "--text", message, "--title", title)
		if err := cmd.Run(); err != nil {
			log.Println("Failed to show dialog:", err)
			return
		}
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		if !isChromeInstalled() {
			showInstallChromeDialog()
			return fmt.Errorf("chrome not installed")
		}
		err = exec.Command("google-chrome", url).Start()
	case "windows":
		if !isChromeInstalled() {
			showInstallChromeDialog()
			return fmt.Errorf("chrome not installed")
		}
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		if !isChromeInstalled() {
			showInstallChromeDialog()
			return fmt.Errorf("chrome not installed")
		}
		err = exec.Command("open", url).Start()
	default:
		showInstallChromeDialog()
		return fmt.Errorf("unsupported platform")
	}
	return err
}
