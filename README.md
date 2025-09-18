## World of Milousques - Guide Technique Détaillé

## Sommaire

1. Introduction et présentation du projet
2. Structure du Projet
3. Explication Détaillée de Chaque Fichier

## 1. Introduction et présentation du projet :

World of Milousques est un RPG textuel solo d'exploration
Dans la version actuelle du jeu, il n'y a pas vraiment d'objectif définit, simplement du contenu comme des quêtes a réaliser, des objets a crafter, des ressources a récolter et des monstres.
Dans de prochaine mise a jour on verra une histoire et un but se dessinera dans les zones suivantes
Pour l'instant la trame du jeu se découpe en deux parties : l'introduction, puis une zone explorable qui fera office de tutoriel lors des versions finales du jeu.
Le jeu a un système de sauvegarde par json très pousser permettant de garde en mémoire la moindre actions du sauvegarde que ce soit l'exploration des zones, les monstres que l'on a vaincu sont sauvegardé comme vaincu, de même pour les ressources que l'on récolte et bien évidement le systeme de sauvegarde comprend ausssi les fonctionalités plus classique comme la sauvegarde de la progression dans les quetes ou de l'inventaire / équipement.
Il y a un système de déplacement sur une map pleine de vie, aux différents biomes, ressources, adverssaires et pnj a croiser.
Tout les menus sont en ASCII Art avec des caractère unicode, on utilise une fonction qui s'adapte et créer le menu en fonction du texte a mettre dedans.


## 2. Structure du projet


World-of-Milousques/
    main.go                    // Point d'entrée du programme
    go.mod                     // Fichier de configuration du projet Go
    saves/                     // Dossier des sauvegardes de jeu
        Nomdupersonnage.json   // Le fichier est automatiquement créer a la création du personnage
    banque/                    // Système de stockage via une banque
        banque.go
    character/                 // Gestion du personnage, de sa création et de la sauvegarde
        character.go
    classe/                    // Système de classe
        classe.go
    commerce/                  // 
        commerce.go
    craft/                     // Système de fabrication
        craft.go
    exploration/               // Exploration du monde
        exploration.go
     fight/                    // Système de combat
        fight.go
    inventory/                 // Inventaires
        inventory.go
    item/                      // Objets du jeu
        item.go
    places/                    // Lieux spéciaux
        places.go
    sorts/                     // Sorts magiques
        sorts.go
    ui/                        // Interface utilisateur
        ui.go
    utils/                     // Fonctions utilitaires
        utils.go
    world/                     // Génération du monde
        world.go


## 📄 Explication Détaillée de Chaque Fichier {#fichiers-detailles}

1. Main.go : Le point d'entrée du code

```go

package main

// Importations de paquets (bibliothèques)
// - packages de la bibliothèque standard (fmt, math/rand, os, strings, time)
// - packages locaux du projet (character, classe, exploration, fight, places, ui, utils)
// Les imports locaux correspondent à d'autres fichiers Go dans le projet.
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

// Fonction principale : point d'entrée du programme
func main() {
	// Initialiser le générateur de nombres aléatoires.
	// rand.Seed prend une valeur (seed) qui est le point de départ des nombres aléatoires.
	// On utilise le temps actuel (UnixNano) pour avoir une seed différente à chaque lancement,
	// sinon rand produirait la même suite de nombres à chaque exécution.
	rand.Seed(time.Now().UnixNano())
	
	// Affiche/retourne un personnage obtenu après que l'utilisateur ait choisi dans le menu principal.
	// La fonction gererMenuPrincipal renvoie un pointeur vers un character.Character ou nil si l'utilisateur quitte.
	c := gererMenuPrincipal()
	if c == nil {
		// Si c est nil, l'utilisateur a choisi de quitter (ou la création a échoué) => on termine le programme.
		return
	}

	// On apelle une fonction du personnage pour initialiser l'état de la map/carte.
	// Ici, InitialiserEtatMap est une fonction qui prépare les données de la carte
	// (ex : positions découvertes, monstres présents, etc.). On utilise c (pointeur) pour modifier l'état du personnage.
	c.InitialiserEtatMap()

	// Gérer l'introduction / tutoriel ou reprendre l'aventure si le joueur a déjà fait l'intro.
	// La fonction retourne true si le jeu peut continuer, false si le joueur est mort pendant le tutoriel.
	if !executerIntroductionOuReprise(c) {
		return // Le joueur a été vaincu pendant le tutoriel => quitter le jeu
	}
	
	// Sauvegarder le personnage avant de commencer/continuer l'exploration
	sauvegarderPersonnageAvecMessage(c, "avant de commencer l'exploration")
	
	// Lancer le système d'exploration (boucle principale d'exploration du jeu)
	exploration.ExplorerMap(c)
}

// executerIntroductionOuReprise gère l'introduction pour un nouveau joueur ou la reprise d'aventure
// Paramètre : c -> pointeur vers le personnage (character.Character)
// Retour : bool -> true si on peut continuer, false si le personnage a été vaincu durant l'introduction
func executerIntroductionOuReprise(c *character.Character) bool {
	// La méthode AIntroEffectuee() est un getter qui indique si l'intro a déjà été faite.
	if !c.AIntroEffectuee() {
		// Si l'introduction n'a pas été faite, lancer le tutoriel.
		return executerTutoriel(c)
	} else {
		// Sinon, afficher un message de reprise et continuer
		reprendreAventure(c)
		return true
	}
}

// executerTutoriel lance le tutoriel d'introduction
// Cette fonction contient la logique d'introduction, de dialogue, d'un combat tutoriel et la gestion
// des conséquences (quête complétée / potion récupérée / mort pendant le tutoriel).
func executerTutoriel(c *character.Character) bool {
	// Récupérer les scènes d'introduction depuis le package places
	scenes := places.GetIntroDialogue()
	// Vérifier s'il y a au moins une scène
	if len(scenes) > 0 {
		// Sert a afficher des scènes, a été implémenter pour pouvoir ajouter d'autres scènes à l'avenir mais ce n'est pas présent par manque de temps
		scene := scenes[0] // Prendre la première (et unique) scène
		
		// Afficher le titre et la description de la scène
		fmt.Println("\n==== " + scene.Titre + " ====")
		fmt.Println(scene.Description)
		
		// Boucle de dialogue jusqu'à ce que le joueur choisisse une option particulière (ex: "On peut y aller !")
		for {
			// Afficher un menu d'options (ui.AfficherMenu affiche un menu avec des numéros et texte)
			ui.AfficherMenu("Choisissez une option", scene.Options)
			// utils.ScanChoice lit un choix depuis l'entrée standard et vérifie la validité
			choix := utils.ScanChoice("Votre choix : ", scene.Options)
			
			// Exécuter l'action correspondante au choix du joueur.
			// scene.Actions est un slice de fonctions.
			// Les indices en Go commencent à 0, d'où l'utilisation choix-1.
			scene.Actions[choix-1](c)
			
			// Si le joueur choisit l'option 4, on sort de la boucle.
			if choix == 4 {
				break
			}
			
			// Pause pour que le joueur lise la réponse. fmt.Scanln attend une entrée (touche Entrée).
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		}
	}

	// Proposer la quête du tutoriel (places.ProposerQueteTutoriel peut demander au joueur et renvoyer true/false)
	queteAcceptee := places.ProposerQueteTutoriel(c)
	
	// Récupérer la quête tutoriel et l'ennemi associé pour le combat.
	// Ici on ignore un des retours (utilisation du "_")
	quete, _, ennemi := places.GetTutorielCombat()
	// Lancer le combat tutoriel. fight.Fight modifiera l'état du personnage (points de vie, etc.).
	fight.Fight(c, ennemi)
	
	// Vérifier si le joueur a survécu au combat
	if c.Pdv > 0 {
		// Si le joueur est encore en vie
		// Si la quête avait été acceptée, la compléter automatiquement
		if queteAcceptee {
			c.CompleterQuete(quete)
			fmt.Println("\n✅ Quête accomplie automatiquement !")
			// Rendre automatiquement la quête du tutoriel : donner récompense, marquer comme rendue, etc.
			c.RendreQuete(quete)
		} else {
			// Si la quête n'avait pas été acceptée, donner une potion comme butin
			fmt.Println("\n💰 Le chacha laisse tomber 1 potion en mourant !")
			c.Inventaire.Potions++
		}
		
		// Marquer l'introduction comme effectuée (cela évite de refaire le tutoriel au prochain lancement)
		c.MarquerIntroEffectuee()
		
		// Messages de félicitations et pause avant de commencer l'aventure
		fmt.Println("\n🎉 Félicitations ! Vous avez terminé le tutoriel !")
		fmt.Println("Le monde s'ouvre maintenant à vous...")
		fmt.Println("\nAppuyez sur Entrée pour commencer votre aventure...")
		fmt.Scanln()
		return true
	} else {
		// Le joueur n'a plus de PV => défaite
		fmt.Println("\n💀 Vous avez été vaincu pendant le tutoriel...")
		fmt.Println("Le jeu se termine ici. Réessayez !")
		return false
	}
}

// reprendreAventure affiche le message de reprise pour un joueur existant
// Ici on récupère la position stockée dans le personnage et on l'affiche (x+1, y+1 pour afficher 1-based coordinates)
func reprendreAventure(c *character.Character) {
	fmt.Println("\n🔄 Reprise de votre aventure...")
	x, y := c.ObtenirPosition()
	// On ajoute +1 pour afficher une position plus naturelle pour l'utilisateur (commence souvent à 1)
	fmt.Printf("Vous êtes à la position (%d, %d) sur la carte.\n", x+1, y+1)
}

// gererMenuPrincipal affiche le menu principal et gère les choix de l'utilisateur
// Retourne un pointeur vers character.Character si l'utilisateur veut jouer, ou nil s'il choisit de quitter
func gererMenuPrincipal() *character.Character {
	for {
		// Affiche le header du menu principal
		afficherMenuPrincipalJeu()
		// Demande le choix de l'utilisateur
		choix := demanderChoixMenuPrincipal()
		
		// Traite le choix et éventuellement crée/charge un personnage
		personnage := traiterChoixMenuPrincipal(choix)
		if personnage != nil || choix == 3 {
			// Si on a un personnage prêt ou si l'utilisateur a choisi Quitter (3), retourner
			return personnage // Retourne le personnage ou nil (pour quitter)
		}
		// Si personnage est nil et choix != 3, on recommence la boucle pour redemander
	}
}

// centrerTexteAvecLargeur centre un texte dans une largeur spécifique
// utile pour l'affichage du header du menu
func centrerTexteAvecLargeur(texte string, largeur int) string {
	if len(texte) >= largeur {
		// Si le texte est plus long que la largeur, on le renvoie tel quel (pas de découpage ici)
		return texte
	}
	// Calculer le nombre d'espaces (padding) à mettre avant le texte pour le centrer
	padding := (largeur - len(texte)) / 2
	return strings.Repeat(" ", padding) + texte
}

// afficherMenuPrincipalJeu affiche le titre et les options du menu principal
// Cette fonction calcule une largeur raisonnable pour le menu et utilise ui.AfficherMenuSimple
func afficherMenuPrincipalJeu() {
	titre := "WORLD OF MILOUSQUES"
	soustitre := "Une aventure pleines de Milousqueries !"
	
	// Définir les options du menu (slice de chaînes)
	options := []string{
		"Créer un nouveau personnage",
		"Charger un personnage existant",
		"Quitter le jeu",
	}
	
	// Simuler le calcul de largeur utilisé par ui.AfficherMenuSimple
	// Ici on calcule la longueur maximale d'une ligne pour donner une belle bordure
	largeurContenu := len("MENU PRINCIPAL")
	for i, opt := range options {
		ligne := fmt.Sprintf("%d) %s", i+1, opt)
		if len(ligne) > largeurContenu {
			largeurContenu = len(ligne)
		}
	}
	// S'assurer d'une largeur minimale pour l'esthétique
	if largeurContenu < 30 {
		largeurContenu = 30
	}
	// Limiter la largeur maximale car cela cause des problèmes d'affichage dans les petits terminaux
	if largeurContenu > 50 {
		largeurContenu = 50
	}
	// Ajouter des marges / bordures
	largeurMenu := largeurContenu + 4 + 2 // +4 pour marge interne, +2 pour les bordures du menu
	
	// Afficher une ligne de séparation, le titre centré, etc.
	fmt.Println("\n" + strings.Repeat("=", largeurMenu))
	fmt.Println(centrerTexteAvecLargeur(titre, largeurMenu))
	fmt.Println(centrerTexteAvecLargeur(soustitre, largeurMenu))
	fmt.Println(strings.Repeat("=", largeurMenu))
	
	// Utiliser la fonction UI pour afficher les options proprement (numérotation, couleurs, etc.)
	ui.AfficherMenuSimple("MENU PRINCIPAL", options)
}

// demanderChoixMenuPrincipal demande et retourne le choix de l'utilisateur
// Elle s'appuie sur utils.ScanChoice qui gère la lecture et la validation
func demanderChoixMenuPrincipal() int {
	options := []string{"Créer un personnage", "Charger un personnage existant", "Quitter"}
	return utils.ScanChoice("Entrez votre choix : ", options)
}

// traiterChoixMenuPrincipal traite le choix de l'utilisateur et exécute l'action correspondante
// Retourne un pointeur vers un personnage si la création/chargement a réussi, nil sinon
func traiterChoixMenuPrincipal(choix int) *character.Character {
	switch choix {
	case 1:
		return gererCreationPersonnage()
	case 2:
		return reprendrePersonnage()
	case 3:
		// L'utilisateur choisit de quitter
		fmt.Println("👋 Au revoir et à bientôt dans World of Milousques !")
		return nil
	default:
		// Choix invalide (par sécurité)
		fmt.Println("❌ Choix invalide, veuillez réessayer.")
		return nil
	}
}

// gererCreationPersonnage gère la création d'un personnage avec validation
// Si la création fonctionne, retourne un pointeur vers character.Character
func gererCreationPersonnage() *character.Character {
	c := creerPersonnage()
	// Si le nom est vide, la création a échoué => retourner nil
	if c.Nom != "" {
		return &c
	}
	return nil // Création échouée
}

// creerPersonnage réalise les étapes concrètes de la création : lecture du nom, choix de la classe, initialisation
func creerPersonnage() character.Character {
	// Lire une chaîne (nom) depuis l'utilisateur ; utils.ScanString protège la longueur minimale
	nom := utils.ScanString("Entrez le nom de votre personnage : ", 1)

	// Récupérer les classes disponibles depuis le package classe
	classes := classe.GetClassesDisponibles()
	// Préparer les options à afficher (slice de string)
	classOptions := make([]string, len(classes))
	for i, cl := range classes {
		// fmt.Sprintf permet de construire une chaîne formatée
		classOptions[i] = fmt.Sprintf("%s (PV max : %d, Mana max : %d)", cl.Nom, cl.Pvmax, cl.ManaMax)
	}

	// Afficher le menu de classes et demander le choix
	ui.AfficherMenu("Choisissez la classe de votre personnage", classOptions)
	choix := utils.ScanChoice("Entrez le numéro de la classe : ", classOptions)

	// Récupérer la classe choisie (attention : indexation 0-based)
	classeChoisie := classes[choix-1]
	// Initialiser le personnage avec les valeurs de base
	c := character.InitCharacter(nom, classeChoisie, 1, classeChoisie.Pvmax, classeChoisie.Pvmax)

	fmt.Println("Personnage créé !")
	// Afficher le personnage (fonction d'affichage locale)
	afficherPersonnage(&c)

	// Créer le dossier de sauvegarde s'il n'existe pas et sauvegarder le personnage
	gererSauvegardePremiereFois(&c)

	return c
}

// reprendrePersonnage permet de charger un personnage existant depuis les fichiers de sauvegarde
func reprendrePersonnage() *character.Character {
	// Afficher les sauvegardes disponibles
	afficherSauvegardesDisponibles()
	
	// Demander le nom du personnage à charger
	nom := utils.ScanString("Entrez le nom du personnage à charger : ", 1)
	c := chargerPersonnageAvecMessage(nom)
	if c == nil {
		// Echec du chargement
		return nil
	}

	// Afficher le personnage chargé
	afficherPersonnage(c)
	return c
}

// afficherPersonnageComplet affiche toutes les informations détaillées d'un personnage
func afficherPersonnageComplet(c *character.Character) {
	// Afficher un résumé compact (nom, classe, PV, mana, XP...)
	afficherPersonnageResume(c)
	
	// Affichage détaillé de l'équipement
	afficherEquipementDetaille(c)
	
	// Afficher les sorts disponibles si la classe en possède
	if len(c.Classe.Sorts) > 0 {
		fmt.Println("\nSorts disponibles :")
		for _, s := range c.Classe.Sorts {
			// Affichage d'une ligne par sort : nom, dégâts, coût en mana
			fmt.Printf("- %s (Dégâts : %d, Coût en mana : %d)\n", s.Nom, s.Degats, s.Cout)
		}
	}
}

// afficherPersonnageResume affiche un résumé compact d'un personnage
func afficherPersonnageResume(c *character.Character) {
	// %s : chaîne, %d : entier (formatage printf)
	fmt.Printf("⚔️  %s (%s niveau %d)\n", c.Nom, c.Classe.Nom, c.Niveau)
	fmt.Printf("   PV: %d/%d | Mana: %d/%d | XP: %d/%d\n", 
		c.Pdv, c.PdvMax, c.Mana, c.ManaMax, c.Experience, c.CalculerXPRequis())
	fmt.Printf("   💰 %d or | 🧆 %d potions | 🗺️ %d/25 zones\n", 
		c.Argent, c.Inventaire.Potions, c.ObtenirNombreZonesDecouvertes())
}

// afficherEquipementDetaille affiche l'équipement d'un personnage de manière lisible
func afficherEquipementDetaille(c *character.Character) {
	equipements := []string{} // slice vide pour collecter les noms d'équipement
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
		// Joindre les éléments du slice avec une virgule pour l'affichage
		fmt.Println("\nÉquipement :", strings.Join(equipements, ", "))
		// Calculer et afficher les bonus d'attaque et de défense apportés par l'équipement
		bonusAttaque := c.CalculerAttaqueBonus()
		bonusDefense := c.CalculerDefenseBonus()
		if bonusAttaque > 0 || bonusDefense > 0 {
			fmt.Printf("Bonus : +%d Attaque, +%d Défense\n", bonusAttaque, bonusDefense)
		}
	}
}

// afficherPersonnage maintient la compatibilité - alias pour afficherPersonnageComplet
// Ceci simplifie l'appel depuis d'autres parties du code qui utilisent encore afficherPersonnage
func afficherPersonnage(c *character.Character) {
	afficherPersonnageComplet(c)
}

// afficherSauvegardesDisponibles affiche un aperçu des personnages sauvegardés
func afficherSauvegardesDisponibles() {
	fmt.Println("\n💾 === SAUVEGARDES DISPONIBLES === 💾")
	
	// Lire le dossier "saves" (où sont stockées les sauvegardes)
	files, err := os.ReadDir("saves")
	if err != nil {
		// Si ReadDir renvoie une erreur, le dossier n'existe probablement pas ou il y a un problème d'accès
		fmt.Println("Aucune sauvegarde trouvée.")
		return
	}
	
	aucuneSauvegarde := true
	for _, file := range files {
		// On ne s'intéresse qu'aux fichiers .json (extension utilisée pour stocker les sauvegardes)
		if strings.HasSuffix(file.Name(), ".json") {
			aucuneSauvegarde = false
			// Retirer l'extension pour avoir le nom du personnage
			nomPersonnage := strings.TrimSuffix(file.Name(), ".json")
			
			// Charger temporairement la sauvegarde pour afficher des informations utiles
			c, err := character.Charger(nomPersonnage)
			if err == nil {
				// Utiliser la fonction réutilisable pour l'affichage de base
				fmt.Println()
				afficherPersonnageResume(c)
				// Afficher des informations additionnelles propres aux sauvegardes
				afficherInfosSauvegarde(c)
			}
		}
	}
	
	if aucuneSauvegarde {
		fmt.Println("Aucune sauvegarde trouvée.")
	}
	
	fmt.Println()
}

// afficherInfosSauvegarde affiche les informations spécifiques aux sauvegardes (quêtes actives, équipement...)
func afficherInfosSauvegarde(c *character.Character) {
	// Compter les quêtes actives : on parcourt c.Quetes et on incrémente si !Rendue
	totalQuetes := 0
	for _, q := range c.Quetes {
		if !q.Rendue {
			totalQuetes++
		}
	}
	
	// Afficher le nombre de quêtes actives
	fmt.Printf("   🎒 %d quêtes actives", totalQuetes)
	
	// Compter le nombre d'équipements (fonction utilitaire locale)
	equipementCount := compterEquipement(c)
	if equipementCount > 0 {
		fmt.Printf(" | 🛡️  %d équipements", equipementCount)
		if c.ArmeEquipee != nil {
			fmt.Printf(" (Arme: %s)", c.ArmeEquipee.Nom)
		}
	}
	fmt.Println()
}

// compterEquipement compte le nombre de pièces d'équipement équipées
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

//Fonctions utilitaires pour la gestion d'erreurs
// gererSauvegardePremiereFois gère la création du dossier de sauvegarde et la première sauvegarde
func gererSauvegardePremiereFois(c *character.Character) {
	// creerDossierSauvegarde s'assure que le dossier "saves" existe
	if err := creerDossierSauvegarde(); err != nil {
		// Si une erreur survient, on l'affiche et on retourne (ne pas paniquer)
		fmt.Println("⚠️  Erreur lors de la création du dossier de sauvegarde :", err)
		return
	}
	
	// Sauvegarder le personnage et afficher un message indiquant le contexte
	sauvegarderPersonnageAvecMessage(c, "lors de la création du personnage")
}

// creerDossierSauvegarde crée le dossier de sauvegarde s'il n'existe pas
// os.MkdirAll crée tous les dossiers parents nécessaires et ne renvoie pas d'erreur si le dossier existe déjà
func creerDossierSauvegarde() error {
	return os.MkdirAll("saves", os.ModePerm)
}

// sauvegarderPersonnageAvecMessage sauvegarde un personnage avec un message contextuel en cas d'erreur
func sauvegarderPersonnageAvecMessage(c *character.Character, contexte string) {
	if err := c.Sauvegarder(); err != nil {
		// fmt.Printf permet d'inclure la variable err dans la chaîne de format
		fmt.Printf("⚠️  Erreur lors de la sauvegarde %s : %v\n", contexte, err)
	} else {
		fmt.Printf("✅ Personnage sauvegardé avec succès\n")
	}
}

// chargerPersonnageAvecMessage charge un personnage avec une gestion d'erreur explicite
func chargerPersonnageAvecMessage(nom string) *character.Character {
	c, err := character.Charger(nom)
	if err != nil {
		// On affiche un message d'erreur utile pour l'utilisateur
		fmt.Printf("❌ Erreur lors du chargement du personnage '%s' : %v\n", nom, err)
		fmt.Println("Vérifiez que le nom est correct et que la sauvegarde existe.")
		return nil
	}
	
	fmt.Printf("✅ Personnage '%s' chargé avec succès !\n", nom)
	return c
}

```

