package commandline

import (
	"fmt"
	"github.com/xhyyd/atm/internal/ui"
)

type cli struct {

}

func (c cli) Println(text string) {
	fmt.Println(text)
}

func NewUI() ui.UI {
	return &cli{}
}