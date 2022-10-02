package mainview

import (
	"github.com/alx99/fly/internal/ui/views/fileview"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	pd = iota // parent directory
	wd        // working directoy
	cd        // child directory
)

type mainView struct {
	fws  []fileview.Window
	h, w int
}

func New() mainView {
	fws := make([]fileview.Window, 3)
	fws[0] = fileview.New("/", 0, 0)
	fws[1] = fileview.New("/var", 0, 0)
	fws[2] = fileview.New("/var/lib", 0, 0)

	return mainView{fws: fws}
}

func (mw mainView) Init() tea.Cmd {
	return tea.Batch(mw.fws[0].Init, mw.fws[1].Init, mw.fws[2].Init)
}

func (mw mainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mw.h, mw.w = msg.Height, msg.Width

	case tea.KeyMsg:
		switch kp := msg.String(); kp {
		case "ctrl+c", "q":
			return mw, tea.Quit

		case "e":
			mw.fws[wd].Move(fileview.Up)
			if mw.fws[wd].GetSelection().IsDir() {
				mw.fws[cd] = fileview.New(mw.fws[1].GetSelectedPath(), mw.w, mw.h)
				return mw, mw.fws[2].Init
			}
			return mw, nil

		case "n":
			mw.fws[wd].Move(fileview.Down)
			if mw.fws[wd].GetSelection().IsDir() {
				mw.fws[cd] = fileview.New(mw.fws[1].GetSelectedPath(), mw.w, mw.h)
				mw.fws[cd].Update(tea.WindowSizeMsg{Height: mw.h, Width: mw.w})
				return mw, mw.fws[2].Init
			}
			return mw, nil
		}
	}

	for i := range mw.fws {
		mw.fws[i], _ = mw.fws[i].Update(msg)
	}

	return mw, nil
}

func (mw mainView) View() string {
	res := make([]string, len(mw.fws))
	for _, fw := range mw.fws {
		res = append(res, fw.View())
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, res...)
}