2. Banque.go : gère le système de stockage personnel pour le joueur

```go

package banque

import (
	"encoding/json" // Pour convertir la structure Banque en JSON (sauvegarde) et lire depuis JSON (chargement)
	"fmt"          // Pour afficher du texte dans le terminal
	"os"           // Pour gérer les fichiers (ouvrir, créer, vérifier l'existence...)
	"world_of_milousques/character" // Import du package qui gère les personnages
	"world_of_milousques/item"      // Import du package qui gère les objets
	"world_of_milousques/ui"        // Import pour afficher des menus
	"world_of_milousques/utils"     // Import pour les fonctions utilitaires (saisie utilisateur...)
)

// Banque représente le coffre-fort d'un joueur
// Elle contient le nom du propriétaire, une liste d'objets et une capacité maximale
// Les champs sont annotés avec `json:"..."` pour que Go sache comment les convertir en JSON lors de la sauvegarde
// Exemple : Proprietaire devient "proprietaire" dans le fichier JSON
//
type Banque struct {
	Proprietaire string      `json:"proprietaire"` // Nom du joueur qui possède la banque
	Objets       []item.Item `json:"objets"`       // Liste des objets stockés dans la banque
	MaxCapacite  int         `json:"max_capacite"` // Capacité maximale (200 objets)
}

// NewBanque crée une nouvelle banque pour un joueur
// Elle initialise la banque avec un propriétaire, une liste vide d'objets et une capacité par défaut de 200
func NewBanque(proprietaire string) *Banque {
	return &Banque{ // Retourne un pointeur vers une nouvelle Banque
		Proprietaire: proprietaire,
		Objets:       []item.Item{}, // Slice vide = aucun objet au début
		MaxCapacite:  200,           // Capacité fixée à 200
	}
}

// ChargerBanque charge la banque d'un joueur depuis un fichier JSON
// Si le fichier n'existe pas, on crée une nouvelle banque vide
func ChargerBanque(proprietaire string) (*Banque, error) {
	filename := "saves/banque_" + proprietaire + ".json" // Nom du fichier basé sur le joueur

	// Vérifier si le fichier existe avec os.Stat
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Si le fichier n'existe pas -> retourner une nouvelle banque vide
		return NewBanque(proprietaire), nil
	}

	// Sinon on ouvre le fichier existant
	file, err := os.Open(filename)
	if err != nil {
		return nil, err // Retourner une erreur si ouverture impossible
	}
	defer file.Close() // On s'assure que le fichier sera fermé à la fin

	var banque Banque // Variable temporaire qui recevra les données JSON
	decoder := json.NewDecoder(file) // Création d'un décodeur JSON
	err = decoder.Decode(&banque) // On remplit la variable banque avec le contenu du fichier
	if err != nil {
		return nil, err
	}

	return &banque, nil // Retourne un pointeur vers la banque chargée
}

// Sauvegarder sauvegarde la banque dans un fichier JSON
func (b *Banque) Sauvegarder() error {
	filename := "saves/banque_" + b.Proprietaire + ".json" // Fichier propre au joueur

	file, err := os.Create(filename) // Crée ou écrase le fichier
	if err != nil {
		return err // Retourne l'erreur si problème
	}
	defer file.Close() // Ferme le fichier automatiquement en fin de fonction

	encoder := json.NewEncoder(file) // Création d'un encodeur JSON
	encoder.SetIndent("", "  ")       // Pour rendre le JSON lisible (indentation)
	return encoder.Encode(b)         // Écriture de la structure Banque dans le fichier
}

// AjouterObjet ajoute un objet à la banque
// Retourne false si la banque est déjà pleine
func (b *Banque) AjouterObjet(objet item.Item) bool {
	if len(b.Objets) >= b.MaxCapacite { // Vérifie si la banque est pleine
		return false
	}

	b.Objets = append(b.Objets, objet) // Ajoute l'objet à la slice
	return true
}

// RetirerObjet retire un objet de la banque en fonction de son index
// Retourne l'objet retiré et true si succès, sinon un objet vide et false
func (b *Banque) RetirerObjet(index int) (item.Item, bool) {
	if index < 0 || index >= len(b.Objets) { // Vérifie si l'index est valide
		return item.Item{}, false
	}

	objet := b.Objets[index] // On garde l'objet à retirer

	// Recrée une nouvelle slice sans l'objet retiré
	nouveauxObjets := make([]item.Item, 0)
	for i, obj := range b.Objets {
		if i != index {
			nouveauxObjets = append(nouveauxObjets, obj)
		}
	}
	b.Objets = nouveauxObjets

	return objet, true
}

// AfficherBanque gère l'interface de la banque (menu, choix du joueur...)
func AfficherBanque(joueur *character.Character) {
	// On charge la banque du joueur
	banque, err := ChargerBanque(joueur.Nom)
	if err != nil {
		fmt.Printf("Erreur lors du chargement de votre coffre : %v\n", err)
		return
	}

	// Boucle infinie jusqu'à ce que le joueur quitte
	for {
		fmt.Println("\n🏦 === BANQUE ROYALE D'ASTRAB === 🏦")
		fmt.Printf("Banquier Salomon : Bienvenue %s ! Votre coffre-fort vous attend.\n", joueur.Nom)
		fmt.Printf("💰 Capacité du coffre : %d/%d objets\n", len(banque.Objets), banque.MaxCapacite)
		fmt.Printf("🎒 Votre inventaire : %d/100 objets\n", len(joueur.Inventaire.Items))

		// Liste des actions possibles
		options := []string{
			"🏦 Déposer des objets",
			"📤 Retirer des objets",
			"📋 Voir le contenu du coffre",
			"🎒 Voir mon inventaire",
			"🚪 Quitter la banque",
		}

		ui.AfficherMenu("Services bancaires", options) // Affiche le menu
		choix := utils.ScanChoice("Que souhaitez-vous faire ? ", options) // Demande le choix

		switch choix {
		case 1:
			deposerObjets(joueur, banque)
		case 2:
			retirerObjets(joueur, banque)
		case 3:
			afficherContenuBanque(banque)
		case 4:
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 5:
			// Sauvegarde avant de quitter
			if err := banque.Sauvegarder(); err != nil {
				fmt.Printf("Erreur lors de la sauvegarde : %v\n", err)
			} else {
				fmt.Println("Banquier Salomon : Vos biens sont en sécurité ! À bientôt !")
			}
			return
		}
	}
}

// deposerObjets gère le dépôt d'objets dans la banque
func deposerObjets(joueur *character.Character, banque *Banque) {
	// Vérifie si l'inventaire du joueur est vide
	if len(joueur.Inventaire.Items) == 0 {
		fmt.Println("❌ Votre inventaire est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	// Vérifie si la banque est pleine
	if len(banque.Objets) >= banque.MaxCapacite {
		fmt.Println("❌ Votre coffre est plein ! Retirez d'abord des objets.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	fmt.Println("\n💰 === DÉPOSER DES OBJETS === 💰")
	fmt.Printf("Espace disponible dans le coffre : %d objets\n\n", banque.MaxCapacite-len(banque.Objets))

	// Grouper les objets identiques pour simplifier l'affichage
	objetsGroupes := make(map[string]GroupeObjet)
	for i, objet := range joueur.Inventaire.Items {
		if groupe, existe := objetsGroupes[objet.Nom]; existe {
			groupe.Quantite++
			groupe.Indices = append(groupe.Indices, i)
			objetsGroupes[objet.Nom] = groupe
		} else {
			objetsGroupes[objet.Nom] = GroupeObjet{
				Item:     objet,
				Quantite: 1,
				Indices:  []int{i},
			}
		}
	}

	// Préparer les options du menu
	options := make([]string, 0)
	groupes := make([]GroupeObjet, 0)

	for _, groupe := range objetsGroupes {
		options = append(options, fmt.Sprintf("%s (%dx)", groupe.Item.Nom, groupe.Quantite))
		groupes = append(groupes, groupe)
	}
	options = append(options, "Retour")

	ui.AfficherMenu("Choisir un objet à déposer", options)
	choix := utils.ScanChoice("Quel objet voulez-vous déposer ? ", options)

	if choix == len(options) { // Si l'utilisateur choisit "Retour"
		return
	}

	groupeChoisi := groupes[choix-1]

	// Si plusieurs exemplaires -> demander combien déposer
	quantiteADeposer := 1
	maxDeposable := min(groupeChoisi.Quantite, banque.MaxCapacite-len(banque.Objets))

	if groupeChoisi.Quantite > 1 && maxDeposable > 1 {
		quantiteADeposer = utils.ScanInt(
			fmt.Sprintf("Combien voulez-vous en déposer ? (max %d) : ", maxDeposable),
			1, maxDeposable)
	}

	// Effectuer le dépôt
	for i := 0; i < quantiteADeposer; i++ {
		if !banque.AjouterObjet(groupeChoisi.Item) {
			fmt.Println("❌ Le coffre est plein !")
			break
		}
	}

	// Retirer les objets déposés de l'inventaire du joueur
	retirerObjetsInventaire(joueur, groupeChoisi.Item.Nom, quantiteADeposer)

	fmt.Printf("✅ %dx %s déposé avec succès dans votre coffre !\n", quantiteADeposer, groupeChoisi.Item.Nom)
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// retirerObjets gère le retrait d'objets de la banque
func retirerObjets(joueur *character.Character, banque *Banque) {
	if len(banque.Objets) == 0 {
		fmt.Println("❌ Votre coffre est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	if len(joueur.Inventaire.Items) >= 100 {
		fmt.Println("❌ Votre inventaire est plein ! Videz d'abord votre inventaire.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	fmt.Println("\n📤 === RETIRER DES OBJETS === 📤")
	fmt.Printf("Espace disponible dans l'inventaire : %d objets\n\n", 100-len(joueur.Inventaire.Items))

	// Affiche le contenu du coffre
	afficherContenuBanque(banque)

	if len(banque.Objets) == 0 {
		return
	}

	choix := utils.ScanInt("Quel objet voulez-vous retirer ? (numéro) : ", 1, len(banque.Objets))

	objet, success := banque.RetirerObjet(choix - 1)
	if success {
		joueur.Inventaire.Items = append(joueur.Inventaire.Items, objet)
		fmt.Printf("✅ %s retiré avec succès de votre coffre !\n", objet.Nom)
	} else {
		fmt.Println("❌ Erreur lors du retrait de l'objet.")
	}

	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// afficherContenuBanque affiche tous les objets actuellement dans la banque
func afficherContenuBanque(banque *Banque) {
	fmt.Printf("\n📋 === CONTENU DU COFFRE === 📋\n")

	if len(banque.Objets) == 0 {
		fmt.Println("Votre coffre est vide.")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	for i, objet := range banque.Objets {
		fmt.Printf("%d. %s | Poids: %d | Effet: %s | Valeur: %d or\n",
			i+1, objet.Nom, objet.Poids, objet.Effet, objet.Valeur)
	}

	fmt.Printf("\nTotal : %d/%d objets\n", len(banque.Objets), banque.MaxCapacite)

	if len(banque.Objets) < 20 { // Si peu d'objets, pas besoin de pause
		return
	}

	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// GroupeObjet représente un groupe d'objets identiques avec leurs indices
// Utile pour gérer plusieurs exemplaires d'un même objet

type GroupeObjet struct {
	Item     item.Item
	Quantite int
	Indices  []int
}

// retirerObjetsInventaire retire une certaine quantité d'un objet de l'inventaire du joueur
func retirerObjetsInventaire(joueur *character.Character, nomObjet string, quantite int) {
	retirees := 0
	nouvelInventaire := make([]item.Item, 0)

	for _, objet := range joueur.Inventaire.Items {
		if objet.Nom == nomObjet && retirees < quantite {
			retirees++ // On retire cet objet (on ne l'ajoute pas au nouvel inventaire)
		} else {
			nouvelInventaire = append(nouvelInventaire, objet)
		}
	}

	joueur.Inventaire.Items = nouvelInventaire // Remplace l'inventaire par le nouveau
}

// min retourne le plus petit entre deux entiers (utile pour limiter une quantité)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

3. Character.go : tout ce qui touche au personnage, aux sauvegardes et aux quetes

```go

