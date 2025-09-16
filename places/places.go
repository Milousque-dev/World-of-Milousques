package places

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"world_of_milousques/character"
	"world_of_milousques/fight"
)

func StartAdventure(name string, class string, joueur *character.Character) {
	fmt.Println("??? : Réveille toi, " + name + " le " + class + "...")
	fmt.Println("??? : Soit le bienvenue dans Acarnam")
	fmt.Println("??? : De grandes aventures t'attendent, mais avant il faut que les dieux s'assurent que tu sois prêt")

	dialogueLoop(name, class, joueur)
}

func dialogueLoop(name string, class string, joueur *character.Character) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + name + " (1) Qui êtes-vous ?  //  (2) Où suis-je ?  //  (3) Je suis prêt, on peut y aller !")
		fmt.Print("Choisis une option (1-3) : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 3 {
			fmt.Println("Entrée invalide, réessaie.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("Mathiouw : Je suis Mathiouw, le berger des jeunes âmes.")
			fmt.Println("Mathiouw : Je serai ton guide pour le début de ton existence")
		case 2:
			fmt.Println("Mathiouw : Un peu long à expliquer et ça le sera plus tard")
			fmt.Println("Mathiouw : Considère toi comme chez toi pour l'instant")
		case 3:
			fmt.Println("Mathiouw : Parfait, je te propose un petit cours sur la vie d'aventurier")
			fmt.Println("Mathiouw : Commençons par les bases du combat")
			CombatTuto(joueur)
			return
		}
	}
}

func CombatTuto(joueur *character.Character) {
	fmt.Println("Mathiouw : Aller " + joueur.Nom + ", montre comment un vrai " + joueur.Classe.Nom + " se débrouille face à ce chachagressif")

	ennemi := fight.Ennemi{"ChachaAgressif", 50, 20}
	fight.Fight(joueur, ennemi)
}
