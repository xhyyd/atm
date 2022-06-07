package main

import (
	"bufio"
	"fmt"
	"github.com/xhyyd/atm/internal/account/impl/local"
	"github.com/xhyyd/atm/internal/atm"
	"github.com/xhyyd/atm/internal/cashier/impl/cashier20"
	"github.com/xhyyd/atm/internal/command"
	"github.com/xhyyd/atm/internal/ui/impl/commandline"
	"os"
)

func newMachine() *atm.ATM {
	cashier := cashier20.NewCashierWithMoney(10000)
	am, err := local.NewManagerWithData([][]string{
		{"2859459814","7386","10.24"},
		{"1434597300","4557","90000.55"},
		{"7089382418","0075","0.00"},
		{"2001377812","5950","60.00"},
	})
	if err != nil {
		fmt.Println("error in initialize account manager:", err.Error())
		os.Exit(1)
	}
	ui := commandline.NewUI()
	return atm.NewATM(am, cashier, ui)
}

func run(atm *atm.ATM) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()

		err := command.Process(atm, text)
		if err != nil && err == command.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}


func main() {
	run(newMachine())
}
