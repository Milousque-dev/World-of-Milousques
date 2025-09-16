package fight

import (
	"fmt"
	"world_of_milousques/character"
)

type Ennemi struct {
	Nom     string
	Pv      int
	Attaque int
}

func Fight(joueur *character.Character, ennemi Ennemi) {
	fmt.Printf("\nUn %s apparaît !\n", ennemi.Nom)

	for joueur.Pdv > 0 && ennemi.Pv > 0 {
		fmt.Println("\nTes PVs :", joueur.Pdv, "/", joueur.Classe.Pvmax, " | Mana :", joueur.Mana, "/", joueur.Classe.ManaMax)
		fmt.Println(ennemi.Nom, "a", ennemi.Pv, "PVs restants.")
		fmt.Println("\nChoisis un sort ou utilise une potion :")

		for i, s := range joueur.Classe.Sorts {
			fmt.Printf("%d - %s (dégâts : %d, coût mana : %d)\n", i+1, s.Nom, s.Degats, s.Cout)
		}

		optionPotion := len(joueur.Classe.Sorts) + 1
		fmt.Printf("%d - Utiliser une potion (+50 PV)\n", optionPotion)

		var choix int
		fmt.Scanln(&choix)

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

	if joueur.Pdv <= 0 {
		fmt.Println("Tu as été vaincu... Game Over.")
	} else {
		joueur.Mana = joueur.Classe.ManaMax
		fmt.Println("Ta mana est maintenant restaurée au maximum !")
	}
}
