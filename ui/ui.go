package ui

import (
	"fmt"
	"strings"
	"world_of_milousques/sorts"
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

func AfficherMenuCombat(joueurNom string, joueurPv, joueurPvMax, joueurMana, joueurManaMax int,
	ennemiNom string, ennemiPv int, sortsList []sorts.Sorts, potions int) {

	lignes := []string{}
	lignes = append(lignes, fmt.Sprintf("%s : PV %d|%d | Mana %d|%d", joueurNom, joueurPv, joueurPvMax, joueurMana, joueurManaMax))
	lignes = append(lignes, fmt.Sprintf("Ennemi %s : PV %d", ennemiNom, ennemiPv))
	lignes = append(lignes, strings.Repeat("-", 10))

	for i, s := range sortsList {
		lignes = append(lignes, fmt.Sprintf("%d) %s  Dégâts: %d Mana: %d", i+1, s.Nom, s.Degats, s.Cout))
	}

	lignes = append(lignes, fmt.Sprintf("%d) Utiliser une potion  (+50 PV) (%d disponibles)", len(sortsList)+1, potions))

	largeur := 0
	for _, ligne := range lignes {
		if len(ligne) > largeur {
			largeur = len(ligne)
		}
	}

	fmt.Print("+")
	fmt.Print(strings.Repeat("-", largeur+2))
	fmt.Println("+")
	for _, ligne := range lignes {
		fmt.Printf("| %-*s |\n", largeur+2, ligne)
	}
	fmt.Print("+")
	fmt.Print(strings.Repeat("-", largeur+2))
	fmt.Println("+")
}
