package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"world_of_milousques/character"
	"world_of_milousques/classe"
	"world_of_milousques/fight"
	"world_of_milousques/places"
	"world_of_milousques/ui"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	c := gererMenuPrincipal(reader)
	if c == nil {
		return
	}

	for _, scene := range places.GetIntroDialogue() {
		fmt.Println("\n==== " + scene.Titre + " ====")
		fmt.Println(scene.Description)
		ui.AfficherMenu("Choisissez une option", scene.Options)
		var choix int
		fmt.Print("Votre choix : ")
		fmt.Scanln(&choix)
		if choix >= 1 && choix <= len(scene.Actions) {
			scene.Actions[choix-1](c)
		}
	}

	quete, recompense, ennemi := places.GetTutorielCombat()
	fmt.Printf("\nQuête proposée : %s\nRécompense : %s\n", quete, recompense)
	c.ProposerEtAjouterQuete(quete, recompense)
	fight.Fight(c, ennemi)
	if ennemi.Pv <= 0 {
		c.CompleterQuete(quete)
	}

	for _, scene := range places.GetScenesAventure() {
		fmt.Println("\n==== " + scene.Titre + " ====")
		fmt.Println(scene.Description)
		ui.AfficherMenu("Que faites-vous ?", scene.Options)
		var choix int
		fmt.Print("Votre choix : ")
		fmt.Scanln(&choix)
		if choix >= 1 && choix <= len(scene.Actions) {
			scene.Actions[choix-1](c)
		}
	}
}

func gererMenuPrincipal(reader *bufio.Reader) *character.Character {
	for {
		options := []string{"Créer un personnage", "Charger un personnage existant", "Quitter"}
		ui.AfficherMenu("World of Milousques", options)
		fmt.Print("Entrez votre choix : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		switch input {
		case "1", "CREER":
			c := creerPersonnage(reader)
			if c.Nom != "" {
				return &c
			}
		case "2", "REPRENDRE":
			c := reprendrePersonnage(reader)
			if c != nil {
				return c
			}
		case "3", "QUITTER":
			fmt.Println("Au revoir !")
			return nil
		default:
			fmt.Println("Commande non reconnue.")
		}
	}
}

func creerPersonnage(reader *bufio.Reader) character.Character {
	fmt.Print("Entrez le nom de votre personnage : ")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)

	classes := classe.GetClassesDisponibles()
	classOptions := []string{}
	for _, cl := range classes {
		classOptions = append(classOptions, fmt.Sprintf("%s (PV max : %d, Mana max : %d)", cl.Nom, cl.Pvmax, cl.ManaMax))
	}

	ui.AfficherMenu("Choisissez la classe de votre personnage", classOptions)
	fmt.Print("Entrez le numéro de la classe : ")
	var choix int
	fmt.Scanln(&choix)

	if choix < 1 || choix > len(classes) {
		fmt.Println("Classe invalide, création annulée.")
		return character.Character{}
	}

	classeChoisie := classes[choix-1]
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
	fmt.Print("Entrez le nom du personnage à charger : ")
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
