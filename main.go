package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Problem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Question string `json:"question"`
}

type model struct {
	list list.Model
}

func main() {
	problems := loadProblems("problems/problems.json")
	items := make([]list.Item, len(problems))
	for i, p := range problems {
		items[i] = listItem{problem: p}
	}

	const defaultWidth = 20
	const listHeight = 10 // Define listHeight
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "Select a Problem"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)

	m := model{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func loadProblems(filePath string) []Problem {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading problems file:", err)
		os.Exit(1)
	}

	var problems []Problem
	if err := json.Unmarshal(file, &problems); err != nil {
		fmt.Println("Error unmarshalling problems:", err)
		os.Exit(1)
	}
	return problems
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedItem := m.list.SelectedItem().(listItem)
			openPythonEditor(selectedItem.problem.ID)
			runPythonSolution(selectedItem.problem.ID)
		case "q", "esc":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

type listItem struct {
	problem Problem
}

func (i listItem) Title() string       { return i.problem.Title }
func (i listItem) Description() string { return i.problem.Question }
func (i listItem) FilterValue() string { return i.problem.Title }

func openPythonEditor(problemID int) {
	termenv.NewOutput(os.Stdout).ClearScreen()
	cmd := exec.Command("pico", "solutions/user_solutions.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Error opening Python editor:", err)
	}
}

func runPythonSolution(problemID int) {
	cmd := exec.Command("python3", "solutions/user_solutions.py", strconv.Itoa(problemID))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running Python solution:", err)
	}
}
