// World of Milousques - Jeu de rôle en ligne de commande
// Ce fichier contient le point d'entrée principal du jeu et la gestion des menus initiaux
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"world_of_milousques/character"
	"world_of_milousques/classe"
	"world_of_milousques/exploration"
	"world_of_milousques/fight"
	"world_of_milousques/places"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	c := gererMenuPrincipal()
	if c == nil {
		return
	}

	c.InitialiserEtatMap()

	// Gérer l'introduction/tutoriel ou reprise d'aventure
	if !executerIntroductionOuReprise(c) {
		return // Le joueur a été vaincu pendant le tutoriel
	}
	
	// Sauvegarder le personnage avant de commencer/continuer l'exploration
	sauvegarderPersonnageAvecMessage(c, "avant de commencer l'exploration")
	
	// Lancer le système d'exploration
	exploration.ExplorerMap(c)
}

// executerIntroductionOuReprise gère l'introduction pour un nouveau joueur ou la reprise d'aventure
// Retourne true si le jeu peut continuer, false si le joueur a été vaincu
func executerIntroductionOuReprise(c *character.Character) bool {
	if !c.AIntroEffectuee() {
		return executerTutoriel(c)
	} else {
		reprendreAventure(c)
		return true
	}
}

// executerTutoriel lance le tutoriel d'introduction
func executerTutoriel(c *character.Character) bool {
	// Gestion du dialogue d'introduction avec boucle
	scenes := places.GetIntroDialogue()
	if len(scenes) > 0 {
		scene := scenes[0] // Prendre la première (et unique) scène
		
		fmt.Println("\n==== " + scene.Titre + " ====")
		fmt.Println(scene.Description)
		
		// Boucle de dialogue jusqu'à ce que le joueur choisisse "On peut y aller !"
		for {
			ui.AfficherMenu("Choisissez une option", scene.Options)
			choix := utils.ScanChoice("Votre choix : ", scene.Options)
			
			// Exécuter l'action correspondante
			scene.Actions[choix-1](c)
			
			// Si le joueur choisit "On peut y aller !" (option 4), sortir de la boucle
			if choix == 4 {
				break
			}
			
			// Pause pour que le joueur puisse lire la réponse
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		}
	}

	// Proposer la quête du tutoriel
	queteAcceptee := places.ProposerQueteTutoriel(c)
	
	quete, _, ennemi := places.GetTutorielCombat()
	fight.Fight(c, ennemi)
	
	// Vérifier si le joueur a survécu
	if c.Pdv > 0 {
		// Compléter la quête si elle avait été acceptée
		if queteAcceptee {
			c.CompleterQuete(quete)
			fmt.Println("\n✅ Quête accomplie automatiquement !")
			c.RendreQuete(quete) // Rendre automatiquement la quête du tutoriel
		} else {
			fmt.Println("\n💰 Le chacha laisse tomber 1 potion en mourant !")
			c.Inventaire.Potions++
		}
		
		// Marquer l'introduction comme effectuée
		c.MarquerIntroEffectuee()
		
		fmt.Println("\n🎉 Félicitations ! Vous avez terminé le tutoriel !")
		fmt.Println("Le monde s'ouvre maintenant à vous...")
		fmt.Println("\nAppuyez sur Entrée pour commencer votre aventure...")
		fmt.Scanln()
		return true
	} else {
		fmt.Println("\n💀 Vous avez été vaincu pendant le tutoriel...")
		fmt.Println("Le jeu se termine ici. Réessayez !")
		return false
	}
}

// reprendreAventure affiche le message de reprise pour un joueur existant
func reprendreAventure(c *character.Character) {
	fmt.Println("\n🔄 Reprise de votre aventure...")
	x, y := c.ObtenirPosition()
	fmt.Printf("Vous êtes à la position (%d, %d) sur la carte.\n", x+1, y+1)
}

// gererMenuPrincipal affiche le menu principal et gère les choix de l'utilisateur
// Retourne un personnage prêt pour l'aventure, ou nil si l'utilisateur souhaite quitter
func gererMenuPrincipal() *character.Character {
	for {
		afficherMenuPrincipalJeu()
		choix := demanderChoixMenuPrincipal()
		
		personnage := traiterChoixMenuPrincipal(choix)
		if personnage != nil || choix == 3 {
			return personnage // Retourne le personnage ou nil (pour quitter)
		}
		// Si personnage est nil et choix != 3, continuer la boucle
	}
}

// centrerTexteAvecLargeur centre un texte dans une largeur spécifique
func centrerTexteAvecLargeur(texte string, largeur int) string {
	if len(texte) >= largeur {
		return texte
	}
	padding := (largeur - len(texte)) / 2
	return strings.Repeat(" ", padding) + texte
}