import (
	"encoding/json"   // Permet d’encoder/décoder les données en JSON (sauvegarde/chargement)
	"fmt"             // Permet d’afficher du texte à l’écran (ex: fmt.Println)
	"os"              // Fournit des fonctions pour manipuler les fichiers (ouvrir, créer…)

	"world_of_milousques/classe"     // Import du package classe (définit les classes de personnages : mage, guerrier…)
	"world_of_milousques/inventory"  // Import du package inventaire (gère les objets du joueur)
)

// Quete représente une mission que le joueur peut accomplir
// Chaque quête a : un nom, un statut (accomplie ou non), et une récompense
// Exemple : "Tuer 5 Moutmouts", accomplie=false, récompense="300 golds"
type Quete struct {
	Nom        string // Nom de la quête
	Accomplie  bool   // true si terminée, false sinon
	Recompense string // Ce que le joueur reçoit s’il termine la quête
}

// Character représente un personnage jouable
// On utilise des tags JSON (`json:"nom"`) pour que les champs puissent être sauvegardés dans un fichier JSON
// Exemple : Nom="Arthur", Niveau=1, Pdv=100, Mana=50, Classe="Guerrier", Inventaire={}, Quetes=[]
type Character struct {
	Nom        string                `json:"nom"`       // Nom du joueur
	Niveau     float64               `json:"niveau"`    // Niveau du personnage (float64 permet d’avoir des demi-niveaux)
	Pdv        int                   `json:"pdv"`      // Points de vie (santé actuelle)
	Mana       int                   `json:"mana"`     // Points de mana (ressource magique)
	Classe     classe.Classe         `json:"classe"`   // Classe du joueur (définie dans le package classe)
	Inventaire inventory.Inventaire  `json:"inventaire"` // Inventaire avec objets, potions, etc.
	Quetes     []Quete               `json:"quetes"`    // Liste de quêtes du joueur
}

// InitCharacter initialise un nouveau personnage
// On lui passe : nom, classe, niveau, pdv actuel, et pdv max (non utilisé ici mais pourrait l’être)
// Retourne un objet Character prêt à jouer
func InitCharacter(nom string, c classe.Classe, niveau float64, pdv int, pdvmax int) Character {
	return Character{
		Nom:        nom,                 // Définit le nom du personnage
		Niveau:     niveau,              // Définit son niveau
		Pdv:        pdv,                 // Points de vie actuels
		Mana:       c.ManaMax,           // Initialise son mana au maximum défini par sa classe
		Classe:     c,                   // Associe la classe choisie
		Inventaire: inventory.Inventaire{}, // Crée un inventaire vide
		Quetes:     []Quete{},           // Aucune quête au départ
	}
}

// Sauvegarder écrit les données du personnage dans un fichier JSON
// Exemple : saves/Arthur.json
func (c *Character) Sauvegarder() error {
	filename := "saves/" + c.Nom + ".json" // Nom du fichier basé sur le nom du joueur

	file, err := os.Create(filename) // Crée ou écrase le fichier
	if err != nil {
		return err // Retourne l’erreur si problème de création
	}
	defer file.Close() // Ferme le fichier automatiquement à la fin

	encoder := json.NewEncoder(file) // Crée un encodeur JSON
	encoder.SetIndent("", "  ")      // Formate joliment le JSON (indentation)
	err = encoder.Encode(c)          // Écrit les données du personnage dans le fichier
	if err != nil {
		return err // Retourne l’erreur si problème d’écriture
	}

	fmt.Println("Personnage sauvegardé dans", filename)
	return nil // Succès
}

// Charger lit un fichier JSON et recrée un personnage à partir de son contenu
func Charger(nom string) (*Character, error) {
	filename := "saves/" + nom + ".json"

	file, err := os.Open(filename) // Ouvre le fichier
	if err != nil {
		return nil, err // Erreur si le fichier n’existe pas
	}
	defer file.Close()

	var c Character
	decoder := json.NewDecoder(file) // Crée un décodeur JSON
	err = decoder.Decode(&c)         // Remplit la variable c avec les données du fichier
	if err != nil {
		return nil, err // Erreur si le JSON est corrompu
	}

	fmt.Println("Personnage chargé depuis", filename)
	return &c, nil // Retourne un pointeur vers le personnage
}

// ProposerEtAjouterQuete ajoute une nouvelle quête à la liste du joueur
// Elle est automatiquement marquée comme non accomplie
tfunc (c *Character) ProposerEtAjouterQuete(nom string, recompense string) {
	c.Quetes = append(c.Quetes, Quete{Nom: nom, Accomplie: false, Recompense: recompense})
}

// CompleterQuete marque une quête comme accomplie et donne sa récompense
func (c *Character) CompleterQuete(nom string) {
	for i := range c.Quetes { // On parcourt toutes les quêtes du joueur
		if c.Quetes[i].Nom == nom { // Si on trouve la bonne quête
			c.Quetes[i].Accomplie = true // On la marque comme accomplie
			fmt.Println("Quête complétée :", nom)
			fmt.Println("Récompense :", c.Quetes[i].Recompense)

			// Exemple de récompense spécifique : donner une potion
			if c.Quetes[i].Recompense == "1 potion" {
				c.Inventaire.Potions++ // Ajoute une potion dans l’inventaire
				fmt.Println("Vous recevez 1 potion !")
			}
			break // On arrête la boucle après avoir trouvé
		}
	}
}

// AfficherQuetes affiche toutes les quêtes du joueur avec leur statut
func (c *Character) AfficherQuetes() {
	if len(c.Quetes) == 0 { // Si aucune quête
		fmt.Println("Aucune quête.")
		return
	}
	fmt.Println("Quêtes :")
	for _, q := range c.Quetes {
		status := "Non accomplie" // Par défaut
		if q.Accomplie {
			status = "Accomplie"
		}
		fmt.Printf("- %s : %s | Récompense : %s\n", q.Nom, status, q.Recompense)
	}
}

```

4. Classe.go : gère les classes jouables

```go

