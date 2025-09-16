package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"world_of_milousques/character"
	"world_of_milousques/classe"
	"world_of_milousques/places"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Bienvenue dans World of Milousques")
	fmt.Println("Voulez-vous créer un personnage ou reprendre un personnage existant ?")
	fmt.Println("Tapez CREER pour créer un nouveau personnage ou REPRENDRE pour sélectionner un personnage existant")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur de lecture :", err)
		return
	}
	input = strings.ToUpper(strings.TrimSpace(input))

	switch input {
	case "CREER":
		c := creerPersonnage(reader)
		places.StartAdventure(&c)
	case "REPRENDRE":
		c := reprendrePersonnage(reader)
		if c != nil {
			places.StartAdventure(c)
		}
	default:
		fmt.Println("Commande non reconnue.")
	}
	fmt.Println("")
}

func creerPersonnage(reader *bufio.Reader) character.Character {
	fmt.Println("Entrez le nom de votre personnage :")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)

	fmt.Println("Voulez-vous voir un aperçu des classes disponibles ou choisir directement une classe ?")
	fmt.Println("Tapez APERCU ou CHOISIR")
	choix, _ := reader.ReadString('\n')
	choix = strings.ToUpper(strings.TrimSpace(choix))

	var classeNom string
	switch choix {
	case "APERCU":
		classes := classe.GetClassesDisponibles()
		fmt.Println("Voici les classes disponibles :")
		for _, cl := range classes {
			fmt.Println("- " + cl.Nom + " (PV max : " + fmt.Sprint(cl.Pvmax) + ", Mana max : " + fmt.Sprint(cl.ManaMax) + ")")
		}
		fmt.Println("Maintenant, entrez la classe de votre personnage :")
		classeNom, _ = reader.ReadString('\n')
		classeNom = strings.TrimSpace(classeNom)
	case "CHOISIR":
		fmt.Println("Entrez la classe de votre personnage (Guerrier, Mage, Voleur) :")
		classeNom, _ = reader.ReadString('\n')
		classeNom = strings.TrimSpace(classeNom)
	default:
		fmt.Println("Commande non reconnue.")
		return character.Character{}
	}

	classeChoisie := classe.GetClasse(classeNom)
	c := character.InitCharacter(nom, classeChoisie, 1, classeChoisie.Pvmax, classeChoisie.Pvmax)

	fmt.Println("Personnage créé !")
	afficherPersonnage(&c)

	err := os.MkdirAll("saves", os.ModePerm)
	if err != nil {
		fmt.Println("Erreur lors de la création du dossier de sauvegarde :", err)
		return c
	}

	if err := c.Sauvegarder(); err != nil {
		fmt.Println("Erreur lors de la sauvegarde :", err)
	} else {
		fmt.Println("Personnage sauvegardé avec succès dans saves/" + c.Nom + ".json")
	}

	return c
}

func reprendrePersonnage(reader *bufio.Reader) *character.Character {
	fmt.Println("Entrez le nom du personnage à charger :")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)

	c, err := character.Charger(nom)
	if err != nil {
		fmt.Println("Erreur lors du chargement du personnage :", err)
		return nil
	}

	fmt.Println("Personnage chargé avec succès !")
	afficherPersonnage(c)
	return c
}

func afficherPersonnage(c *character.Character) {
	fmt.Println("Nom :", c.Nom)
	fmt.Println("Classe :", c.Classe.Nom)
	fmt.Println("Niveau :", c.Niveau)
	fmt.Println("PV :", c.Pdv, "/", c.Classe.Pvmax)
	fmt.Println("Mana :", c.Mana, "/", c.Classe.ManaMax)
	fmt.Println("Potions dans l'inventaire :", c.Inventaire.Potions)
	if len(c.Classe.Sorts) > 0 {
		fmt.Println("Sorts disponibles :")
		for _, s := range c.Classe.Sorts {
			fmt.Printf("- %s (Dégâts : %d, Coût en mana : %d)\n", s.Nom, s.Degats, s.Cout)
		}
	}
}
