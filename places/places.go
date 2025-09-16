package places

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/fight"
	"world_of_milousques/ui"
)

func StartAdventure(joueur *character.Character) {
	fmt.Printf("??? : Réveille toi, %s le %s...\n", joueur.Nom, joueur.Classe.Nom)
	fmt.Println("??? : Sois le bienvenu dans Acarnam.")
	fmt.Println("??? : De grandes aventures t'attendent, mais avant il faut que les dieux s'assurent que tu sois prêt.")

	for {
		options := []string{
			"Qui êtes-vous ?",
			"Où suis-je ?",
			"Je suis prêt, on peut y aller !",
		}
		ui.AfficherMenu("Dialogue avec Mathiouw", options)

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("Mathiouw : Je suis Mathiouw, le berger des jeunes âmes.")
		case 2:
			fmt.Println("Mathiouw : Un peu long à expliquer, cela viendra plus tard.")
		case 3:
			fmt.Println("Mathiouw : Très bien, ta première quête consiste à vaincre un Chacha Agressif !")
			TutorielCombatQuete(joueur)
			return
		default:
			fmt.Println("Commande non reconnue.")
		}
	}
}

func TutorielCombatQuete(joueur *character.Character) {
	quete := "Vaincre le Chacha Agressif"
	recompense := "50 pièces d'or et 1 potion"
	fmt.Printf("\nMathiouw : Quête proposée : %s\n", quete)

	options := []string{
		"Accepter la quête",
		"Refuser la quête",
	}
	ui.AfficherMenu("Décision du joueur", options)

	var choix int
	fmt.Scanln(&choix)
	if choix == 1 {
		joueur.ProposerEtAjouterQuete(quete, recompense)
		fmt.Println("Quête acceptée !")
	} else {
		fmt.Println("Tu as refusé la quête, mais le chacha l'a prise pour toi !")
		fmt.Println("Qu'il est malin ce chacha, il empochera la récompense si tu le bats")
		joueur.ProposerEtAjouterQuete(quete, recompense)
	}

	ennemi := fight.Ennemi{Nom: "Chacha Agressif", Pv: 50, Attaque: 15}
	fight.Fight(joueur, ennemi)

	if ennemi.Pv <= 0 {
		joueur.CompleterQuete(quete)
	}

	fmt.Println("\nMathiouw : Bien joué ! Tu es maintenant prêt pour tes aventures.")
}
