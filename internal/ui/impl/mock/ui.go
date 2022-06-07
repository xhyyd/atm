package mock

import "github.com/xhyyd/atm/internal/ui"

type MockUI struct {
	Lines []string
}

var v = ui.UI((*MockUI)(nil))
func NewUI() *MockUI {
	return &MockUI{}
}

func (m *MockUI) Println(text string) {
	m.Lines = append(m.Lines, text)
}

func (m *MockUI)Clear() {
	m.Lines = nil
}