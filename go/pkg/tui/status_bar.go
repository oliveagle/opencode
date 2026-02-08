package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusBar struct {
	view    *tview.TextView
	status  string
	line    int
	column  int
	fileLen int
}

func NewStatusBar() *StatusBar {
	tv := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	sb := &StatusBar{
		view: tv,
	}

	sb.updateDisplay()

	// Update time every second
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			sb.updateDisplay()
		}
	}()

	return sb
}

func (sb *StatusBar) SetStatus(status string) {
	sb.status = status
	sb.updateDisplay()
}

func (sb *StatusBar) SetCursor(line, column int) {
	sb.line = line
	sb.column = column
	sb.updateDisplay()
}

func (sb *StatusBar) SetFileLength(length int) {
	sb.fileLen = length
	sb.updateDisplay()
}

func (sb *StatusBar) View() *tview.TextView {
	return sb.view
}

func (sb *StatusBar) updateDisplay() {
	now := time.Now().Format("15:04:05")
	text := fmt.Sprintf(" %s | L:%d C:%d | %s | %s",
		sb.status,
		sb.line+1,
		sb.column,
		sb.fileLen,
		now,
	)

	sb.view.SetText(text)
	sb.view.SetBackgroundColor(tcell.ColorDarkSlateGray)
	sb.view.SetTextColor(tcell.ColorWhite)
}
