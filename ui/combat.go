package ui

import (
	"fmt"
	"strings"
	"world_of_milousques/sorts"
)

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