import (
	"world_of_milousques/classeitem" // Package local pour les objets spécifiques aux classes
	"world_of_milousques/sorts"      // Package local qui contient les définitions de sorts
)

// Classe regroupe les informations de base d'une classe de personnage. Chaque champ est exporté (première lettre en 
// majuscule) pour être accessible depuis d'autres packages et pour la sauvegarde json.
type Classe struct {
	Nom         string                     `json:"nom"`         // Nom lisible de la classe (ex: "Guerrier")
	Pvmax       int                        `json:"pv_max"`      // Points de vie maximum de la classe
	ManaMax     int                        `json:"mana_max"`    // Mana maximum disponible
	Sorts       []sorts.Sorts              `json:"sorts"`       // Slice (liste) des sorts disponibles pour la classe
	ClasseItems []classeitem.Classeitem   `json:"classe_items"`// Objets ou bonus spécifiques à la classe
}

// GetClasse retourne une valeur de type Classe en fonction d'un nom fourni.
func GetClasse(nom string) Classe {
	switch nom {
	case "Guerrier":
		// Exemple d'initialisation d'une structure : on précise chaque champ
		return Classe{
			Nom:     "Guerrier",
			Pvmax:   130,
			ManaMax: 70,
			// Sorts est une slice : on appelle sorts.GetSorts pour récupérer chaque sort
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Fracasser"),
				sorts.GetSorts("Briser"),
			},
			// Ici on initialise une slice vide d'objets de classe
			ClasseItems: []classeitem.Classeitem{},
		}

	case "Mage":
		return Classe{
			Nom:     "Mage",
			Pvmax:   70,
			ManaMax: 130,
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Boule de feu"),
				sorts.GetSorts("Explosion"),
			},
			ClasseItems: []classeitem.Classeitem{},
		}

	case "Voleur":
		return Classe{
			Nom:     "Voleur",
			Pvmax:   100,
			ManaMax: 100,
			Sorts: []sorts.Sorts{
				sorts.GetSorts("Coup bas"),
				sorts.GetSorts("Fourberie"),
			},
			ClasseItems: []classeitem.Classeitem{},
		}

	default:
		// Cas par défaut : si le nom n'est pas reconnu, on retourne une classe générique
		// C'est pratique pour éviter des erreurs panics, mais dans un jeu réel on
		// préférerait probablement renvoyer une erreur ou utiliser des constantes.
		return Classe{
			Nom:     nom, // on renvoie le nom fourni pour transparence
			Pvmax:   100,
			ManaMax: 100,
			Sorts:       []sorts.Sorts{},
			ClasseItems: []classeitem.Classeitem{},
		}
	}
}

// GetClassesDisponibles retourne un slice de Classe contenant toutes les classes proposées dans le jeu. Cette fonction est   
// utile pour afficher un menu de choix des classes lors de la création d'un personnage.
func GetClassesDisponibles() []Classe {
	// Ici on définit d'abord la liste des noms de classes disponibles.
	classesNoms := []string{"Guerrier", "Mage", "Voleur"}

	// Déclaration d'une slice vide de Classe
	var classes []Classe

	// Boucle for ... range : pour chaque nom dans classesNoms, on récupère la structure correspondante avec GetClasse et on  
    // l'ajoute à la slice classes.
	for _, nom := range classesNoms {
		classes = append(classes, GetClasse(nom))
	}

	return classes
}

```

5. Commerce.go : Permet l'achat d'équipements et de consommables a une marchand

```go

import (
	"fmt" // Sert pour afficher du texte dans la console
	"world_of_milousques/character" // On importe les personnages (pour gérer leur argent et inventaire)
	"world_of_milousques/item"      // On importe les objets (items) utilisables
	"world_of_milousques/ui"        // Pour afficher des menus interactifs
	"world_of_milousques/utils"     // Pour gérer les entrées utilisateur (choix, nombres...)
)

// Article représente un article vendu par le marchand, chaque article a un objet, un prix, un stock et une option illimité.
type Article struct {
	Item     item.Item // L'objet vendu
	Prix     int       // Le prix en pièces d'or
	Stock    int       // Combien d'exemplaires sont disponibles
	Illimite bool      // Si true, alors le stock est infini (ex: potions)
}

// Marchand représente un marchand avec son inventaire d'articles, contient son nom, une phrase de salutation, et la liste de 
// ses articles.
type Marchand struct {
	Nom      string    // Nom du marchand
	Salut    string    // Message de bienvenue
	Articles []Article // Liste des articles en vente
}

