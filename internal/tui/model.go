package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/diegosantochi/wid/internal/item"
	"github.com/diegosantochi/wid/internal/store"
)

type viewState int

const (
	listView viewState = iota
	createView
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("19")).
			Padding(1, 2)

	categoryStyle = lipgloss.NewStyle().
			Bold(true).
		//Foreground(lipgloss.Color("240")).
		Foreground(lipgloss.Color("45")).
		MarginTop(1)

	// Normal item
	normalTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "16", Dark: "255"})
	normalDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Selected item
	selectedTitleStyle = lipgloss.NewStyle().
				Bold(true).
		//Foreground(lipgloss.Color("#EE6FF8"))
		Foreground(lipgloss.Color("171"))
	selectedDescStyle = lipgloss.NewStyle().
		//Foreground(lipgloss.Color("#AD58B4"))
		Foreground(lipgloss.Color("177"))

	// Done item
	doneTitleStyle = lipgloss.NewStyle().
		//Foreground(lipgloss.Color("#04B575"))
		//Foreground(lipgloss.Color("34"))
		Foreground(lipgloss.Color("244"))

	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	labelStyle = lipgloss.NewStyle().Width(30)
)

var inputLabels = [3]string{
	"Category",
	"What is this about?",
	"What were you doing with it?",
}

type Model struct {
	store    *store.Store
	items    []item.Item
	cursor   int
	state    viewState
	inputs   [3]textinput.Model
	focusIdx int
	width    int
	height   int
}

func New() (Model, error) {
	s, err := store.Load()
	if err != nil {
		return Model{}, err
	}

	placeholders := [3]string{"e.g. Work", "e.g. Project Refactor", "e.g. Splitting the auth module into packages"}
	var inputs [3]textinput.Model
	for i := range inputs {
		t := textinput.New()
		t.Placeholder = placeholders[i]
		t.CharLimit = 200
		t.Width = 50
		inputs[i] = t
	}

	return Model{
		store:  s,
		items:  sortedItems(s.List()),
		inputs: inputs,
	}, nil
}

func sortedItems(items []item.Item) []item.Item {
	sorted := make([]item.Item, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Category != sorted[j].Category {
			return sorted[i].Category < sorted[j].Category
		}
		if sorted[i].Title != sorted[j].Title {
			return sorted[i].Title < sorted[j].Title
		}
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})
	return sorted
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.state {
	case listView:
		return m.updateList(msg)
	case createView:
		return m.updateCreate(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "enter":
			if len(m.items) > 0 {
				m.store.Toggle(m.items[m.cursor].ID)
				m.store.Save()
				m.items = sortedItems(m.store.List())
			}
		case "d":
			if len(m.items) > 0 {
				m.store.Delete(m.items[m.cursor].ID)
				m.store.Save()
				m.items = sortedItems(m.store.List())
				if m.cursor >= len(m.items) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "n":
			m.state = createView
			m.focusIdx = 0
			for i := range m.inputs {
				m.inputs[i].SetValue("")
				m.inputs[i].Blur()
			}
			return m, m.inputs[0].Focus()
		}
	}
	return m, nil
}

func (m Model) updateCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = listView
			return m, nil
		case "tab":
			m.focusIdx = (m.focusIdx + 1) % len(m.inputs)
			return m, m.focusInput()
		case "shift+tab":
			m.focusIdx = (m.focusIdx - 1 + len(m.inputs)) % len(m.inputs)
			return m, m.focusInput()
		case "enter":
			if m.focusIdx < len(m.inputs)-1 {
				m.focusIdx++
				return m, m.focusInput()
			}
			category := strings.TrimSpace(m.inputs[0].Value())
			title := strings.TrimSpace(m.inputs[1].Value())
			desc := strings.TrimSpace(m.inputs[2].Value())
			if category != "" && title != "" {
				m.store.Add(item.New(category, title, desc))
				m.store.Save()
				m.items = sortedItems(m.store.List())
			}
			m.state = listView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
	return m, cmd
}

func (m *Model) focusInput() tea.Cmd {
	var cmd tea.Cmd
	for i := range m.inputs {
		if i == m.focusIdx {
			cmd = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return cmd
}

func (m Model) View() string {
	var inner string
	switch m.state {
	case listView:
		inner = m.viewList()
	case createView:
		inner = m.viewCreate()
	}
	return appStyle.Render(inner)
}

func (m Model) viewList() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("What Was I Doing") + "\n")

	if len(m.items) == 0 {
		sb.WriteString("\n" + dimStyle.Render("No items yet. Press n to add one.") + "\n")
	} else {
		prevCategory := ""
		for i, it := range m.items {
			if it.Category != prevCategory {
				sb.WriteString(categoryStyle.Render(it.Category) + "\n")
				prevCategory = it.Category
			}

			selected := i == m.cursor
			done := it.Status == item.StatusDone

			date := dimStyle.Render(it.CreatedAt.Format("2006-01-02"))

			var titleLine, descLine string
			titlePrefix := "  "
			if selected {
				titlePrefix = "> "
			}
			if selected {
				titleLine = selectedTitleStyle.Render(fmt.Sprintf("%s%s", titlePrefix, it.Title)) + "  " + date
				descLine = selectedDescStyle.Render(fmt.Sprintf("  %s", it.Description))
			} else if done {
				titleLine = doneTitleStyle.Render(fmt.Sprintf("%s%s", titlePrefix, it.Title)) + "  " + date
				descLine = normalDescStyle.Render(fmt.Sprintf("  %s", it.Description))
			} else {
				titleLine = normalTitleStyle.Render(fmt.Sprintf("%s%s", titlePrefix, it.Title)) + "  " + date
				descLine = normalDescStyle.Render(fmt.Sprintf("  %s", it.Description))
			}

			sb.WriteString(titleLine + "\n")
			if it.Description != "" {
				sb.WriteString(descLine + "\n")
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n" + dimStyle.Render("↑/↓: navigate  space: toggle  d: delete  n: new  q: quit"))

	return sb.String()
}

func (m Model) viewCreate() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("New Item") + "\n\n")
	for i, input := range m.inputs {
		sb.WriteString(labelStyle.Render(inputLabels[i]) + input.View() + "\n")
	}
	sb.WriteString("\n" + dimStyle.Render("tab: next field  enter: confirm  esc: cancel"))

	return sb.String()
}
