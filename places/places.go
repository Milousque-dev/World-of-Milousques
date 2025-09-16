package places

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/fight"
)

func StartAdventure(joueur *character.Character) {
	fmt.Printf("??? : Réveille toi, %s le %s...\n", joueur.Nom, joueur.Classe.Nom)
	fmt.Println("??? : Sois le bienvenu dans Acarnam.")
	fmt.Println("??? : De grandes aventures t'attendent, mais avant il faut que les dieux s'assurent que tu sois prêt.")

	for {
		fmt.Printf("\nQue veux-tu demander ?\n1 - Qui êtes-vous ?\n2 - Où suis-je ?\n3 - Je suis prêt, on peut y aller !\n")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("Mathiouw : Je suis Mathiouw, le berger des jeunes âmes.")
			fmt.Println("Mathiouw : Je serais ton guide pour le début de ton existence.")
		case 2:
			fmt.Println("Mathiouw : Un peu long à expliquer, cela viendra plus tard.")
			fmt.Println("Mathiouw : Considère-toi comme chez toi pour l'instant.")
		case 3:
			fmt.Println("Mathiouw : Très bien, ta première quête consiste à vaincre un Chacha Agressif !")
			TutorielCombatQuete(joueur)
			return
		default:
			fmt.Println("Commande non reconnue, réessaie.")
		}
	}
}

func TutorielCombatQuete(joueur *character.Character) {
	quete := "Vaincre le Chacha Agressif"
	fmt.Printf("\nMathiouw : Quête proposée : %s\n", quete)
	fmt.Println("Veux-tu accepter cette quête ? (1 = Oui, 2 = Non)")

	var choix int
	fmt.Scanln(&choix)
	if choix == 1 {
		joueur.ProposerEtAjouterQuete(quete)
		fmt.Println("Quête acceptée !")
	} else {
		fmt.Println("Tu as refusé la quête, mais le chacha l'a prise pour toi !")
		fmt.Println("Qu'il est malin ce chacha, il empochera la récompenses si tu le bat")
	}

	ennemi := fight.Ennemi{Nom: "Chacha Agressif", Pv: 50, Attaque: 15}
	fight.Fight(joueur, ennemi)

	if ennemi.Pv <= 0 && choix == 1 {
		joueur.CompleterQuete(quete)
		fmt.Println("Quête terminée : Tu as vaincu le Chacha !")
	}

	fmt.Println("\nMathiouw : Bien joué ! Tu es maintenant prêt pour tes aventures.")
}