// GetMarchandAstrab retourne le marchand principal de la ville d'Astrab, ici, on définit ses articles manuellement.
func GetMarchandAstrab() Marchand {
	return Marchand{
		Nom:   "Maître Karim le Marchand",
		Salut: "Bienvenue dans ma boutique ! J'ai tout l'équipement qu'il vous faut, aventurier !",
		Articles: []Article{
			// Équipements en cuir (150 or chacun)
			{Item: item.NewItem("Casque en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			{Item: item.NewItem("Torse en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			{Item: item.NewItem("Jambières en Cuir"), Prix: 150, Stock: 5, Illimite: false},
			// Armes simples (250 or chacune)
			{Item: item.NewItem("Bâton Simple"), Prix: 250, Stock: 3, Illimite: false},
			{Item: item.NewItem("Épée Simple"), Prix: 250, Stock: 3, Illimite: false},
			{Item: item.NewItem("Dague Simple"), Prix: 250, Stock: 3, Illimite: false},
			// Potions (stock illimité)
			{Item: item.NewItem("Potion de Vie"), Prix: 50, Stock: 0, Illimite: true},
			{Item: item.NewItem("Potion de Mana"), Prix: 50, Stock: 0, Illimite: true},
		},
	}
}

// AfficherMarchand affiche le menu principal d'interaction avec le marchand, le joueur peut voir, acheter, vendre ou quitter.
func AfficherMarchand(joueur *character.Character) {
	// On récupère le marchand prédéfini (Astrab)
	marchand := GetMarchandAstrab()

	// Boucle infinie (while true) : tant que le joueur n'a pas choisi de quitter
	for {
		fmt.Printf("\n💰 === %s === 💰\n", marchand.Nom)
		fmt.Printf("%s\n", marchand.Salut)
		fmt.Printf("💳 Votre argent : %d pièces d'or\n", joueur.Argent)

		// Menu des options possibles
		options := []string{
			"Voir les articles à vendre",
			"Acheter un article",
			"Vendre mes objets",
			"Voir mon inventaire",
			"Quitter la boutique",
		}

		// On affiche le menu et récupère le choix du joueur
		ui.AfficherMenu("Boutique", options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)

		// On exécute l'action choisie via un switch
		switch choix {
		case 1:
			afficherArticles(marchand)
		case 2:
			acheterArticle(joueur, &marchand)
		case 3:
			vendreObjets(joueur)
		case 4:
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 5:
			fmt.Printf("%s : Merci de votre visite ! Revenez quand vous voulez !\n", marchand.Nom)
			return // Quitter la fonction => sortie de la boutique
		}
	}
}

// afficherArticles affiche tous les articles avec leur prix et leur disponibilité.
func afficherArticles(marchand Marchand) {
	fmt.Println("\n🛒 === ARTICLES DISPONIBLES === 🛒")

	// Parcours des articles avec leur index (i) et valeur (article)
	for i, article := range marchand.Articles {
		stockInfo := fmt.Sprintf("(%d en stock)", article.Stock)
		if article.Illimite {
			stockInfo = "(Stock illimité)"
		}

		disponible := "✅"
		if article.Stock == 0 && !article.Illimite {
			disponible = "❌ RUPTURE"
			stockInfo = ""
		}

		// Affiche ligne principale
		fmt.Printf("%d. %s %s - %d pièces d'or %s\n",
			i+1, disponible, article.Item.Nom, article.Prix, stockInfo)
		// Affiche effet de l'objet
		fmt.Printf("   %s\n", article.Item.Effet)
	}

	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// acheterArticle permet au joueur d'acheter un article du marchand.
func acheterArticle(joueur *character.Character, marchand *Marchand) {
	fmt.Println("\n💳 === ACHAT D'ARTICLE === 💳")
	fmt.Printf("Votre argent : %d pièces d'or\n", joueur.Argent)

	// Création du menu des options
	options := make([]string, 0)
	for i, article := range marchand.Articles {
		disponible := "✅"
		prixInfo := fmt.Sprintf("- %d or", article.Prix)

		// Si rupture
		if article.Stock == 0 && !article.Illimite {
			disponible = "❌"
			prixInfo = "- RUPTURE"
		// Si trop cher
		} else if joueur.Argent < article.Prix {
			disponible = "💸"
			prixInfo = fmt.Sprintf("- %d or (trop cher)", article.Prix)
		}

		// Ajout de l'option formatée
		options = append(options, fmt.Sprintf("%s %s %s",
			disponible, article.Item.Nom, prixInfo))
		_ = i // évite une erreur "unused variable"
	}
	options = append(options, "Retour")

	// Affichage du menu et choix
	ui.AfficherMenu("Choisir un article à acheter", options)
	choix := utils.ScanChoice("Quel article voulez-vous acheter ? ", options)

	if choix == len(options) {
		return // Retour
	}

	articleChoisi := &marchand.Articles[choix-1]

	// Vérification du stock et argent disponible
	if articleChoisi.Stock == 0 && !articleChoisi.Illimite {
		fmt.Println("❌ Cet article n'est plus en stock !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
	if joueur.Argent < articleChoisi.Prix {
		fmt.Printf("💸 Vous n'avez pas assez d'argent ! Il vous faut %d pièces d'or.\n",
			articleChoisi.Prix)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	// Confirmation avant achat
	fmt.Printf("\n🛒 Acheter : %s\n", articleChoisi.Item.Nom)
	fmt.Printf("Prix : %d pièces d'or\n", articleChoisi.Prix)
	fmt.Printf("Argent restant : %d pièces d'or\n", joueur.Argent-articleChoisi.Prix)

	options = []string{"Confirmer l'achat", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)

	if confirmation == 1 {
		// Effectuer l'achat : on retire l'argent
		joueur.Argent -= articleChoisi.Prix

		// Cas spéciaux : potions (elles augmentent un compteur, pas un item normal)
		if articleChoisi.Item.Nom == "Potion de Vie" {
			joueur.Inventaire.Potions++
		} else if articleChoisi.Item.Nom == "Potion de Mana" {
			joueur.Inventaire.PotionsMana++
		} else {
			// Autres objets ajoutés directement dans Items
			joueur.Inventaire.Items = append(joueur.Inventaire.Items, articleChoisi.Item)
		}

		// Réduction du stock si pas illimité
		if !articleChoisi.Illimite {
			articleChoisi.Stock--
		}

		fmt.Printf("✅ %s acheté avec succès !\n", articleChoisi.Item.Nom)
		fmt.Printf("Argent restant : %d pièces d'or\n", joueur.Argent)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// vendreObjets permet au joueur de vendre ses objets de l'inventaire.
func vendreObjets(joueur *character.Character) {
	if len(joueur.Inventaire.Items) == 0 {
		fmt.Println("❌ Votre inventaire est vide !")
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	fmt.Println("\n💰 === VENTE D'OBJETS === 💰")
	fmt.Printf("Votre argent actuel : %d pièces d'or\n", joueur.Argent)

	// Grouper les objets identiques dans une map (clé = nom de l'objet)
	objetsGroupes := make(map[string]GroupeObjet)
	for _, objet := range joueur.Inventaire.Items {
		if groupe, existe := objetsGroupes[objet.Nom]; existe {
			groupe.Quantite++
			objetsGroupes[objet.Nom] = groupe
		} else {
			objetsGroupes[objet.Nom] = GroupeObjet{
				Item:     objet,
				Quantite: 1,
				PrixVente: objet.Valeur / 2, // Vente à 50% du prix d'achat
			}
		}
	}

	// Créer les options de menu à partir de la map
	options := make([]string, 0)
	groupes := make([]GroupeObjet, 0)

	for _, groupe := range objetsGroupes {
		options = append(options, fmt.Sprintf("%s (%dx) - %d or chacun",
			groupe.Item.Nom, groupe.Quantite, groupe.PrixVente))
		groupes = append(groupes, groupe)
	}
	options = append(options, "Retour")

	ui.AfficherMenu("Choisir un objet à vendre", options)
	choix := utils.ScanChoice("Quel objet voulez-vous vendre ? ", options)

	if choix == len(options) {
		return // Retour
	}

	groupeChoisi := groupes[choix-1]

	// Demander combien vendre si plusieurs exemplaires
	quantiteAVendre := 1
	if groupeChoisi.Quantite > 1 {
		quantiteAVendre = utils.ScanInt(
			fmt.Sprintf("Combien voulez-vous en vendre ? (max %d) : ", groupeChoisi.Quantite),
			1, groupeChoisi.Quantite)
	}

	prixTotal := groupeChoisi.PrixVente * quantiteAVendre

	// Confirmation de la vente
	fmt.Printf("\n💰 Vendre : %dx %s\n", quantiteAVendre, groupeChoisi.Item.Nom)
	fmt.Printf("Prix total : %d pièces d'or\n", prixTotal)
	fmt.Printf("Argent après vente : %d pièces d'or\n", joueur.Argent+prixTotal)

	options = []string{"Confirmer la vente", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)

	if confirmation == 1 {
		// On retire les objets de l'inventaire
		retirerObjets(joueur, groupeChoisi.Item.Nom, quantiteAVendre)
		// On ajoute l'argent
		joueur.Argent += prixTotal

		fmt.Printf("✅ %dx %s vendu avec succès !\n", quantiteAVendre, groupeChoisi.Item.Nom)
		fmt.Printf("Vous avez gagné : %d pièces d'or\n", prixTotal)
		fmt.Printf("Argent total : %d pièces d'or\n", joueur.Argent)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// GroupeObjet représente un groupe d'objets identiques pour simplifier la vente.
type GroupeObjet struct {
	Item      item.Item // L'objet
	Quantite  int       // Combien d'exemplaires
	PrixVente int       // Prix de vente unitaire (50% du prix normal)
}

// retirerObjets retire une certaine quantité d'un objet de l'inventaire du joueur.
func retirerObjets(joueur *character.Character, nomObjet string, quantite int) {
	retirees := 0
	nouvelInventaire := make([]item.Item, 0)

	for _, objet := range joueur.Inventaire.Items {
		if objet.Nom == nomObjet && retirees < quantite {
			// On saute cet objet pour le retirer
			retirees++
		} else {
			nouvelInventaire = append(nouvelInventaire, objet)
		}
	}

	// Mise à jour de l'inventaire du joueur
	joueur.Inventaire.Items = nouvelInventaire
}

```
6. Craft.go : Gère tout le système de craft

```go

package craft

import (
	"fmt"                               // Pour afficher des textes et interagir avec l'utilisateur
	"world_of_milousques/character"    // Pour manipuler le personnage du joueur (inventaire, etc.)
	"world_of_milousques/item"         // Pour manipuler les objets (items) du jeu
	"world_of_milousques/ui"           // Pour afficher des menus interactifs
	"world_of_milousques/utils"        // Pour des fonctions utilitaires (scan de choix utilisateur)
)

// Ingredient représente un ingrédient nécessaire pour une recette c,haque ingrédient est défini par un item (nom, caractéristiques) et une quantité requise
 type Ingredient struct {
	Item     item.Item   // L’objet requis (ex : "Fer", "Bois", etc.)
	Quantite int         // La quantité de cet objet nécessaire pour la recette
}

// Recette représente une recette de craft, elle contient son nom, une description, la liste des ingrédients, le produit créé et la quantité produite
 type Recette struct {
	Nom             string        // Nom de la recette (ex : "Épée d'Expert")
	Description     string        // Description de l’objet crafté
	Ingredients     []Ingredient  // Liste des ingrédients nécessaires
	Produit         item.Item     // L’objet résultant du craft
	QuantiteProduit int           // Combien d’exemplaires sont produits par craft
}

// GetRecettesDisponibles retourne toutes les recettes de craft disponibles dans le jeu
func GetRecettesDisponibles() []Recette {
	return []Recette{
		// On définit ici toutes les recettes de manière statique
		// Exemple : un casque en métal requiert 10 bois, 10 fer, etc.
		{
			Nom:         "Casque en Métal",
			Description: "Protection de tête métallique résistante (+10 défense)",
			Ingredients: []Ingredient{
				{Item: item.NewItem("Bois"), Quantite: 10},
				{Item: item.NewItem("Fer"), Quantite: 10},
				{Item: item.NewItem("Blé"), Quantite: 10},
				{Item: item.NewItem("Laitue Vireuse"), Quantite: 10},
				{Item: item.NewItem("Pichon"), Quantite: 10},
			},
			Produit:         item.NewItem("Casque en Métal"),
			QuantiteProduit: 1,
		},
		// ... (les autres recettes suivent le même schéma)
	}
}

// AfficherForge affiche le menu principal de la forge
func AfficherForge(joueur *character.Character) {
	for {
		// Boucle infinie tant que le joueur reste dans la forge
		fmt.Println("\n🔨 === FORGE D'ASTRAB === 🔨")
		fmt.Println("Maître Forgeron : Bienvenue dans ma forge ! Que puis-je créer pour vous ?")

		// Options proposées au joueur
		options := []string{
			"Voir les recettes disponibles",
			"Crafter un objet",
			"Voir mon inventaire",
			"Quitter la forge",
		}

		// Affichage du menu via ui.AfficherMenu
		ui.AfficherMenu("Forge", options)
		// On demande le choix de l’utilisateur avec utils.ScanChoice
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)

		switch choix {
		case 1:
			// Voir les recettes disponibles
			afficherRecettes()
		case 2:
			// Crafter un objet
			crafterObjet(joueur)
		case 3:
			// Affiche l’inventaire du joueur
			joueur.Inventaire.Afficher()
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
		case 4:
			// Quitte la forge
			fmt.Println("Maître Forgeron : Revenez quand vous voulez !")
			return
		}
	}
}

// afficherRecettes affiche toutes les recettes disponibles au joueur
func afficherRecettes() {
	recettes := GetRecettesDisponibles() // On récupère la liste des recettes

	fmt.Println("\n📜 === RECETTES DISPONIBLES === 📜")
	for i, recette := range recettes {
		// Pour chaque recette, on affiche ses détails
		fmt.Printf("\n%d. %s\n", i+1, recette.Nom)
		fmt.Printf("   Description: %s\n", recette.Description)
		fmt.Printf("   Produit: %dx %s\n", recette.QuantiteProduit, recette.Produit.Nom)
		fmt.Printf("   Ingrédients requis:\n")
		for _, ingredient := range recette.Ingredients {
			fmt.Printf("     - %dx %s\n", ingredient.Quantite, ingredient.Item.Nom)
		}
	}

	// Pause avant retour
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// crafterObjet permet au joueur de crafter un objet s’il a les bons ingrédients
func crafterObjet(joueur *character.Character) {
	recettes := GetRecettesDisponibles()

	fmt.Println("\n⚒️  === CRÉATION D'OBJET === ⚒️")

	// Création de la liste d’options affichées au joueur
	options := make([]string, 0)
	for i, recette := range recettes {
		disponible := "✅"
		if !peutCrafter(joueur, recette) {
			disponible = "❌" // Marque la recette comme non craftable
		}
		options = append(options, fmt.Sprintf("%s %s (%dx %s)",
			disponible, recette.Nom, recette.QuantiteProduit, recette.Produit.Nom))
		_ = i // pour éviter une erreur si i est inutilisé
	}
	options = append(options, "Retour")

	// On affiche le menu des recettes
	ui.AfficherMenu("Choisir une recette à crafter", options)
	choix := utils.ScanChoice("Quelle recette voulez-vous utiliser ? ", options)

	if choix == len(options) {
		return // Retour
	}

	recetteChoisie := recettes[choix-1]

	// Vérification si le joueur a les ingrédients nécessaires
	if !peutCrafter(joueur, recetteChoisie) {
		fmt.Println("\n❌ Vous n'avez pas les ingrédients nécessaires pour cette recette !")
		fmt.Println("\nIngrédients requis :")
		for _, ingredient := range recetteChoisie.Ingredients {
			quantitePossedee := compterItem(joueur, ingredient.Item.Nom)
			fmt.Printf("  - %s : %d/%d %s\n",
				ingredient.Item.Nom,
				quantitePossedee,
				ingredient.Quantite,
				getStatusIcon(quantitePossedee >= ingredient.Quantite))
		}
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}

	// Confirmation avant le craft
	fmt.Printf("\n🔨 Crafter : %s\n", recetteChoisie.Nom)
	fmt.Printf("Produit : %dx %s\n", recetteChoisie.QuantiteProduit, recetteChoisie.Produit.Nom)

	options = []string{"Confirmer le craft", "Annuler"}
	ui.AfficherMenu("Confirmation", options)
	confirmation := utils.ScanChoice("Êtes-vous sûr ? ", options)

	if confirmation == 1 {
		// Si confirmé → on retire les ingrédients et on ajoute le produit
		retirerIngredients(joueur, recetteChoisie)
		ajouterProduit(joueur, recetteChoisie)

		fmt.Printf("\n✅ %s créé avec succès !\n", recetteChoisie.Nom)
		fmt.Printf("Vous avez reçu : %dx %s\n", recetteChoisie.QuantiteProduit, recetteChoisie.Produit.Nom)
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// peutCrafter vérifie si le joueur possède tous les ingrédients d’une recette
func peutCrafter(joueur *character.Character, recette Recette) bool {
	for _, ingredient := range recette.Ingredients {
		if compterItem(joueur, ingredient.Item.Nom) < ingredient.Quantite {
			return false // Pas assez d’un ingrédient
		}
	}
	return true // Tous les ingrédients sont présents
}

// compterItem compte combien d’exemplaires d’un objet le joueur possède
func compterItem(joueur *character.Character, nomItem string) int {
	count := 0
	for _, item := range joueur.Inventaire.Items {
		if item.Nom == nomItem {
			count++
		}
	}
	return count
}

// getStatusIcon renvoie une icône selon si un ingrédient est disponible ou non
func getStatusIcon(disponible bool) string {
	if disponible {
		return "✅"
	}
	return "❌"
}

// retirerIngredients enlève les ingrédients consommés du sac du joueur
func retirerIngredients(joueur *character.Character, recette Recette) {
	for _, ingredient := range recette.Ingredients {
		retirees := 0
		nouvelInventaire := make([]item.Item, 0)

		// On recrée un inventaire sans les ingrédients utilisés
		for _, item := range joueur.Inventaire.Items {
			if item.Nom == ingredient.Item.Nom && retirees < ingredient.Quantite {
				retirees++ // On consomme l’objet
				// On ne l’ajoute pas au nouvel inventaire
			} else {
				nouvelInventaire = append(nouvelInventaire, item)
			}
		}

		// On met à jour l’inventaire du joueur
		joueur.Inventaire.Items = nouvelInventaire
	}
}

// ajouterProduit ajoute le produit crafté à l’inventaire du joueur
func ajouterProduit(joueur *character.Character, recette Recette) {
	// Cas spéciaux : les potions vont dans des compteurs séparés
	if recette.Produit.Nom == "Potion de Vie" {
		joueur.Inventaire.Potions += recette.QuantiteProduit
		return
	}
	if recette.Produit.Nom == "Potion de Mana" {
		joueur.Inventaire.PotionsMana += recette.QuantiteProduit
		return
	}

	// Cas normal : on ajoute l’objet directement à la liste Items
	for i := 0; i < recette.QuantiteProduit; i++ {
		joueur.Inventaire.Items = append(joueur.Inventaire.Items, recette.Produit)
	}
}

```

7. Exploration.go : Gère la boucle de gameplay principale et les déplacements

```go

package exploration

import (
	"fmt"
	"strings"
	"world_of_milousques/banque"
	"world_of_milousques/character"
	"world_of_milousques/commerce"
	"world_of_milousques/craft"
	"world_of_milousques/fight"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
	"world_of_milousques/world"
)

// ExplorerMap lance la boucle principale d'exploration du monde ouvert pour le joueur
func ExplorerMap(joueur *character.Character) {
	// Création d'une nouvelle carte du monde
	gameMap := world.NewMap()
	
	// Restaurer la position et les zones découvertes du joueur depuis la sauvegarde
	x, y := joueur.ObtenirPosition()
	gameMap.RestaurerPosition(x, y)
	gameMap.RestaurerEtatDecouverte(joueur.ZonesDecouvertes)
	
	// Restaurer l'état des ressources de la carte selon ce que le joueur a déjà récolté
	gameMap.RestaurerEtatRessources(joueur)
	
	// Marquer la zone actuelle comme découverte
	joueur.MarquerZoneDecouverte(x, y)
	
	// Message d'accueil du joueur
	fmt.Println("\n🗺️  === BIENVENUE DANS LE MONDE OUVERT === 🗺️")
	fmt.Println("Vous pouvez maintenant explorer le monde librement !")
	fmt.Println("Utilisez les menus pour vous déplacer et interagir avec l'environnement.")
	
	// Afficher le nombre de zones découvertes
	nombreZones := joueur.ObtenirNombreZonesDecouvertes()
	fmt.Printf("Vous avez déjà découvert %d zones sur 25.\n", nombreZones)
	
	actionCount := 0          // Compteur global des actions
	maxActions := 1000        // Limite pour éviter les boucles infinies
	
	for actionCount < maxActions {
		actionCount++
		
		// Vérifier si le joueur est mort
		if joueur.Pdv <= 0 {
			fmt.Println("\n💀 Vous êtes mort ! Le jeu se termine.")
			break
		}
		
		// Afficher la carte
		gameMap.AfficherMap()
		
		// Afficher le menu principal d'exploration
		if !menuPrincipalExploration(gameMap, joueur) {
			break // Le joueur a choisi de quitter
		}
		
		// Sauvegarde automatique toutes les 50 actions
		if actionCount%50 == 0 {
			fmt.Printf("\n💾 Sauvegarde automatique... (Action %d/%d)\n", actionCount, maxActions)
			if err := joueur.Sauvegarder(); err != nil {
				fmt.Println("⚠️  Erreur lors de la sauvegarde automatique:", err)
			} else {
				fmt.Println("✅ Sauvegarde réussie !")
			}
		}
	}
	
	// Si la limite d'actions est atteinte
	if actionCount >= maxActions {
		fmt.Printf("\n⚠️  Limite d'actions atteinte (%d). Le jeu se ferme pour éviter une surcharge.\n", maxActions)
		fmt.Println("Votre progression a été sauvegardée automatiquement.")
	}
}

// menuPrincipalExploration affiche le menu principal avec les actions possibles du joueur
func menuPrincipalExploration(gameMap *world.Map, joueur *character.Character) bool {
	options := []string{
		"Explorer cette zone",
		"Se déplacer",
		"Voir la carte complète",
		"Afficher le statut du personnage",
		"Quitter le jeu",
	}
	
	// Afficher le menu et demander le choix
	ui.AfficherMenu("Que voulez-vous faire ?", options)
	choix := utils.ScanChoice("Votre choix : ", options)
	
	switch choix {
	case 1:
		explorerZoneActuelle(gameMap, joueur) // Explorer la zone actuelle
	case 2:
		seDeplacer(gameMap, joueur)           // Déplacer le joueur
	case 3:
		gameMap.AfficherMap()                 // Afficher la carte complète
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	case 4:
		afficherStatutPersonnage(joueur)      // Afficher les stats du personnage
	case 5:
		fmt.Println("Merci d'avoir joué à World of Milousques !")
		return false
	}
	
	return true
}

// explorerZoneActuelle ouvre le menu d'actions pour la zone où le joueur se trouve
func explorerZoneActuelle(gameMap *world.Map, joueur *character.Character) {
	zone := gameMap.GetCurrentZone() // Obtenir la zone actuelle
	
	fmt.Printf("\n🏠  === %s === 🏠\n", zone.Nom)
	fmt.Println(zone.Description)
	fmt.Println()
	
	zoneActionCount := 0
	maxZoneActions := 50 // Limite d'actions par zone
	
	for zoneActionCount < maxZoneActions {
		zoneActionCount++
		
		options := []string{}
		estAstrab := strings.Contains(zone.Nom, "Astrab") // Vérifier si on est dans la ville Astrab
		
		// Ajouter options selon contenu de la zone
		if len(zone.Ressources) > 0 {
			options = append(options, fmt.Sprintf("Récolter des ressources (%d disponibles)", len(zone.Ressources)))
		}
		if len(zone.Monstres) > 0 {
			options = append(options, fmt.Sprintf("Affronter un monstre (%d disponibles)", len(zone.Monstres)))
		}
		if len(zone.PNJs) > 0 {
			options = append(options, fmt.Sprintf("Parler aux habitants (%d présents)", len(zone.PNJs)))
		}
		
		// Options spéciales pour Astrab
		if estAstrab {
			options = append(options, "🔨 Aller à la forge")
			options = append(options, "💰 Aller chez le marchand")
			options = append(options, "🏦 Aller à la banque")
		}
		
		options = append(options, "Retour à la carte")
		
		// Sécurité : toujours au moins "Retour"
		if len(options) == 0 {
			fmt.Println("Erreur : Aucune option disponible. Retour automatique.")
			return
		}
		
		// Si la zone est vide
		if len(options) == 1 {
			fmt.Println("Cette zone semble vide... Il n'y a rien d'intéressant ici.")
			fmt.Println("Appuyez sur Entrée pour retourner à la carte.")
			fmt.Scanln()
			return
		}
		
		ui.AfficherMenu(fmt.Sprintf("Explorer %s", zone.Nom), options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		// Sécurité pour choix invalide
		if choix < 1 || choix > len(options) {
			fmt.Printf("Choix invalide (%d). Retour automatique.\n", choix)
			return
		}
		
		currentIndex := 0
		
		// Récolter des ressources
		if len(zone.Ressources) > 0 {
			currentIndex++
			if choix == currentIndex {
				recolterRessources(zone, joueur)
				continue
			}
		}
		
		// Affronter un monstre
		if len(zone.Monstres) > 0 {
			currentIndex++
			if choix == currentIndex {
				affronterMonstre(zone, joueur)
				if joueur.Pdv <= 0 { // Vérifier si le joueur est mort
					fmt.Println("\n💀 Vous avez été vaincu...")
					return
				}
				continue
			}
		}
		
		// Parler aux PNJs
		if len(zone.PNJs) > 0 {
			currentIndex++
			if choix == currentIndex {
				parlerAuxPNJs(zone, joueur)
				continue
			}
		}
		
		// Options spéciales pour Astrab
		if estAstrab {
			currentIndex++
			if choix == currentIndex {
				craft.AfficherForge(joueur)
				continue
			}
			currentIndex++
			if choix == currentIndex {
				commerce.AfficherMarchand(joueur)
				continue
			}
			currentIndex++
			if choix == currentIndex {
				banque.AfficherBanque(joueur)
				continue
			}
		}
		
		// Retour à la carte
		currentIndex++
		if choix == currentIndex {
			return
		}
		
		// Sécurité : erreur logique
		fmt.Printf("⚠️  Erreur de logique avec le choix %d. Retour automatique.\n", choix)
		return
	}
	
	// Limite d'actions dans la zone
	if zoneActionCount >= maxZoneActions {
		fmt.Printf("\n⚠️  Trop d'actions dans cette zone (%d). Retour automatique à la carte.\n", maxZoneActions)
		fmt.Println("Appuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
}

// seDeplacer gère le déplacement du joueur avec les touches ZQSD
func seDeplacer(gameMap *world.Map, joueur *character.Character) {
	fmt.Println("\nDéplacements possibles :")
	fmt.Println("Z = Nord | S = Sud | Q = Ouest | D = Est | A = Annuler")
	
	// Options disponibles selon déplacements possibles
	optionsDisponibles := []string{}
	if gameMap.CanMoveTo("NORD") {
		optionsDisponibles = append(optionsDisponibles, "Z (Nord)")
	}
	if gameMap.CanMoveTo("SUD") {
		optionsDisponibles = append(optionsDisponibles, "S (Sud)")
	}
	if gameMap.CanMoveTo("OUEST") {
		optionsDisponibles = append(optionsDisponibles, "Q (Ouest)")
	}
	if gameMap.CanMoveTo("EST") {
		optionsDisponibles = append(optionsDisponibles, "D (Est)")
	}
	optionsDisponibles = append(optionsDisponibles, "A (Annuler)")
	
	// Vérifier si déplacement impossible
	if len(optionsDisponibles) == 1 {
		fmt.Println("Vous ne pouvez pas vous déplacer d'ici !")
		return
	}
	
	ui.AfficherMenu("Choisir une direction", optionsDisponibles)
	choixInput := utils.ScanString("Tapez Z/Q/S/D pour vous déplacer (ou A pour annuler) : ", 1)
	choixInput = strings.ToUpper(strings.TrimSpace(choixInput))
	
	direction := ""
	nomDirection := ""
	
	// Vérification de la direction choisie et permission de déplacement
	switch choixInput {
	case "Z":
		if gameMap.CanMoveTo("NORD") {
			direction = "NORD"
			nomDirection = "Nord"
		} else {
			fmt.Println("Vous ne pouvez pas aller au Nord !")
			return
		}
	case "S":
		if gameMap.CanMoveTo("SUD") {
			direction = "SUD"
			nomDirection = "Sud"
		} else {
			fmt.Println("Vous ne pouvez pas aller au Sud !")
			return
		}
	case "Q":
		if gameMap.CanMoveTo("OUEST") {
			direction = "OUEST"
			nomDirection = "Ouest"
		} else {
			fmt.Println("Vous ne pouvez pas aller à l'Ouest !")
			return
		}
	case "D":
		if gameMap.CanMoveTo("EST") {
			direction = "EST"
			nomDirection = "Est"
		} else {
			fmt.Println("Vous ne pouvez pas aller à l'Est !")
			return
		}
	case "A":
		return // Annuler
	default:
		fmt.Println("Direction invalide ! Utilisez Z/Q/S/D ou A.")
		return
	}
	
	// Déplacement effectif et mise à jour de la zone
	if gameMap.MoveToWithCharacter(direction, joueur) {
		newZone := gameMap.GetCurrentZone()
		fmt.Printf("\n🚶 Vous vous déplacez vers le %s...\n", nomDirection)
		fmt.Printf("📍 Vous arrivez à : %s\n", newZone.Nom)
		
		nombreZones := joueur.ObtenirNombreZonesDecouvertes()
		fmt.Printf("🗺️  Zones découvertes : %d/25\n", nombreZones)
		
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

```
8. Fight.go : Le système de combat au tour par tour

```go

package fight

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/ui"
)

// Ennemi représente un monstre ou adversaire dans un combat
type Ennemi struct {
	Nom     string // Nom de l'ennemi
	Pv      int    // Points de vie
	Attaque int    // Points de dégâts infligés par attaque
}

// Fight lance un combat entre le joueur et un ennemi
func Fight(joueur *character.Character, ennemi Ennemi) {
	// Boucle principale du combat : continue tant que le joueur et l'ennemi ont des PV
	for joueur.Pdv > 0 && ennemi.Pv > 0 {

		// Affichage de l'état du combat
		ui.AfficherMenuCombat(
			joueur.Nom, joueur.Pdv, joueur.Classe.Pvmax, joueur.Mana, joueur.Classe.ManaMax,
			ennemi.Nom, ennemi.Pv, joueur.Classe.Sorts, joueur.Inventaire.Potions,
		)

		// Demander l'action du joueur
		fmt.Print("Choisis ton action : ")
		var choix int
		fmt.Scanln(&choix)

		// L'option potion est toujours après les sorts
		optionPotion := len(joueur.Classe.Sorts) + 1

		if choix == optionPotion { // Le joueur choisit d'utiliser une potion
			if joueur.Inventaire.Potions > 0 {
				joueur.Pdv += 50 // Restaurer 50 PV
				if joueur.Pdv > joueur.Classe.Pvmax {
					joueur.Pdv = joueur.Classe.Pvmax // Limiter aux PV max
				}
				joueur.Inventaire.Potions-- // Décrémenter potion
				fmt.Println("Tu utilises une potion et récupères 50 PVs !")
			} else {
				fmt.Println("Tu n'as pas de potion !")
				continue // Recommencer le tour sans que l'ennemi attaque
			}
		} else if choix >= 1 && choix <= len(joueur.Classe.Sorts) { // Le joueur choisit un sort
			s := joueur.Classe.Sorts[choix-1] // Récupérer le sort choisi
			if joueur.Mana < s.Cout {
				fmt.Println("Pas assez de mana pour lancer ce sort !")
				continue // Recommencer le tour
			}
			joueur.Mana -= s.Cout        // Dépenser la mana du sort
			ennemi.Pv -= s.Degats        // Infliger les dégâts à l'ennemi
			fmt.Printf("Tu lances %s et infliges %d dégâts !\n", s.Nom, s.Degats)
		} else {
			fmt.Println("Choix invalide, réessaie.") // Entrée non valide
			continue
		}

		// Vérifier si l'ennemi est mort après l'attaque du joueur
		if ennemi.Pv <= 0 {
			fmt.Printf("%s est vaincu !\n", ennemi.Nom)
			break // Sortir de la boucle, le combat est terminé
		}

		// Tour de l'ennemi : infliger des dégâts au joueur
		joueur.Pdv -= ennemi.Attaque
		fmt.Printf("%s t'attaque et inflige %d dégâts !\n", ennemi.Nom, ennemi.Attaque)
	}

	// Combat terminé : vérifier le résultat
	if joueur.Pdv > 0 {
		// Restaurer la mana du joueur après victoire
		joueur.Mana = joueur.Classe.ManaMax
		fmt.Println("Ta mana est maintenant restaurée au maximum !")
	} else {
		// Le joueur est mort
		fmt.Println("Tu as été vaincu... Game Over.")
	}
}

```

9. Inventory.go : Tout ce qui nous faut pour gérer l'inventaire

```go

package inventory

import (
	"fmt"
	"world_of_milousques/item"
)

// Inventaire représente l'inventaire du joueur
type Inventaire struct {
	Potions int         `json:"potions"` // Nombre de potions disponibles
	Items   []item.Item `json:"items"`   // Liste des objets possédés
}

// AddItem ajoute un objet à l'inventaire, plusieurs fois si nécessaire
func (inv *Inventaire) AddItem(it item.Item, quantity int) {
	for i := 0; i < quantity; i++ {
		inv.Items = append(inv.Items, it) // Ajout de l'objet à la slice Items
	}
}

// Recolter ajoute une liste de ressources récoltées dans l'inventaire
func (inv *Inventaire) Recolter(ressources []item.Item) {
	if len(ressources) == 0 {
		fmt.Println("Aucune ressource à récolter ici.")
		return
	}

	fmt.Println("Vous récoltez :")
	for _, it := range ressources {
		fmt.Printf("- %s\n", it.Nom) // Affichage du nom de chaque ressource
		inv.AddItem(it, 1)           // Ajout de l'objet à l'inventaire
	}
	fmt.Printf("Votre inventaire contient maintenant %d objets.\n", len(inv.Items))
}

// Afficher imprime le contenu complet de l'inventaire
func (inv *Inventaire) Afficher() {
	if len(inv.Items) == 0 {
		fmt.Println("Votre inventaire est vide.")
		return
	}
	fmt.Println("Inventaire :")
	for i, it := range inv.Items {
		fmt.Printf("%d) %s | Poids: %d | Effet: %s | Valeur: %d\n", 
			i+1, it.Nom, it.Poids, it.Effet, it.Valeur) // Détails de chaque objet
	}
	fmt.Printf("Total d'objets : %d\n", len(inv.Items))
}

```
10. Item.go : on définit simplement les items

```go

package item

// Item représente un objet dans le jeu
type Item struct {
	Nom    string // Nom de l'objet
	Poids  int    // Poids de l'objet
	Effet  string // Description ou effet de l'objet
	Valeur int    // Valeur de l'objet en pièces
}

// NewItem crée un nouvel objet selon son nom
func NewItem(nom string) Item {
	switch nom {
	case "Bois":
		return Item{
			Nom: "Bois", 
			Poids: 10, 
			Effet: "Manger du bois vous fera mal aux dents", 
			Valeur: 5,
		}
	case "Fer":
		return Item{
			Nom: "Fer", 
			Poids: 10, 
			Effet: "Pas le matériau le plus adéquat pour fabriquer un lit", 
			Valeur: 10,
		}
	case "Blé":
		return Item{
			Nom: "Blé", 
			Poids: 2, 
			Effet: "Votre meilleur ami pour être accepté à Ynuv", 
			Valeur: 1,
		}
	case "Laitue Vireuse":
		return Item{
			Nom: "Laitue Vireuse", 
			Poids: 1, 
			Effet: "La plante de secours du Grand Alchimiste Yelram Bob", 
			Valeur: 1,
		}
	case "Pichon":
		return Item{
			Nom: "Pichon", 
			Poids: 2, 
			Effet: "Piche qui glisse n'amasse pas de risques!", 
			Valeur: 2,
		}
	default:
		// Pour tout autre nom, renvoie un objet générique avec un effet d'erreur
		return Item{
			Nom: nom, 
			Poids: 10, 
			Effet: "oops erreur", 
			Valeur: 10,
		}
	}
}

```
11. Places.go : Gère notre tutoriel (oui le nom est giga mauvais mais j'avais peur de tout casser en le renommant et j'avais pas le temps de réparer les erreurs qui aller venir)

```go

package places

import (
	"fmt" // Pour afficher des messages à l'écran

	"world_of_milousques/character" // Pour manipuler les personnages
	"world_of_milousques/fight"     // Pour gérer les combats
	"world_of_milousques/item"      // Pour gérer les objets
	"world_of_milousques/ui"        // Pour l'affichage des menus
	"world_of_milousques/utils"     // Pour les entrées utilisateur
)

// Scene représente une scène du jeu avec dialogue et options
type Scene struct {
	Titre       string                   // Titre de la scène
	Description string                   // Texte ou description affichée au joueur
	Options     []string                 // Choix possibles pour le joueur
	Actions     []func(*character.Character) // Fonctions exécutées selon le choix
	Ressources  []item.Item              // Objets disponibles dans la scène
}

// GetIntroDialogue retourne le dialogue d'introduction du jeu
func GetIntroDialogue() []Scene {
	return []Scene{
		{
			Titre: "Réveil mystérieux", // Nom de la scène
			Description: "??? : Réveille toi, aventurier...", // Message initial
			Options: []string{"Qui êtes-vous ?", "Où suis-je ?", "Que sont les Milousques ?", "On peut y aller !"}, // Choix du joueur
			Actions: []func(*character.Character){ // Actions correspondant aux options
				// Option 1 : Qui êtes-vous ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Je suis Mathiouw, Le berger des jeunes âmes. Mon but est de faire de toi un aventurier assez puissant pour partir en quête des milousques.")
				},
				// Option 2 : Où suis-je ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Tu es à Astrab, le lieu d'apparition des chasseurs de milousques comme toi.")
				},
				// Option 3 : Que sont les Milousques ?
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Les milousques sont de puissantes chimères qui donne à ceux capable de les dompter un pouvoir incommensurable !")
				},
				// Option 4 : On peut y aller !
				func(c *character.Character) {
					fmt.Println("\nMathiouw : Parfait ! Commençons par un petit test de tes capacités...")
					// Cette action sera gérée dans main.go pour passer à la suite
				},
			},
			Ressources: []item.Item{}, // Pas de ressources dans cette scène
		},
	}
}

// GetTutorielCombat retourne les informations du tutoriel de combat
func GetTutorielCombat() (string, string, *fight.Ennemi) {
	quete := "Vaincre le Chacha Agressif" // Nom de la quête
	recompense := "1 potion"               // Récompense pour le joueur
	ennemi := &fight.Ennemi{               // Création de l'ennemi
		Nom: "Chacha Agressif",
		Pv: 50,
		Attaque: 15,
	}
	return quete, recompense, ennemi // Retourne les infos
}

// ProposerQueteTutoriel propose la quête du tutoriel avec option de refus
func ProposerQueteTutoriel(c *character.Character) bool {
	fmt.Println("\n🎒 === QUÊTE PROPOSÉE === 🎒")
	fmt.Println("Mathiouw : Alors, acceptes-tu de m'aider à vaincre le Chacha Agressif ?")
	fmt.Println("Récompense : 1 potion de soin")

	options := []string{"Accepter la quête", "Refuser la quête"} // Les choix du joueur
	ui.AfficherMenu("Décision", options)                          // Affiche le menu
	choix := utils.ScanChoice("Votre décision : ", options)      // Récupère le choix

	if choix == 1 { // Si le joueur accepte
		fmt.Println("\nMathiouw : Parfait ! Je savais que je pouvais compter sur toi.")
		c.ProposerEtAjouterQueteAvecPNJ("Vaincre le Chacha Agressif", "1 potion", "Mathiouw") // Ajoute la quête au personnage
		return true
	} else { // Si le joueur refuse
		fmt.Println("\nMathiouw : Le chacha a pris la quête à ta place, si tu le bats il empochera la récompense de sa défaite. Qu'il est malin ce chacha !")
		fmt.Println("\n(La quête continue quand même pour le tutoriel)")
		return false
	}
}

```

12. Sorts.go : Tout ce qui est en rapports avec les sorts

```go

package sorts // Déclare le package "sorts", utilisé pour gérer les sorts du jeu

// Sorts représente un sort avec son nom, ses dégâts et son coût en mana
type Sorts struct {
	Nom    string // Nom du sort
	Degats int    // Dégâts que le sort inflige
	Cout   int    // Coût en mana pour lancer le sort
}

// InitSorts initialise un nouveau sort avec les valeurs fournies
func InitSorts(nom string, degats int, cout int) Sorts {
	return Sorts{
		Nom:    nom,    // Nom du sort
		Degats: degats, // Dégâts infligés
		Cout:   cout,   // Coût en mana
	}
}

// GetSorts retourne un sort prédéfini selon son nom
func GetSorts(nom string) Sorts {
	switch nom {
	case "Boule de feu": // Sort offensif classique
		return Sorts{Nom: "Boule de feu", Degats: 30, Cout: 20}
	case "Explosion": // Sort puissant avec gros dégâts
		return Sorts{Nom: "Explosion", Degats: 50, Cout: 40}
	case "Coup bas": // Sort moins coûteux mais efficace
		return Sorts{Nom: "Coup bas", Degats: 25, Cout: 15}
	case "Fourberie": // Sort tactique, faible dégât mais gratuit
		return Sorts{Nom: "Fourberie", Degats: 10, Cout: 0}
	case "Fracasser": // Sort moyen, coût raisonnable
		return Sorts{Nom: "Fracasser", Degats: 20, Cout: 10}
	case "Briser": // Sort puissant mais coûteux
		return Sorts{Nom: "Briser", Degats: 40, Cout: 20}
	default: // Sort par défaut si le nom n'est pas reconnu
		return Sorts{Nom: nom, Degats: 25, Cout: 15}
	}
}

```

13. Ui.go : Gère l'interface en ASCII Art avec caractères unicode 

```go

package ui // Déclare le package "ui" pour gérer l'affichage du jeu

import (
	"fmt"       // Pour afficher du texte dans le terminal
	"strings"   // Pour manipuler les chaînes de caractères
	"world_of_milousques/sorts" // Pour accéder aux sorts du jeu
)

// calculerLargeurAffichage calcule la largeur réelle d'affichage d'une chaîne en tenant compte des caractères Unicode (emojis) qui prennent plus d'espace
func calculerLargeurAffichage(s string) int {
	largeur := 0
	for _, r := range s { // Boucle sur chaque rune (caractère Unicode) de la chaîne
		if estEmoji(r) { // Si c'est un emoji, compter comme 2 colonnes
			largeur += 2
		} else { // Sinon, compter comme 1 colonne
			largeur += 1
		}
	}
	return largeur
}

// estEmoji détermine si un caractère (rune) est un emoji
func estEmoji(r rune) bool {
	// Plages Unicode communes pour les emojis
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) ||   // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) ||   // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) ||   // Indicateurs régionaux
		(r >= 0x2600 && r <= 0x26FF) ||     // Symboles divers
		(r >= 0x2700 && r <= 0x27BF) ||     // Dingbats
		(r >= 0xFE0F && r <= 0xFE0F)       // Variation Selector
}

// centrerTexte centre un texte dans une largeur donnée
func centrerTexte(texte string, largeur int) string {
	w := calculerLargeurAffichage(texte) // Largeur réelle du texte
	if w >= largeur { // Si le texte est plus large que la largeur donnée, renvoyer tel quel
		return texte
	}
	paddingTotal := largeur - w              // Calculer le nombre total d'espaces à ajouter
	paddingGauche := paddingTotal / 2       // Répartir moitié à gauche
	paddingDroite := paddingTotal - paddingGauche // Le reste à droite
	return strings.Repeat(" ", paddingGauche) + texte + strings.Repeat(" ", paddingDroite)
}

// alignerGauche aligne un texte à gauche dans une largeur donnée
func alignerGauche(texte string, largeur int) string {
	w := calculerLargeurAffichage(texte)
	if w >= largeur { // Si le texte est plus large, le renvoyer tel quel
		return texte
	}
	return texte + strings.Repeat(" ", largeur-w) // Ajouter des espaces à droite
}

// AfficherMenuSimple affiche un menu simple avec bordures Unicode
func AfficherMenuSimple(titre string, options []string) {
	largeurContenu := calculerLargeurAffichage(titre) // Déterminer largeur initiale
	for i, opt := range options { // Boucler sur toutes les options pour trouver la largeur max
		ligne := fmt.Sprintf(" %d) %s", i+1, opt) // Préfixe option avec numéro
		lw := calculerLargeurAffichage(ligne)
		if lw > largeurContenu {
			largeurContenu = lw
		}
	}

	// Limites de largeur pour éviter le wrapping dans les terminaux
	if largeurContenu < 30 {
		largeurContenu = 30
	}
	if largeurContenu > 50 {
		largeurContenu = 50
	}

	largeurTotale := largeurContenu + 4 // Ajouter marge (2 espaces de chaque côté)

	// Ligne supérieure avec bordure Unicode
	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)

	// Titre centré
	titreCentre := centrerTexte(titre, largeurTotale)
	ligneTitre := "\u2502" + titreCentre + "\u2502"
	fmt.Println(ligneTitre)

	// Ligne de séparation
	ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
	fmt.Println(ligneSeparation)

	// Options du menu
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt)
		ligneAlignee := alignerGauche(ligne, largeurTotale)
		ligneOption := "\u2502" + ligneAlignee + "\u2502"
		fmt.Println(ligneOption)
	}

	// Ligne inférieure
	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}

// AfficherMenu similaire à AfficherMenuSimple mais version plus générique
func AfficherMenu(titre string, options []string) {
	largeurContenu := calculerLargeurAffichage(titre)
	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt)
		if calculerLargeurAffichage(ligne) > largeurContenu {
			largeurContenu = calculerLargeurAffichage(ligne)
		}
	}

	if largeurContenu < 30 {
		largeurContenu = 30
	}
	if largeurContenu > 50 {
		largeurContenu = 50
	}

	largeurTotale := largeurContenu + 4

	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)

	titreCentre := centrerTexte(titre, largeurTotale)
	ligneTitre := "\u2502" + titreCentre + "\u2502"
	fmt.Println(ligneTitre)

	ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
	fmt.Println(ligneSeparation)

	for i, opt := range options {
		ligne := fmt.Sprintf(" %d) %s", i+1, opt)
		ligneAlignee := alignerGauche(ligne, largeurTotale)
		ligneOption := "\u2502" + ligneAlignee + "\u2502"
		fmt.Println(ligneOption)
	}

	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}

// AfficherMenuCombat affiche un menu spécifique au combat
func AfficherMenuCombat(joueurNom string, joueurPv, joueurPvMax, joueurMana, joueurManaMax int,
	ennemiNom string, ennemiPv int, sortsList []sorts.Sorts, potions, potionsMana int) {

	lignes := []string{} // Slice pour stocker toutes les lignes à afficher

	// Ajouter les informations du joueur
	lignes = append(lignes, fmt.Sprintf("%s : PV %d/%d | Mana %d/%d", joueurNom, joueurPv, joueurPvMax, joueurMana, joueurManaMax))
	// Ajouter les informations de l'ennemi
	lignes = append(lignes, fmt.Sprintf("Ennemi %s : PV %d", ennemiNom, ennemiPv))
	lignes = append(lignes, "") // Ligne vide pour séparation

	// Ajouter les sorts disponibles
	for i, s := range sortsList {
		lignes = append(lignes, fmt.Sprintf("%d) %s - Dégâts: %d, Mana: %d", i+1, s.Nom, s.Degats, s.Cout))
	}

	// Ajouter options potions et fuite
	lignes = append(lignes, fmt.Sprintf("%d) Utiliser une potion de vie (+50 PV) [%d disponibles]", len(sortsList)+1, potions))
	lignes = append(lignes, fmt.Sprintf("%d) Utiliser une potion de mana (+50 Mana) [%d disponibles]", len(sortsList)+2, potionsMana))
	lignes = append(lignes, fmt.Sprintf("%d) Fuir le combat", len(sortsList)+3))

	// Calculer largeur maximale pour alignement
	largeurContenu := 0
	for _, ligne := range lignes {
		if calculerLargeurAffichage(ligne) > largeurContenu {
			largeurContenu = calculerLargeurAffichage(ligne)
		}
	}

	largeurTotale := largeurContenu + 4 // Ajouter marge

	// Ligne supérieure
	ligneBordure := "\u250C" + strings.Repeat("\u2500", largeurTotale) + "\u2510"
	fmt.Println(ligneBordure)

	// Afficher chaque ligne avec bordures
	for _, ligne := range lignes {
		if ligne == "" { // Ligne vide = ligne de séparation
			ligneSeparation := "\u251C" + strings.Repeat("\u2500", largeurTotale) + "\u2524"
			fmt.Println(ligneSeparation)
		} else { // Ligne de contenu normale
			ligneAlignee := alignerGauche(ligne, largeurTotale)
			ligneOption := "\u2502" + ligneAlignee + "\u2502"
			fmt.Println(ligneOption)
		}
	}

	// Ligne inférieure
	ligneBordureInf := "\u2514" + strings.Repeat("\u2500", largeurTotale) + "\u2518"
	fmt.Println(ligneBordureInf)
}

```

14. Utils.go : Quelques fonctions utiles pour que l'on réapelle ailleurs (un peu un fourre-tout de fonctions que j'arriver pas a ranger)

```go

package utils // Déclare le package "utils" pour regrouper des fonctions utilitaires

import (
	"bufio"    // Pour lire l'entrée utilisateur de manière plus contrôlée
	"fmt"      // Pour afficher du texte dans le terminal
	"os"       // Pour accéder à l'entrée standard
	"strconv"  // Pour convertir des chaînes en nombres
	"strings"  // Pour manipuler et nettoyer des chaînes
)

// ScanInt lit un entier depuis l'entrée standard avec validation et gestion d'erreur
func ScanInt(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin) // Crée un lecteur pour l'entrée utilisateur
	for {
		fmt.Print(prompt) // Affiche le message à l'utilisateur
		input, err := reader.ReadString('\n') // Lit la ligne entière entrée par l'utilisateur
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Réessayez.\n", err)
			continue // Recommencer la lecture si erreur
		}

		input = strings.TrimSpace(input) // Supprime espaces et retour chariot
		if input == "" {                 // Si l'utilisateur n'a rien entré
			fmt.Println("Veuillez entrer une valeur.")
			continue
		}

		value, err := strconv.Atoi(input) // Tenter de convertir la chaîne en entier
		if err != nil {
			fmt.Printf("'%s' n'est pas un nombre valide. Réessayez.\n", input)
			continue
		}

		if value < min || value > max { // Vérifier que le nombre est dans les bornes
			fmt.Printf("Veuillez entrer un nombre entre %d et %d.\n", min, max)
			continue
		}

		return value // Retourner le nombre valide
	}
}

// ScanString lit une chaîne de caractères depuis l'entrée standard avec validation
func ScanString(prompt string, minLength int) string {
	reader := bufio.NewReader(os.Stdin) // Crée un lecteur pour l'entrée utilisateur
	for {
		fmt.Print(prompt) // Affiche le message
		input, err := reader.ReadString('\n') // Lit la ligne complète
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Réessayez.\n", err)
			continue
		}

		input = strings.TrimSpace(input) // Supprime espaces et retour chariot
		if len(input) < minLength {      // Vérifie que la chaîne est assez longue
			if minLength == 1 {
				fmt.Println("Veuillez entrer au moins un caractère.")
			} else {
				fmt.Printf("Veuillez entrer au moins %d caractères.\n", minLength)
			}
			continue
		}

		return input // Retourner la chaîne valide
	}
}

// ScanChoice lit un choix parmi des options avec validation flexible (numéro ou texte)
func ScanChoice(prompt string, options []string) int {
	// Vérification de sécurité : il doit y avoir au moins une option
	if len(options) == 0 {
		fmt.Println("Erreur : Aucune option disponible.")
		return 1 // Retourne 1 par défaut pour éviter erreurs d'index
	}

	reader := bufio.NewReader(os.Stdin) // Crée le lecteur
	attemptsCount := 0                  // Compteur de tentatives
	maxAttempts := 5                     // Nombre maximal de tentatives avant sélection automatique

	for {
		attemptsCount++ // Incrémenter le compteur de tentatives

		if attemptsCount > maxAttempts { // Si trop de tentatives
			fmt.Printf("\n⚠️  Trop de tentatives invalides (%d). Sélection automatique de l'option 1.\n", maxAttempts)
			fmt.Println("Appuyez sur Entrée pour continuer...")
			fmt.Scanln() // Pause pour que l'utilisateur voie le message
			return 1     // Retourne la première option par défaut
		}

		fmt.Print(prompt)
		input, err := reader.ReadString('\n') // Lire la saisie
		if err != nil {
			fmt.Printf("Erreur de lecture : %v. Tentative %d/%d.\n", err, attemptsCount, maxAttempts)
			continue
		}

		input = strings.TrimSpace(input) // Nettoyer la saisie
		if input == "" {                 // Vérifier qu'il y a quelque chose
			fmt.Printf("Veuillez faire un choix. Tentative %d/%d.\n", attemptsCount, maxAttempts)
			continue
		}

		// Essayer d'abord d'interpréter comme un numéro
		if value, err := strconv.Atoi(input); err == nil {
			if value >= 1 && value <= len(options) { // Vérifier que le numéro est valide
				return value
			} else {
				fmt.Printf("⚠️  Numéro hors limites (%d). Choisissez entre 1 et %d. Tentative %d/%d.\n", value, len(options), attemptsCount, maxAttempts)
				continue
			}
		}

		// Essayer comme correspondance textuelle (insensible à la casse)
		inputUpper := strings.ToUpper(input)
		for i, option := range options {
			if strings.Contains(strings.ToUpper(option), inputUpper) { // Si l'entrée correspond au début d'une option
				return i + 1 // Retourner l'index correspondant
			}
		}

		// Aucun choix valide
		fmt.Printf("⚠️  Choix invalide '%s'. Entrez un numéro (1-%d) ou le début d'une option. Tentative %d/%d.\n", input, len(options), attemptsCount, maxAttempts)
	}
}

```

15. World.go : Gère la map, les sous zones et tout ce qui trouve dedans

```go

// Package world gère la carte du jeu, les zones, les PNJs et la génération de contenu
// Structure la carte 5x5 avec différents biomes : champs, forêts, mines, rivières et routes
package world

import (
	"fmt" // Pour l'affichage de la carte
	"world_of_milousques/character" // Pour gérer les personnages et l'état des zones
	"world_of_milousques/fight"     // Pour gérer les ennemis
	"world_of_milousques/item"      // Pour gérer les objets et ressources
)

// PNJ représente un personnage non-joueur
type PNJ struct {
	Nom        string // Nom du PNJ
	Dialogue   string // Dialogue affiché au joueur
	Quete      string // Quête associée au PNJ
	Recompense string // Récompense donnée par le PNJ
}

// Zone représente une sous-zone de la map
type Zone struct {
	Nom         string         // Nom de la zone
	Description string         // Description de la zone
	Ressources  []item.Item    // Liste des ressources disponibles
	Monstres    []fight.Ennemi // Liste des ennemis dans la zone
	PNJs        []PNJ          // Liste des PNJs présents
	Visitee     bool           // Indique si le joueur a visité cette zone
}

// Position représente la position du joueur sur la map
type Position struct {
	X, Y int // Coordonnées X et Y
}

// Map représente la grille 5x5 du monde
type Map struct {
	Zones    [5][5]Zone // Grille de zones 5x5
	Position Position   // Position actuelle du joueur
}

// NewMap crée une nouvelle map et initialise toutes les zones
func NewMap() *Map {
	m := &Map{
		Position: Position{X: 2, Y: 2}, // Position centrale au départ
	}

	// Remplir toutes les zones avec du contenu
	m.initializeZones()

	return m
}

// RestaurerPosition met à jour la position du joueur sur la map
func (m *Map) RestaurerPosition(x, y int) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		m.Position = Position{X: x, Y: y} // Met à jour la position seulement si elle est valide
	}
}

// RestaurerEtatDecouverte met à jour l'état de découverte des zones
func (m *Map) RestaurerEtatDecouverte(zonesDecouvertes [5][5]bool) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if zonesDecouvertes[y][x] {
				m.Zones[y][x].Visitee = true // Marque la zone comme visitée si le joueur l'a découverte
			}
		}
	}
}