// afficherMenuPrincipalJeu affiche le titre et les options du menu principal
func afficherMenuPrincipalJeu() {
	titre := "WORLD OF MILOUSQUES"
	soustitre := "Une aventure pleines de Milousqueries !"
	
	// Calculer la largeur du menu comme le fait AfficherMenuSimple
	options := []string{
		"Créer un nouveau personnage",
		"Charger un personnage existant", 
		"Quitter le jeu",
	}
	
	// Simuler le calcul de AfficherMenuSimple
	largeurContenu := len("MENU PRINCIPAL")
	for i, opt := range options {
		ligne := fmt.Sprintf("%d) %s", i+1, opt)
		if len(ligne) > largeurContenu {
			largeurContenu = len(ligne)
		}
	}
	if largeurContenu < 30 {
		largeurContenu = 30
	}
	// Largeur maximale pour éviter le wrapping dans les terminaux Windows (plus conservative)
	if largeurContenu > 50 {
		largeurContenu = 50
	}
	largeurMenu := largeurContenu + 4 + 2 // +4 pour marge interne, +2 pour les bordures du menu
	
	fmt.Println("\n" + strings.Repeat("=", largeurMenu))
	fmt.Println(centrerTexteAvecLargeur(titre, largeurMenu))
	fmt.Println(centrerTexteAvecLargeur(soustitre, largeurMenu))
	fmt.Println(strings.Repeat("=", largeurMenu))
	
	ui.AfficherMenuSimple("MENU PRINCIPAL", options)
}

// demanderChoixMenuPrincipal demande et retourne le choix de l'utilisateur
func demanderChoixMenuPrincipal() int {
	options := []string{"Créer un personnage", "Charger un personnage existant", "Quitter"}
	return utils.ScanChoice("Entrez votre choix : ", options)
}

// traiterChoixMenuPrincipal traite le choix de l'utilisateur et exécute l'action correspondante
func traiterChoixMenuPrincipal(choix int) *character.Character {
	switch choix {
	case 1:
		return gererCreationPersonnage()
	case 2:
		return reprendrePersonnage()
	case 3:
		fmt.Println("👋 Au revoir et à bientôt dans World of Milousques !")
		return nil
	default:
		fmt.Println("❌ Choix invalide, veuillez réessayer.")
		return nil
	}
}

// gererCreationPersonnage gère la création d'un personnage avec validation
func gererCreationPersonnage() *character.Character {
	c := creerPersonnage()
	if c.Nom != "" {
		return &c
	}
	return nil // Création échouée
}

func creerPersonnage() character.Character {
	nom := utils.ScanString("Entrez le nom de votre personnage : ", 1)

	classes := classe.GetClassesDisponibles()
	classOptions := make([]string, len(classes))
	for i, cl := range classes {
		classOptions[i] = fmt.Sprintf("%s (PV max : %d, Mana max : %d)", cl.Nom, cl.Pvmax, cl.ManaMax)
	}

	ui.AfficherMenu("Choisissez la classe de votre personnage", classOptions)
	choix := utils.ScanChoice("Entrez le numéro de la classe : ", classOptions)

	classeChoisie := classes[choix-1]
	c := character.InitCharacter(nom, classeChoisie, 1, classeChoisie.Pvmax, classeChoisie.Pvmax)

	fmt.Println("Personnage créé !")
	afficherPersonnage(&c)

	// Créer le dossier de sauvegarde et sauvegarder
	gererSauvegardePremiereFois(&c)

	return c
}

func reprendrePersonnage() *character.Character {
	afficherSauvegardesDisponibles()
	
	nom := utils.ScanString("Entrez le nom du personnage à charger : ", 1)
	c := chargerPersonnageAvecMessage(nom)
	if c == nil {
		return nil
	}

	afficherPersonnage(c)
	return c
}

// afficherPersonnageComplet affiche toutes les informations détaillées d'un personnage
func afficherPersonnageComplet(c *character.Character) {
	afficherPersonnageResume(c)
	
	// Équipement détaillé
	afficherEquipementDetaille(c)
	
	// Sorts disponibles
	if len(c.Classe.Sorts) > 0 {
		fmt.Println("\nSorts disponibles :")
		for _, s := range c.Classe.Sorts {
			fmt.Printf("- %s (Dégâts : %d, Coût en mana : %d)\n", s.Nom, s.Degats, s.Cout)
		}
	}
}

// afficherPersonnageResume affiche un résumé compact d'un personnage
func afficherPersonnageResume(c *character.Character) {
	fmt.Printf("⚔️  %s (%s niveau %d)\n", c.Nom, c.Classe.Nom, c.Niveau)
	fmt.Printf("   PV: %d/%d | Mana: %d/%d | XP: %d/%d\n", 
		c.Pdv, c.PdvMax, c.Mana, c.ManaMax, c.Experience, c.CalculerXPRequis())
	fmt.Printf("   💰 %d or | 🧆 %d potions | 🗺️ %d/25 zones\n", 
		c.Argent, c.Inventaire.Potions, c.ObtenirNombreZonesDecouvertes())
}

