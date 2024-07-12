package tui

import (
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
)

type model struct {
	focus  focusState
	ready  bool
	panel  panel
	lists  [3]list.Model
	width  int
	height int
	client *jira.Client
}

func NewTuiModel(client *jira.Client) *model {
	m := &model{
		ready:  false,
		client: client,
	}

	return m
}

func (m *model) Init() tea.Cmd {
	defaultList := list.NewDefaultDelegate()
	defaultList.ShowDescription = false

	wg := sync.WaitGroup{}
	wg.Add(3)

	go m.createAssignedIssuesList(defaultList, &wg)
	go m.createProjectsList(defaultList, &wg)
	go m.createEpicsList(defaultList, &wg)

	wg.Wait()

	m.panel.tabs = []string{"Tab 1", "Tab 2", "Tab 3"}
	m.panel.tabContent = []string{
		"Content for Tab 1",
		"Content for Tab 2",
		"Content for Tab 3",
	}

	loadIssueDescriptionTab(m)

	m.ready = true

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % 3
		case "shift+tab":
			m.focus = (m.focus + 2) % 3
		// TODO: this will need redoing in the future but for demo purposes works atm
		case "]":
			m.panel.activeTab = (m.panel.activeTab + 1) % len(m.panel.tabs)
			return m, nil
		case "[":
			if m.panel.activeTab > 0 {
				m.panel.activeTab--
			} else {
				m.panel.activeTab = len(m.panel.tabs) - 1
			}
			return m, nil
		}

		if m.focus >= IssuesList && m.focus <= EpicsList {
			listIndex := int(m.focus) - int(IssuesList)
			if listIndex >= 0 && listIndex < len(m.lists) {
				m.lists[listIndex], cmd = m.lists[listIndex].Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		thirdWidth := m.width / 3
		for i := range m.lists {
			m.lists[i].SetWidth(thirdWidth)
			m.lists[i].SetHeight((m.height / 3) - 2)
		}
	}

	m.panel.tabs = []string{"Issues", "Tab 2", "Tab 3"}
	if m.focus == IssuesList {
		m.panel.tabs = []string{"Description", "Tab 2", "Tab 3"}
	}

	loadIssueDescriptionTab(m)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var views []string
	for i, list := range m.lists {
		listStyle := unfocusedStyle
		if m.focus == focusState(IssuesList+focusState(i)) {
			listStyle = focusedStyle
		}
		views = append(views, listStyle.Render(list.View()))
	}

	listView := lipgloss.JoinVertical(lipgloss.Top, views...)
	panelView := m.panel.View()
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, listView, panelView)
	return mainView
}

func loadIssueDescriptionTab(m *model) {
	if m.focus == 0 {
		selectedItem := m.lists[m.focus].SelectedItem().(jira.Issue)
		converter := md.NewConverter("", true, nil)
		markdown, err := converter.ConvertString(selectedItem.GetRenderedDescription())

		if err != nil {
			markdown = ""
		}

		m.panel.tabContent[0] = markdown
	}
}