// RestaurerEtatRessources met à jour les ressources et monstres d'une zone selon un personnage
func (m *Map) RestaurerEtatRessources(characterInterface interface{}) {
	// Si c'est un personnage moderne
	if char, ok := characterInterface.(*character.Character); ok {
		m.restaurerAvecCharacter(char)
		return
	}

	// Sinon fallback pour l'ancienne méthode
	if charInterface, ok := characterInterface.(interface {
		ZoneRessourcesRecoltees(int, int) bool
		ZoneMonstresVaincus(int, int) bool
	}); ok {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				if charInterface.ZoneRessourcesRecoltees(x, y) {
					m.Zones[y][x].Ressources = []item.Item{} // Ressources vides si déjà récoltées
				}
				if charInterface.ZoneMonstresVaincus(x, y) {
					m.Zones[y][x].Monstres = []fight.Ennemi{} // Monstres vides si déjà vaincus
				}
			}
		}
	}
}

// restaurerAvecCharacter restaure les ressources et monstres avec l'état détaillé
func (m *Map) restaurerAvecCharacter(char *character.Character) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if etat, existe := char.ObtenirEtatZone(x, y); existe && etat.Visitee {
				m.restaurerRessourcesZone(&m.Zones[y][x], etat.RessourcesRestantes)
				m.restaurerMonstresZone(&m.Zones[y][x], etat.MonstresRestants)
			}
		}
	}
}