// afficherEquipementDetaille affiche l'équipement d'un personnage
func afficherEquipementDetaille(c *character.Character) {
	equipements := []string{}
	if c.ArmeEquipee != nil {
		equipements = append(equipements, "⚔️  "+c.ArmeEquipee.Nom)
	}
	if c.CasqueEquipe != nil {
		equipements = append(equipements, "🪖 "+c.CasqueEquipe.Nom)
	}
	if c.TorseEquipe != nil {
		equipements = append(equipements, "👕 "+c.TorseEquipe.Nom)
	}
	if c.JambiereEquipee != nil {
		equipements = append(equipements, "👖 "+c.JambiereEquipee.Nom)
	}
	
	if len(equipements) > 0 {
		fmt.Println("\nÉquipement :", strings.Join(equipements, ", "))
		bonusAttaque := c.CalculerAttaqueBonus()
		bonusDefense := c.CalculerDefenseBonus()
		if bonusAttaque > 0 || bonusDefense > 0 {
			fmt.Printf("Bonus : +%d Attaque, +%d Défense\n", bonusAttaque, bonusDefense)
		}
	}
}

// afficherPersonnage maintient la compatibilité - alias pour afficherPersonnageComplet
func afficherPersonnage(c *character.Character) {
	afficherPersonnageComplet(c)
}

// afficherSauvegardesDisponibles affiche un aperçu des personnages sauvegardés
func afficherSauvegardesDisponibles() {
	fmt.Println("\n💾 === SAUVEGARDES DISPONIBLES === 💾")
	
	// Lire le dossier saves
	files, err := os.ReadDir("saves")
	if err != nil {
		fmt.Println("Aucune sauvegarde trouvée.")
		return
	}
	
	aucuneSauvegarde := true
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			aucuneSauvegarde = false
			nomPersonnage := strings.TrimSuffix(file.Name(), ".json")
			
			// Charger temporairement pour afficher les infos
			c, err := character.Charger(nomPersonnage)
			if err == nil {
				// Utiliser la fonction réutilisable pour l'affichage de base
				fmt.Println()
				afficherPersonnageResume(c)
				
				// Ajouter les informations spécifiques aux sauvegardes
				afficherInfosSauvegarde(c)
			}
		}
	}
	
	if aucuneSauvegarde {
		fmt.Println("Aucune sauvegarde trouvée.")
	}
	
	fmt.Println()
}

// afficherInfosSauvegarde affiche les informations spécifiques aux sauvegardes
func afficherInfosSauvegarde(c *character.Character) {
	// Compter les quêtes actives
	totalQuetes := 0
	for _, q := range c.Quetes {
		if !q.Rendue {
			totalQuetes++
		}
	}
	
	// Afficher quêtes et équipement de façon concise
	fmt.Printf("   🎒 %d quêtes actives", totalQuetes)
	
	equipementCount := compterEquipement(c)
	if equipementCount > 0 {
		fmt.Printf(" | 🛡️  %d équipements", equipementCount)
		if c.ArmeEquipee != nil {
			fmt.Printf(" (Arme: %s)", c.ArmeEquipee.Nom)
		}
	}
	fmt.Println()
}

// compterEquipement compte le nombre de pièces d'équipement
func compterEquipement(c *character.Character) int {
	count := 0
	if c.ArmeEquipee != nil {
		count++
	}
	if c.CasqueEquipe != nil {
		count++
	}
	if c.TorseEquipe != nil {
		count++
	}
	if c.JambiereEquipee != nil {
		count++
	}
	return count
}

// === FONCTIONS UTILITAIRES POUR LA GESTION D'ERREUR ===

// gererSauvegardePremiereFois gère la création du dossier de sauvegarde et la première sauvegarde
func gererSauvegardePremiereFois(c *character.Character) {
	if err := creerDossierSauvegarde(); err != nil {
		fmt.Println("⚠️  Erreur lors de la création du dossier de sauvegarde :", err)
		return
	}
	
	sauvegarderPersonnageAvecMessage(c, "lors de la création du personnage")
}

// creerDossierSauvegarde crée le dossier de sauvegarde s'il n'existe pas
func creerDossierSauvegarde() error {
	return os.MkdirAll("saves", os.ModePerm)
}

// sauvegarderPersonnageAvecMessage sauvegarde un personnage avec un message contextuel en cas d'erreur
func sauvegarderPersonnageAvecMessage(c *character.Character, contexte string) {
	if err := c.Sauvegarder(); err != nil {
		fmt.Printf("⚠️  Erreur lors de la sauvegarde %s : %v\n", contexte, err)
	} else {
		fmt.Printf("✅ Personnage sauvegardé avec succès\n")
	}
}

// chargerPersonnageAvecMessage charge un personnage avec une gestion d'erreur explicite
func chargerPersonnageAvecMessage(nom string) *character.Character {
	c, err := character.Charger(nom)
	if err != nil {
		fmt.Printf("❌ Erreur lors du chargement du personnage '%s' : %v\n", nom, err)
		fmt.Println("Vérifiez que le nom est correct et que la sauvegarde existe.")
		return nil
	}
	
	fmt.Printf("✅ Personnage '%s' chargé avec succès !\n", nom)
	return c
}
