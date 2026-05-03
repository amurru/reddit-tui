package modal

import (
	"reddittui/components/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImagePreviewModal struct {
	Content string
	w       int
	h       int
}

func NewImagePreviewModal() ImagePreviewModal {
	return ImagePreviewModal{}
}

func (m ImagePreviewModal) View() string {
	title := lipgloss.NewStyle().Bold(true).Render("Image preview")
	help := lipgloss.NewStyle().Faint(true).Render("Press esc to close")

	contentStyle := lipgloss.NewStyle().Width(m.w).MaxHeight(m.h)
	content := contentStyle.Render(m.Content)

	container := lipgloss.JoinVertical(lipgloss.Left, title, help, "", content)
	return lipgloss.NewStyle().Width(m.w).Render(container)
}

func (m ImagePreviewModal) Update(msg tea.Msg) (ImagePreviewModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter", "h", "backspace":
			return m, messages.ExitModal
		}
	}

	return m, nil
}

func (m *ImagePreviewModal) SetSize(w, h int) {
	// leave room for title/help rows in the modal frame
	m.w = max(1, w-1)
	m.h = max(1, h-3)
}
