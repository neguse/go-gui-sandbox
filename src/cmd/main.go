package main

import (
	"fmt"
	"log"

	"os/exec"

	"bytes"

	"io"

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
	)

	i, err := MainWindow{
		Title:   "sandbox",
		MinSize: Size{Width: 320, Height: 240},
		Size:    Size{Width: 1024, Height: 800},
		// Layout:  Grid{Rows: 3, Columns: 3},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo: &pingButton,
						Text:     "ping 8.8.8.8",
						OnClicked: func() {
							pingButton.SetEnabled(false)
							pingTE.SetText("ping実行中...")
							go func() {
								defer pingButton.SetEnabled(true)
								output, err := exec.Command("cmd.exe", "/c", "ping 8.8.8.8").Output()
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
						Row:    0,
						Column: 0,
					},
					TextEdit{
						AssignTo: &pingTE,
						Row:      0,
						Column:   1,
						ReadOnly: true,
					},
				},
			},
		},
	}.Run()
	log.Println(i, err)

}
