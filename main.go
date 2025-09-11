package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"world_of_milousques/character"
	"world_of_milousques/classe"
)

func main() {
	fmt.Println("Bienvenue dans World of Milousques")
	fmt.Println("Voulez-vous créer un personnage ou reprendre un personnage existant ?")
	fmt.Println("Tapez CREER pour créer un nouveau personnage ou REPRENDRE pour sélectionner un personnage existant")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur de lecture :", err)
		return
	}
	input = strings.TrimSpace(input)

	if input == "CREER" {
		fmt.Println("Entrez le nom de votre personnage :")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		nom := input

		fmt.Println("Voulez-vous voir un aperçu des classes disponibles ou choisir directement une classe ?")
		fmt.Println("Tapez APERCU ou CHOISIR")

		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "APERCU" {
			classes := classe.GetClassesDisponibles()
			fmt.Println("Voici les classes disponibles :")
			for _, cl := range classes {
				fmt.Println("- " + cl.Nom + " (PV max : " + fmt.Sprint(cl.Pvmax) + ")")
			}
			fmt.Println("Maintenant, entrez la classe de votre personnage :")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
		} else if input == "CHOISIR" {
			fmt.Println("Entrez la classe de votre personnage (Guerrier, Mage, Voleur) :")
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
		} else {
			fmt.Println("Commande non reconnue.")
			return
		}

		classeChoisie := classe.GetClasse(input)
		c := character.InitCharacter(nom, classeChoisie, 1, classeChoisie.Pvmax, classeChoisie.Pvmax)

		fmt.Println("Personnage créé !")
		fmt.Println("Nom :", c.Nom)
		fmt.Println("Classe :", c.Classe.Nom)
		fmt.Println("Niveau :", c.Niveau)
		fmt.Println("PV :", c.Pdv, "/", c.PdvMax)
		fmt.Println("Potions dans l'inventaire :", c.Inventaire.Potions)

		err := os.MkdirAll("saves", os.ModePerm)
		if err != nil {
			fmt.Println("Erreur lors de la création du dossier de sauvegarde :", err)
			return
		}

		saveErr := c.Sauvegarder()
		if saveErr != nil {
			fmt.Println("Erreur lors de la sauvegarde :", saveErr)
		} else {
			fmt.Println("Personnage sauvegardé avec succès dans saves/" + c.Nom + ".json")
		}

	} else if input == "REPRENDRE" {
		fmt.Println("Entrez le nom du personnage à charger :")
		nom, _ := reader.ReadString('\n')
		nom = strings.TrimSpace(nom)

		c, err := character.Charger(nom)
		if err != nil {
			fmt.Println("Erreur lors du chargement du personnage :", err)
			return
		}

		fmt.Println("Personnage chargé avec succès !")
		fmt.Println("Nom :", c.Nom)
		fmt.Println("Classe :", c.Classe.Nom)
		fmt.Println("Niveau :", c.Niveau)
		fmt.Println("PV :", c.Pdv, "/", c.PdvMax)
		fmt.Println("Potions dans l'inventaire :", c.Inventaire.Potions)

	} else {
		fmt.Println("Commande non reconnue.")
	}
}
