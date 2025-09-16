package fight

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/ui"
)

type Ennemi struct {
	Nom     string
	Pv      int
	Attaque int
}

func Fight(joueur *character.Character, ennemi Ennemi) {
	for joueur.Pdv > 0 && ennemi.Pv > 0 {

		fmt.Printf("\nTes PV : %d / %d | Mana : %d / %d\n", joueur.Pdv, joueur.Classe.Pvmax, joueur.Mana, joueur.Classe.ManaMax)
		fmt.Printf("%s a %d PV\n", ennemi.Nom, ennemi.Pv)

		options := []string{}
		for _, s := range joueur.Classe.Sorts {
			options = append(options, fmt.Sprintf("%s (dégâts : %d, coût mana : %d)", s.Nom, s.Degats, s.Cout))
		}
		options = append(options, fmt.Sprintf("Utiliser une potion (+50 PV) (%d disponibles)", joueur.Inventaire.Potions))

		ui.AfficherMenu("Combat contre "+ennemi.Nom, options)
		fmt.Print("Choisis ton action : ")

		var choix int
		fmt.Scanln(&choix)
		optionPotion := len(joueur.Classe.Sorts) + 1

		if choix == optionPotion {
			if joueur.Inventaire.Potions > 0 {
				joueur.Pdv += 50
				if joueur.Pdv > joueur.Classe.Pvmax {
					joueur.Pdv = joueur.Classe.Pvmax
				}
				joueur.Inventaire.Potions--
				fmt.Println("Tu utilises une potion et récupères 50 PVs !")
			} else {
				fmt.Println("Tu n'as pas de potion !")
				continue
			}
		} else if choix >= 1 && choix <= len(joueur.Classe.Sorts) {
			s := joueur.Classe.Sorts[choix-1]
			if joueur.Mana < s.Cout {
				fmt.Println("Pas assez de mana pour lancer ce sort !")
				continue
			}
			joueur.Mana -= s.Cout
			ennemi.Pv -= s.Degats
			fmt.Printf("Tu lances %s et infliges %d dégâts !\n", s.Nom, s.Degats)
		} else {
			fmt.Println("Choix invalide, réessaie.")
			continue
		}

		if ennemi.Pv <= 0 {
			fmt.Printf("%s est vaincu !\n", ennemi.Nom)
			break
		}

		joueur.Pdv -= ennemi.Attaque
		fmt.Printf("%s t'attaque et inflige %d dégâts !\n", ennemi.Nom, ennemi.Attaque)
	}

	if joueur.Pdv > 0 {
		joueur.Mana = joueur.Classe.ManaMax
		fmt.Println("Ta mana est maintenant restaurée au maximum !")
	} else {
		fmt.Println("Tu as été vaincu... Game Over.")
	}
}
