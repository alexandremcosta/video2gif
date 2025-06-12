package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var supportedVideoExts = map[string]bool{
	".mp4": true, ".mov": true, ".mkv": true, ".avi": true, ".webm": true,
}

func main() {
	// Ensure ffmpeg and ffprobe paths are accessible even from GUI .app
	_ = os.Setenv("PATH", os.Getenv("PATH")+":/opt/homebrew/bin:/usr/local/bin")

	myApp := app.New()
	w := myApp.NewWindow("video2gif")
	w.Resize(fyne.NewSize(500, 320))

	var videoFile string
	var outputPath string

	startSecEntry := widget.NewEntry()
	endSecEntry := widget.NewEntry()

	fileLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("")

	viewButton := widget.NewButton("View GIF", func() {
		if outputPath != "" {
			openWithDefaultBrowser(outputPath)
		}
	})
	viewButton.Hide()

	fileButton := widget.NewButton("Choose Video File", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if r != nil {
				videoFile = r.URI().Path()
				ext := strings.ToLower(filepath.Ext(videoFile))
				if !supportedVideoExts[ext] {
					dialog.ShowError(fmt.Errorf("Unsupported file format: %s. Please use .mov, .mp4, .webm, etc.", ext), w)
					videoFile = ""
					fileLabel.SetText("No file selected")
					return
				}
				fileLabel.SetText(filepath.Base(videoFile))

				durationSec, err := getVideoDuration(videoFile)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				startSecEntry.SetText("1")
				endSecEntry.SetText(strconv.Itoa(int(durationSec)))
			}
		}, w)
	})

	runButton := widget.NewButton("Generate GIF", func() {
		startSec, err1 := strconv.Atoi(startSecEntry.Text)
		endSec, err2 := strconv.Atoi(endSecEntry.Text)

		if videoFile == "" || err1 != nil || err2 != nil || endSec <= startSec {
			dialog.ShowError(fmt.Errorf("Check inputs. File must be selected and seconds must be valid."), w)
			return
		}

		duration := endSec - startSec
		startTime := fmt.Sprintf("00:00:%02d.000", startSec)
		output := strings.TrimSuffix(videoFile, filepath.Ext(videoFile)) + ".gif"

		cmd := exec.Command("ffmpeg",
			"-ss", startTime,
			"-t", fmt.Sprintf("%d", duration),
			"-i", videoFile,
			"-filter_complex",
			"[0:v]fps=10,scale=960:-1:flags=lanczos,split[x][z];[z]palettegen=stats_mode=diff:max_colors=128[p];[x][p]paletteuse=dither=sierra2_4a",
			"-y", output,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			dialog.ShowError(fmt.Errorf("FFmpeg error:\n%s", stderr.String()), w)
			return
		}

		outputPath = output
		outputLabel.SetText("GIF saved as: " + filepath.Base(output))
		viewButton.Show()
	})

	w.SetContent(container.NewVBox(
		fileButton,
		fileLabel,
		widget.NewForm(
			widget.NewFormItem("Start Second", startSecEntry),
			widget.NewFormItem("End Second", endSecEntry),
		),
		runButton,
		viewButton,
		outputLabel,
	))

	w.ShowAndRun()
}

func getVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %v", err)
	}
	seconds, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse duration: %v", err)
	}
	return seconds, nil
}

func openWithDefaultBrowser(filePath string) {
	// Copy GIF to a safe temp location
	targetGIF := filepath.Join(os.TempDir(), "video2gif-preview.gif")
	inputBytes, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	_ = os.WriteFile(targetGIF, inputBytes, 0644)

	// Create minimal HTML wrapper
	html := fmt.Sprintf(`<html><body style="margin:0"><img src="file://%s" style="width:100%%;max-width:100vw"/></body></html>`, targetGIF)
	tmpHTML := filepath.Join(os.TempDir(), "video2gif-preview.html")
	_ = os.WriteFile(tmpHTML, []byte(html), 0644)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", tmpHTML)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", tmpHTML)
	default:
		cmd = exec.Command("xdg-open", tmpHTML)
	}
	_ = cmd.Start()
}