// restaurerRessourcesZone restaure les ressources d'une zone
func (m *Map) restaurerRessourcesZone(zone *Zone, ressourcesRestantes []string) {
	nouvellesRessources := []item.Item{}
	for _, nomRessource := range ressourcesRestantes {
		nouvellesRessources = append(nouvellesRessources, item.NewItem(nomRessource))
	}
	zone.Ressources = nouvellesRessources
}

// restaurerMonstresZone restaure les monstres d'une zone
func (m *Map) restaurerMonstresZone(zone *Zone, monstresRestants []character.MonstreState) {
	nouveauxMonstres := []fight.Ennemi{}
	for _, monstreState := range monstresRestants {
		nouveauxMonstres = append(nouveauxMonstres, fight.Ennemi{
			Nom:     monstreState.Nom,
			Pv:      monstreState.Pv,
			Attaque: monstreState.Attaque,
		})
	}
	zone.Monstres = nouveauxMonstres
}

// GetCurrentZone retourne la zone actuelle du joueur
func (m *Map) GetCurrentZone() *Zone {
	return &m.Zones[m.Position.Y][m.Position.X]
}

// GetZoneAt retourne la zone à une position donnée (utilisé pour tests)
func (m *Map) GetZoneAt(x, y int) *Zone {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return &m.Zones[y][x]
	}
	return nil
}

