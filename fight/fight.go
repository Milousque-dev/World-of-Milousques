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
		fmt.Printf("\n%s - PV : %d/%d, Mana : %d/%d\n", joueur.Nom, joueur.Pdv, joueur.PdvMax, joueur.Mana, joueur.ManaMax)
		fmt.Printf("%s - PV : %d\n", ennemi.Nom, ennemi.Pv)

		fmt.Println("\nChoisis un sort :")
		for i, s := range joueur.Classe.Sorts {
			fmt.Printf("%d - %s (Dégâts : %d, Coût mana : %d)\n", i+1, s.Nom, s.Degats, s.Cout)
		}

		var choix int
		fmt.Print("Entrez le numéro du sort : ")
		fmt.Scanln(&choix)

		if choix < 1 || choix > len(joueur.Classe.Sorts) {
			fmt.Println("Choix invalide, réessaie.")
			continue
		}

		sortChoisi := joueur.Classe.Sorts[choix-1]

		if sortChoisi.Cout > joueur.Mana {
			fmt.Println("Pas assez de mana pour lancer ce sort !")
			continue
		}

		joueur.Mana -= sortChoisi.Cout

		ennemi.Pv -= sortChoisi.Degats
		if ennemi.Pv < 0 {
			ennemi.Pv = 0
		}

		fmt.Printf("Tu lances %s et infliges %d dégâts !\n", sortChoisi.Nom, sortChoisi.Degats)
		fmt.Printf("%s a maintenant %d PV\n", ennemi.Nom, ennemi.Pv)

		if ennemi.Pv <= 0 {
			fmt.Printf("%s est vaincu !\n", ennemi.Nom)
			return
		}

		joueur.Pdv -= ennemi.Attaque
		if joueur.Pdv < 0 {
			joueur.Pdv = 0
		}
		fmt.Printf("%s t'attaque et inflige %d dégâts !\n", ennemi.Nom, ennemi.Attaque)
		fmt.Printf("Tu as maintenant %d PV\n", joueur.Pdv)

		if joueur.Pdv <= 0 {
			fmt.Println("Tu as été vaincu...")
			return
		}
	}
}
