package ui

import (
	"fmt"
	"strings"
)

func AfficherMenu(titre string, options []string) {
	longueurMax := len(titre)
	for i, opt := range options {
		ligne := fmt.Sprintf("%d) %s", i+1, opt)
		if len(ligne) > longueurMax {
			longueurMax = len(ligne)
		}
	}
	longueurMax += 4

	fmt.Print("+")
	fmt.Print(strings.Repeat("-", longueurMax))
	fmt.Println("+")

	titrePadded := fmt.Sprintf("  %-*s", longueurMax-2, titre)
	fmt.Printf("|%s|\n", titrePadded)

	fmt.Print("|")
	fmt.Print(strings.Repeat("-", longueurMax))
	fmt.Println("|")

	for i, opt := range options {
		optionPadded := fmt.Sprintf(" %d) %-*s", i+1, longueurMax-4, opt)
		fmt.Printf("|%s|\n", optionPadded)
	}

	fmt.Print("+")
	fmt.Print(strings.Repeat("-", longueurMax))
	fmt.Println("+")
}