// CanMoveTo vérifie si le joueur peut se déplacer dans une direction
func (m *Map) CanMoveTo(direction string) bool {
	newX, newY := m.Position.X, m.Position.Y

	switch direction {
	case "NORD":
		newY--
	case "SUD":
		newY++
	case "OUEST":
		newX--
	case "EST":
		newX++
	default:
		return false // Direction invalide
	}

	return newX >= 0 && newX < 5 && newY >= 0 && newY < 5
}

// MoveTo déplace le joueur dans une direction sans personnage
func (m *Map) MoveTo(direction string) bool {
	return m.MoveToWithCharacter(direction, nil)
}

// MoveToWithCharacter déplace le joueur et met à jour la sauvegarde
func (m *Map) MoveToWithCharacter(direction string, character interface{}) bool {
	if !m.CanMoveTo(direction) {
		return false
	}

	switch direction {
	case "NORD":
		m.Position.Y--
	case "SUD":
		m.Position.Y++
	case "OUEST":
		m.Position.X--
	case "EST":
		m.Position.X++
	}

	m.GetCurrentZone().Visitee = true // Marquer la zone comme visitée

	// Sauvegarder la position si personnage fourni
	if char, ok := character.(interface {
		SauvegarderPositionMap(int, int)
		MarquerZoneDecouverte(int, int)
	}); ok {
		char.SauvegarderPositionMap(m.Position.X, m.Position.Y)
		char.MarquerZoneDecouverte(m.Position.X, m.Position.Y)
	}

	return true
}

// AfficherMap affiche la carte ASCII avec la position du joueur
func (m *Map) AfficherMap() {
	fmt.Println("\n=== CARTE DU MONDE ===\n")

	for y := 0; y < 5; y++ {
		// Ligne supérieure de chaque rangée
		for x := 0; x < 5; x++ {
			fmt.Print("+-------")
		}
		fmt.Println("+")
		
		// Ligne centrale avec symbole de la zone
		for x := 0; x < 5; x++ {
			zone := &m.Zones[y][x]
			symbol := " "
			if x == m.Position.X && y == m.Position.Y {
				symbol = "♦" // Joueur
			} else if zone.Visitee {
				symbol = "○" // Zone visitée
			} else {
				symbol = "?" // Zone inconnue
			}
			fmt.Printf("|   %s   ", symbol)
		}
		fmt.Println("|")
	}

	// Ligne du bas
	for x := 0; x < 5; x++ {
		fmt.Print("+-------")
	}
	fmt.Println("+")

	// Légende et position actuelle
	fmt.Println("\nLégende: ♦ = Vous | ○ = Visitée | ? = Inconnue")
	fmt.Printf("Position actuelle: %s (%d,%d)\n",
		m.GetCurrentZone().Nom, m.Position.X+1, m.Position.Y+1)
}

// initializeZones remplit la map avec du contenu cohérent géographiquement
func (m *Map) initializeZones() {
	zoneTypes := [5][5]string{
		{"Champs", "Champs", "Transition Nord", "Mines", "Mines"},
		{"Champs", "Forêt", "Transition Ouest", "Transition Est", "Mines"},
		{"Forêt", "Forêt", "Astrab", "Rivière", "Rivière"},
		{"Forêt", "Forêt", "Transition Sud", "Rivière", "Rivière"},
		{"Forêt", "Forêt", "Rivière", "Rivière", "Rivière"},
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			zone := &m.Zones[y][x]

			// Déterminer le type de zone spéciale
			isRoute := (x == 2 && (y == 0 || y == 1 || y == 3 || y == 4)) || (y == 2 && (x == 0 || x == 1 || x == 3 || x == 4))
			isSpecialChamps := (x == 4 && (y == 0 || y == 1)) || (x == 3 && (y == 0 || y == 1))
			isSpecialMine := (x == 0 && (y == 0 || y == 1)) || (x == 1 && (y == 0 || y == 1))
			isSpecialForet := (x == 0 && (y == 3 || y == 4)) || (x == 1 && (y == 3 || y == 4))
			isSpecialRiviere := (x == 3 && y == 3) || (x == 4 && y == 3) || (x == 3 && y == 4) || (x == 4 && y == 4)

			if isRoute {
				// Configuration pour les routes
				zone.Nom = "Route"
				zone.Description = "A la croisée des chemins, on trouve tous les gros malins !"
				zone.Ressources =
