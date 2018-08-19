package main

import (
	"bytes"
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

func sjisToUtf8(sjis []byte) ([]byte, error) {
	rd := transform.NewReader(bytes.NewReader(sjis), japanese.ShiftJIS.NewDecoder())
	wr := bytes.NewBuffer([]byte{})
	_, err := io.Copy(wr, rd)
	if err != nil {
		return nil, err
	}
	return wr.Bytes(), nil
}

func main() {

	var (
		pingTE     *walk.TextEdit
		pingButton *walk.PushButton
		pingCB     *walk.ComboBox
	)

	setEnabled := func(enabled bool) {
		pingTE.SetEnabled(enabled)
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
							pingTE.SetText(fmt.Sprint("executing ", command))
							go func() {
								defer setEnabled(true)
								output, err := exec.Command("cmd.exe", "/c", command).Output()
								if err != nil {
									pingTE.SetText(fmt.Sprint(err))
									return
								}
								u8output, err := sjisToUtf8(output)
								if err != nil {
									pingTE.SetText(fmt.Sprint(err))
									return
								}
								pingTE.SetText(string(u8output))
							}()
						},
						Row:    1,
						Column: 1,
					},
				},
			},
			TextEdit{
				AssignTo: &pingTE,
				Row:      0,
				Column:   1,
				ReadOnly: true,
			},
		},
	}
	i, err := wnd.Run()
	log.Println(i, err)

}
