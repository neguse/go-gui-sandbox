package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// あとこのへん調べたい
// Canvas
// CheckBox
// ComboBox
// Dialog
// FileDialog
// Form
// Label
// GridLayout
// LineEdit
// ListBox
// Menu
// PushButton
// RadioButton
// StatusBar
// TableView

// Must treats unwanted error.
// If error has occurred, Must will panic.
func Must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// OutputConverter returns CP932 to UTF8 reader, and writer
func OutputConverter() (io.Reader, io.WriteCloser) {
	r, w := io.Pipe()
	jr := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	return jr, w
}

func main() {

	var (
		pingStatusTE               *walk.LineEdit
		pingStdoutTE, pingStderrTE *walk.TextEdit
		pingButton                 *walk.PushButton
		pingCB                     *walk.ComboBox
	)

	setEnabled := func(enabled bool) {
		pingButton.SetEnabled(enabled)
		pingCB.SetEnabled(enabled)
	}

	wnd := MainWindow{
		Title:   "sandbox",
		MinSize: Size{Width: 320, Height: 240},
		Size:    Size{Width: 1024, Height: 800},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text:   "address",
						Row:    0,
						Column: 0,
					},
					ComboBox{
						AssignTo: &pingCB,
						Row:      0,
						Column:   1,
						Model: []string{
							"8.8.8.8", "1.1.1.1", "127.0.0.1",
						},
					},

					Label{
						Text: "command",
						Row:  1, Column: 0,
					},
					PushButton{
						AssignTo: &pingButton,
						Text:     "ping",
						OnClicked: func() {
							setEnabled(false)
							command := fmt.Sprint("ping ", pingCB.Text())
							Must(pingStatusTE.SetText(fmt.Sprint("executing ", command)))
							Must(pingStdoutTE.SetText(""))
							Must(pingStderrTE.SetText(""))
							go func() {
								defer setEnabled(true)
								cmd := exec.Command("cmd.exe", "/c", command)
								stdoutR, stdoutW := OutputConverter()
								stderrR, stderrW := OutputConverter()
								cmd.Stdout = stdoutW
								cmd.Stderr = stderrW
								err := cmd.Start()
								if err != nil {
									Must(pingStatusTE.SetText(fmt.Sprint(err)))
									return
								}
								go func() {
									s := bufio.NewScanner(stdoutR)
									for s.Scan() {
										log.Println(s.Text())
										pingStdoutTE.AppendText(s.Text() + "\r\n")
									}
								}()
								go func() {
									s := bufio.NewScanner(stderrR)
									for s.Scan() {
										log.Println(s.Text())
										pingStderrTE.AppendText(s.Text() + "\r\n")
									}
								}()
								err = cmd.Wait()
								if err != nil {
									Must(pingStatusTE.SetText(fmt.Sprint(err)))
								}
								Must(pingStatusTE.SetText("finished"))
							}()
						},
						Row:    1,
						Column: 1,
					},
				},
			},
			LineEdit{
				AssignTo: &pingStatusTE,
				ReadOnly: true,
			},
			TextEdit{
				AssignTo: &pingStdoutTE,
				ReadOnly: true,
			},
			TextEdit{
				AssignTo: &pingStderrTE,
				ReadOnly: true,
			},
		},
	}
	i, err := wnd.Run()
	log.Println(i, err)

}
