package ui

import (
	"io/fs"
	"os"
	"strings"

	"github.com/alx99/fly/internal/util"
	tea "github.com/charmbracelet/bubbletea"
)

type fileWindow struct {
	path     string
	files    []fs.DirEntry
	moveDown bool

	h, w           int
	pos            int
	prevFileStart  int
	fileStart      int
	visibleFileLen int

	// Configurable settings
	padding int
}

func NewFileWindow(path string) fileWindow {
	return fileWindow{path: path, padding: 2}
}

func (fw fileWindow) Init() tea.Cmd {
	return fw.readFiles
}

func (fw fileWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []fs.DirEntry:
		fw.files = msg
		fw.visibleFileLen = len(fw.files)

	case tea.WindowSizeMsg:
		fw.h, fw.w = msg.Height, msg.Width

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return fw, tea.Quit
		}

		switch kp := msg.String(); kp {
		case "ctrl+c", "q":
			return fw, tea.Quit

		case "e":
			if fw.pos > 0 {
				fw.pos -= 1
				fw.calcViewPort(false)
			}

		case "n":
			if fw.pos < fw.visibleFileLen-1 {
				fw.pos += 1
				fw.calcViewPort(true)
			}
		}
	}

	return fw, nil
}

func (fw *fileWindow) calcViewPort(moveDown bool) {
	if moveDown && fw.pos-fw.fileStart+1 > (fw.h-fw.padding) {
		fw.fileStart = util.Min(fw.visibleFileLen-fw.h, fw.fileStart+1)
	}

	if !moveDown && fw.fileStart > (fw.pos-fw.padding) {
		fw.fileStart = util.Max(0, fw.fileStart-1)
	}
}

func (fw fileWindow) View() string {
	var v strings.Builder
	drawn := 0

	for i := fw.fileStart; i < fw.visibleFileLen && drawn < fw.h; i++ {
		drawn++
		if i == fw.pos {
			v.WriteString("> ")
		}

		v.WriteString(fw.files[i].Name())

		if fw.files[i].IsDir() {
			v.WriteString("/")
		}

		if drawn+1 <= fw.h && i < fw.visibleFileLen-1 {
			v.WriteString("\n")
		}
	}

	return v.String()
}

func (fw fileWindow) readFiles() tea.Msg {
	files, err := os.ReadDir(fw.path)
	if err != nil {
		panic(err) // todo
	}
	return files
}