// package main
//
// import (
// 	"bytes"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"runtime"
// 	"strconv"
// 	"strings"
//
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/widget"
// )
//
// func main() {
// 	_ = os.Setenv("PATH", os.Getenv("PATH")+":/opt/homebrew/bin:/usr/local/bin")
//
// 	myApp := app.New()
// 	w := myApp.NewWindow("video2gif")
// 	w.Resize(fyne.NewSize(500, 320))
//
// 	var videoFile string
// 	var outputPath string
//
// 	startSecEntry := widget.NewEntry()
// 	endSecEntry := widget.NewEntry()
//
// 	fileLabel := widget.NewLabel("No file selected")
// 	outputLabel := widget.NewLabel("")
//
// 	viewButton := widget.NewButton("View GIF", func() {
// 		if outputPath != "" {
// 			openWithDefaultApp(outputPath)
// 		}
// 	})
// 	viewButton.Hide()
//
// 	fileButton := widget.NewButton("Choose Video File", func() {
// 		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
// 			if r != nil {
// 				videoFile = r.URI().Path()
// 				fileLabel.SetText(filepath.Base(videoFile))
//
// 				durationSec, err := getVideoDuration(videoFile)
// 				if err != nil {
// 					dialog.ShowError(err, w)
// 					return
// 				}
//
// 				startSecEntry.SetText("1")
// 				endSecEntry.SetText(strconv.Itoa(int(durationSec)))
// 			}
// 		}, w)
// 	})
//
// 	runButton := widget.NewButton("Generate GIF", func() {
// 		startSec, err1 := strconv.Atoi(startSecEntry.Text)
// 		endSec, err2 := strconv.Atoi(endSecEntry.Text)
//
// 		if videoFile == "" || err1 != nil || err2 != nil || endSec <= startSec {
// 			dialog.ShowError(fmt.Errorf("Check inputs. File must be selected and seconds must be valid."), w)
// 			return
// 		}
//
// 		duration := endSec - startSec
// 		startTime := fmt.Sprintf("00:00:%02d.000", startSec)
// 		output := strings.TrimSuffix(videoFile, filepath.Ext(videoFile)) + ".gif"
//
// 		cmd := exec.Command("ffmpeg",
// 			"-ss", startTime,
// 			"-t", fmt.Sprintf("%d", duration),
// 			"-i", videoFile,
// 			"-filter_complex",
// 			"[0:v]fps=10,scale=960:-1:flags=lanczos,split[x][z];[z]palettegen=stats_mode=diff:max_colors=128[p];[x][p]paletteuse=dither=sierra2_4a",
// 			"-y", output,
// 		)
//
// 		var stderr bytes.Buffer
// 		cmd.Stderr = &stderr
//
// 		err := cmd.Run()
// 		if err != nil {
// 			dialog.ShowError(fmt.Errorf("FFmpeg error: %s", stderr.String()), w)
// 			return
// 		}
//
// 		outputPath = output
// 		outputLabel.SetText("GIF saved as: " + filepath.Base(output))
// 		viewButton.Show()
// 	})
//
// 	w.SetContent(container.NewVBox(
// 		fileButton,
// 		fileLabel,
// 		widget.NewForm(
// 			widget.NewFormItem("Start Second", startSecEntry),
// 			widget.NewFormItem("End Second", endSecEntry),
// 		),
// 		runButton,
// 		viewButton,
// 		outputLabel,
// 	))
//
// 	w.ShowAndRun()
// }
//
// func getVideoDuration(filePath string) (float64, error) {
// 	cmd := exec.Command("ffprobe",
// 		"-v", "error",
// 		"-show_entries", "format=duration",
// 		"-of", "default=noprint_wrappers=1:nokey=1",
// 		filePath,
// 	)
// 	out, err := cmd.Output()
// 	if err != nil {
// 		return 0, fmt.Errorf("ffprobe failed: %v", err)
// 	}
// 	seconds, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
// 	if err != nil {
// 		return 0, fmt.Errorf("cannot parse duration: %v", err)
// 	}
// 	return seconds, nil
// }
//
// func openWithDefaultApp(filePath string) {
// 	var cmd *exec.Cmd
// 	switch runtime.GOOS {
// 	case "darwin":
// 		cmd = exec.Command("open", filePath)
// 	case "windows":
// 		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
// 	default: // Linux
// 		cmd = exec.Command("xdg-open", filePath)
// 	}
// 	_ = cmd.Start()
// }
