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

	input, _ := reader.ReadString('\n')
	input = strings.ToUpper(strings.TrimSpace(input))

	switch input {
	case "CREER":
		creerPersonnage(reader)
	case "REPRENDRE":
		reprendrePersonnage(reader)
	default:
		fmt.Println("Commande non reconnue.")
	}
	fmt.Println("")
}

func creerPersonnage(reader *bufio.Reader) {
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
		return
	}

	classeChoisie := classe.GetClasse(classeNom)
	c := character.InitCharacter(
		nom,
		classeChoisie,
		1,
		classeChoisie.Pvmax,
		classeChoisie.Pvmax,
		classeChoisie.ManaMax,
		classeChoisie.ManaMax,
	)

	fmt.Println("Personnage créé !")
	afficherPersonnage(&c)

	_ = os.MkdirAll("saves", os.ModePerm)

	if err := c.Sauvegarder(); err != nil {
		fmt.Println("Erreur lors de la sauvegarde :", err)
	} else {
		fmt.Println("Personnage sauvegardé avec succès dans saves/" + c.Nom + ".json")
	}

	places.StartAdventure(c.Nom, c.Classe.Nom, &c)
}

func reprendrePersonnage(reader *bufio.Reader) {
	fmt.Println("Entrez le nom du personnage à charger :")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)

	c, err := character.Charger(nom)
	if err != nil {
		fmt.Println("Erreur lors du chargement du personnage :", err)
		return
	}

	fmt.Println("Personnage chargé avec succès !")
	afficherPersonnage(c)

	places.StartAdventure(c.Nom, c.Classe.Nom, c)
}

func afficherPersonnage(c *character.Character) {
	fmt.Println("Nom :", c.Nom)
	fmt.Println("Classe :", c.Classe.Nom)
	fmt.Println("Niveau :", c.Niveau)
	fmt.Println("PV :", c.Pdv, "/", c.PdvMax)
	fmt.Println("Mana :", c.Mana, "/", c.ManaMax)
	fmt.Println("Potions dans l'inventaire :", c.Inventaire.Potions)
}
