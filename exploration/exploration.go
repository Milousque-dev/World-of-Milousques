// Package exploration gère la boucle principale du jeu en monde ouvert
// Permet au joueur de se déplacer, explorer, récolter et interagir avec l'environnement
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

// ExplorerMap lance la boucle principale d'exploration
func ExplorerMap(joueur *character.Character) {
	gameMap := world.NewMap()
	
	// Restaurer la position et l'état de découverte du joueur
	x, y := joueur.ObtenirPosition()
	gameMap.RestaurerPosition(x, y)
	gameMap.RestaurerEtatDecouverte(joueur.ZonesDecouvertes)
	
	// Restaurer l'état des ressources récoltées depuis le personnage
	gameMap.RestaurerEtatRessources(joueur)
	
	// Marquer la zone actuelle comme découverte
	joueur.MarquerZoneDecouverte(x, y)
	
	fmt.Println("\n🗺️  === BIENVENUE DANS LE MONDE OUVERT === 🗺️")
	fmt.Println("Vous pouvez maintenant explorer le monde librement !")
	fmt.Println("Utilisez les menus pour vous déplacer et interagir avec l'environnement.")
	
	// Afficher le nombre de zones découvertes
	nombreZones := joueur.ObtenirNombreZonesDecouvertes()
	fmt.Printf("Vous avez déjà découvert %d zones sur 25.\n", nombreZones)
	
	actionCount := 0
	maxActions := 1000 // Limite le nombre d'actions pour éviter les boucles infinies
	
	for actionCount < maxActions {
		actionCount++
		
		// Vérifier si le joueur est mort
		if joueur.Pdv <= 0 {
			fmt.Println("\n💀 Vous êtes mort ! Le jeu se termine.")
			break
		}
		
		// Afficher la map
		gameMap.AfficherMap()
		
		// Afficher le menu principal d'exploration
		if !menuPrincipalExploration(gameMap, joueur) {
			break // Le joueur veut quitter
		}
		
		// Sauvegarder automatiquement tous les 50 actions
		if actionCount%50 == 0 {
			fmt.Printf("\n💾 Sauvegarde automatique... (Action %d/%d)\n", actionCount, maxActions)
			if err := joueur.Sauvegarder(); err != nil {
				fmt.Println("⚠️  Erreur lors de la sauvegarde automatique:", err)
			} else {
				fmt.Println("✅ Sauvegarde réussie !")
			}
		}
	}
	
	if actionCount >= maxActions {
		fmt.Printf("\n⚠️  Limite d'actions atteinte (%d). Le jeu se ferme pour éviter une surcharge.\n", maxActions)
		fmt.Println("Votre progression a été sauvegardée automatiquement.")
	}
}

// menuPrincipalExploration affiche le menu principal d'exploration
func menuPrincipalExploration(gameMap *world.Map, joueur *character.Character) bool {
	options := []string{
		"Explorer cette zone",
		"Se déplacer",
		"Voir la carte complète",
		"Afficher le statut du personnage",
		"Quitter le jeu",
	}
	
	ui.AfficherMenu("Que voulez-vous faire ?", options)
	choix := utils.ScanChoice("Votre choix : ", options)
	
	switch choix {
	case 1:
		explorerZoneActuelle(gameMap, joueur)
	case 2:
		seDeplacer(gameMap, joueur)
	case 3:
		gameMap.AfficherMap()
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	case 4:
		afficherStatutPersonnage(joueur)
	case 5:
		fmt.Println("Merci d'avoir joué à World of Milousques !")
		return false
	}
	
	return true
}

// explorerZoneActuelle ouvre le menu d'exploration de la zone actuelle
func explorerZoneActuelle(gameMap *world.Map, joueur *character.Character) {
	zone := gameMap.GetCurrentZone()
	
	fmt.Printf("\n🏠  === %s === 🏠\n", zone.Nom)
	fmt.Println(zone.Description)
	fmt.Println()
	
	zoneActionCount := 0
	maxZoneActions := 50 // Limite les actions dans une zone spécifique
	
	for zoneActionCount < maxZoneActions {
		zoneActionCount++
		
		options := []string{}
		
		// Vérifier si on est à Astrab pour les options spéciales
		estAstrab := strings.Contains(zone.Nom, "Astrab")
		
		// Ajouter les options disponibles selon le contenu de la zone
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
		
		// Vérification de sécurité - toujours au moins l'option "Retour"
		if len(options) == 0 {
			fmt.Println("Erreur : Aucune option disponible. Retour automatique.")
			return
		}
		
		if len(options) == 1 {
			fmt.Println("Cette zone semble vide... Il n'y a rien d'intéressant ici.")
			fmt.Println("Appuyez sur Entrée pour retourner à la carte.")
			fmt.Scanln()
			return
		}
		
		ui.AfficherMenu(fmt.Sprintf("Explorer %s", zone.Nom), options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		// Vérification de sécurité pour le choix
		if choix < 1 || choix > len(options) {
			fmt.Printf("Choix invalide (%d). Retour automatique.\n", choix)
			return
		}
		
		// Utilisation d'un index plus clair pour éviter les erreurs de logique
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
				// Vérifier si le joueur est mort
				if joueur.Pdv <= 0 {
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
			// Forge
			currentIndex++
			if choix == currentIndex {
				craft.AfficherForge(joueur)
				continue
			}
			
			// Marchand
			currentIndex++
			if choix == currentIndex {
				commerce.AfficherMarchand(joueur)
				continue
			}
			
			// Banque
			currentIndex++
			if choix == currentIndex {
				banque.AfficherBanque(joueur)
				continue
			}
		}
		
		// Retour à la carte (toujours la dernière option)
		currentIndex++
		if choix == currentIndex {
			return
		}
		
		// Si aucune condition n'est remplie, sortir par sécurité
		fmt.Printf("⚠️  Erreur de logique avec le choix %d. Retour automatique.\n", choix)
		return
	}
	
	// Vérifier si on a atteint la limite d'actions dans cette zone
	if zoneActionCount >= maxZoneActions {
		fmt.Printf("\n⚠️  Trop d'actions dans cette zone (%d). Retour automatique à la carte.\n", maxZoneActions)
		fmt.Println("Appuyez sur Entrée pour continuer...")
		fmt.Scanln()
		return
	}
}

// seDeplacer gère le déplacement du joueur sur la map avec ZQSD
func seDeplacer(gameMap *world.Map, joueur *character.Character) {
	fmt.Println("\nDéplacements possibles :")
	fmt.Println("Z = Nord | S = Sud | Q = Ouest | D = Est | A = Annuler")
	
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
	
	if len(optionsDisponibles) == 1 {
		fmt.Println("Vous ne pouvez pas vous déplacer d'ici !")
		return
	}
	
	ui.AfficherMenu("Choisir une direction", optionsDisponibles)
	choixInput := utils.ScanString("Tapez Z/Q/S/D pour vous déplacer (ou A pour annuler) : ", 1)
	choixInput = strings.ToUpper(strings.TrimSpace(choixInput))
	
	direction := ""
	nomDirection := ""
	
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
		return
	default:
		fmt.Println("Direction invalide ! Utilisez Z/Q/S/D ou A.")
		return
	}
	
	if gameMap.MoveToWithCharacter(direction, joueur) {
		newZone := gameMap.GetCurrentZone()
		fmt.Printf("\n🚶 Vous vous déplacez vers le %s...\n", nomDirection)
		fmt.Printf("📍 Vous arrivez à : %s\n", newZone.Nom)
		
		// Afficher le nombre total de zones découvertes
		nombreZones := joueur.ObtenirNombreZonesDecouvertes()
		fmt.Printf("🗺️  Zones découvertes : %d/25\n", nombreZones)
		
		fmt.Println("\nAppuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// recolterRessources permet au joueur de récolter des ressources
func recolterRessources(zone *world.Zone, joueur *character.Character) {
	if len(zone.Ressources) == 0 {
		fmt.Println("Il n'y a pas de ressources à récolter ici.")
		return
	}
	
	fmt.Println("\n🌿 === RÉCOLTE DE RESSOURCES === 🌿")
	fmt.Println("Ressources disponibles dans cette zone :")
	
	for i, ressource := range zone.Ressources {
		fmt.Printf("%d. %s (Valeur: %d pièces)\n", i+1, ressource.Nom, ressource.Valeur)
	}
	
	options := []string{"Récolter toutes les ressources", "Retour"}
	ui.AfficherMenu("Récolte", options)
	choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
	
	if choix == 1 {
		// Récolter toutes les ressources
		joueur.Inventaire.Recolter(zone.Ressources)
		fmt.Printf("✅ Vous avez récolté %d ressources !\n", len(zone.Ressources))
		
		// Sauvegarder l'état complet de la zone dans le personnage
		x, y := joueur.ObtenirPosition()
		
		// Convertir les ressources et monstres restants
		ressourcesNoms := []string{}
		for _, res := range zone.Ressources {
			ressourcesNoms = append(ressourcesNoms, res.Nom)
		}
		
		monstresRestants := []character.MonstreState{}
		for _, mon := range zone.Monstres {
			monstresRestants = append(monstresRestants, character.MonstreState{
				Nom: mon.Nom,
				Pv: mon.Pv,
				Attaque: mon.Attaque,
			})
		}
		
		joueur.SauvegarderEtatZoneComplete(x, y, ressourcesNoms, monstresRestants)
		
		// Vider les ressources de la zone
		zone.Ressources = []item.Item{}
		
		// Sauvegarde automatique après récolte
		if err := joueur.Sauvegarder(); err != nil {
			fmt.Println("⚠️  Erreur lors de la sauvegarde automatique:", err)
		} else {
			fmt.Println("💾 Progression sauvegardée automatiquement")
		}
		
		fmt.Println("Appuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// affronterMonstre permet au joueur d'affronter les monstres de la zone
func affronterMonstre(zone *world.Zone, joueur *character.Character) {
	if len(zone.Monstres) == 0 {
		fmt.Println("Il n'y a pas de monstres à affronter ici.")
		return
	}
	
	fmt.Println("\n⚔️  === MONSTRES DE LA ZONE === ⚔️")
	
	options := make([]string, 0)
	for i, monstre := range zone.Monstres {
		options = append(options, fmt.Sprintf("%s (PV: %d, Attaque: %d)", 
			monstre.Nom, monstre.Pv, monstre.Attaque))
		_ = i // Éviter le warning unused variable
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir un adversaire", options)
	choix := utils.ScanChoice("Quel monstre voulez-vous affronter ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	// Combat
	monstreChoisi := &zone.Monstres[choix-1]
	fmt.Printf("\n🥊 Combat contre %s !\n", monstreChoisi.Nom)
	
	fight.Fight(joueur, monstreChoisi)
	
	// Si le monstre est vaincu, le retirer de la zone
	if monstreChoisi.Pv <= 0 {
		// Créer une nouvelle slice sans le monstre vaincu
		nouveauxMonstres := make([]fight.Ennemi, 0)
		for i, m := range zone.Monstres {
			if i != choix-1 {
				nouveauxMonstres = append(nouveauxMonstres, m)
			}
		}
		zone.Monstres = nouveauxMonstres
		fmt.Println("🏆 Le monstre a été vaincu et ne reviendra plus dans cette zone !")
		
		// Sauvegarder l'état complet de la zone après modification des monstres
		x, y := joueur.ObtenirPosition()
		
		// Convertir les ressources et monstres restants
		ressourcesNoms := []string{}
		for _, res := range zone.Ressources {
			ressourcesNoms = append(ressourcesNoms, res.Nom)
		}
		
		monstresRestants := []character.MonstreState{}
		for _, mon := range zone.Monstres {
			monstresRestants = append(monstresRestants, character.MonstreState{
				Nom: mon.Nom,
				Pv: mon.Pv,
				Attaque: mon.Attaque,
			})
		}
		
		joueur.SauvegarderEtatZoneComplete(x, y, ressourcesNoms, monstresRestants)
	}
	
	// Sauvegarde automatique après combat (victoire ou fuite)
	if err := joueur.Sauvegarder(); err != nil {
		fmt.Println("⚠️  Erreur lors de la sauvegarde automatique:", err)
	} else {
		fmt.Println("💾 Progression sauvegardée automatiquement")
	}
	
	fmt.Println("Appuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// parlerAuxPNJs permet d'interagir avec les PNJs de la zone
func parlerAuxPNJs(zone *world.Zone, joueur *character.Character) {
	if len(zone.PNJs) == 0 {
		fmt.Println("Il n'y a personne à qui parler ici.")
		return
	}
	
	fmt.Println("\n💬 === HABITANTS DE LA ZONE === 💬")
	
	options := make([]string, 0)
	for _, pnj := range zone.PNJs {
		options = append(options, pnj.Nom)
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Parler à", options)
	choix := utils.ScanChoice("À qui voulez-vous parler ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	pnj := zone.PNJs[choix-1]
	fmt.Printf("\n🗣️  %s :\n", pnj.Nom)
	fmt.Printf("\"%s\"\n", pnj.Dialogue)
	
	// Vérifier si le joueur a une quête à rendre à ce PNJ
	queteARendreExiste := false
	for _, q := range joueur.Quetes {
		if q.DonneurPNJ == pnj.Nom && q.Accomplie && !q.Rendue {
			queteARendreExiste = true
			break
		}
	}
	
	// Vérifier si le joueur a déjà cette quête
	queteExiste := false
	if pnj.Quete != "" {
		for _, q := range joueur.Quetes {
			if q.Nom == pnj.Quete && !q.Rendue {
				queteExiste = true
				break
			}
		}
	}
	
	if queteARendreExiste {
		// Proposer de rendre les quêtes
		fmt.Println("\n🎉 Ce PNJ a des récompenses pour vous !")
		options := []string{"Rendre quête(s)", "Retour"}
		ui.AfficherMenu("Actions", options)
		choixAction := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		if choixAction == 1 {
			// Rendre toutes les quêtes complétées pour ce PNJ
			queteRendue := false
			for _, q := range joueur.Quetes {
				if q.DonneurPNJ == pnj.Nom && q.Accomplie && !q.Rendue {
					joueur.RendreQuete(q.Nom)
					queteRendue = true
				}
			}
			
			// Sauvegarde automatique après rendu de quête
			if queteRendue {
				if err := joueur.Sauvegarder(); err != nil {
					fmt.Println("⚠️  Erreur lors de la sauvegarde automatique:", err)
				} else {
					fmt.Println("💾 Progression sauvegardée automatiquement")
				}
			}
			
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
			return
		}
	} else if pnj.Quete != "" && !queteExiste {
		// Proposer une nouvelle quête
		fmt.Printf("\n📜 Quête proposée : %s\n", pnj.Quete)
		fmt.Printf("🎁 Récompense : %s\n", pnj.Recompense)
		
		options := []string{"Accepter la quête", "Refuser", "Retour"}
		ui.AfficherMenu("Quête", options)
		choixQuete := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		if choixQuete == 1 {
			// Gérer les nouvelles quêtes de combat
			gererNouvelleQuete(joueur, pnj)
			fmt.Println("✅ Quête acceptée et ajoutée à votre journal !")
		}
	} else if queteExiste {
		fmt.Println("\nℹ️ Vous avez déjà cette quête en cours.")
	} else {
		fmt.Println("\n😊 Ce PNJ n'a pas de quête pour le moment.")
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// afficherStatutPersonnage affiche les informations du personnage avec options
func afficherStatutPersonnage(joueur *character.Character) {
	statutActionCount := 0
	maxStatutActions := 20 // Limite les actions dans le menu de statut
	
	for statutActionCount < maxStatutActions {
		statutActionCount++
		
		// Utiliser l'affichage standard avec quelques ajouts
		fmt.Println("\n📊 === STATUT DU PERSONNAGE === 📊")
		afficherStatutComplet(joueur)
		
		// Menu d'actions
		options := []string{}
		if joueur.Inventaire.Potions > 0 {
			options = append(options, "🧪 Utiliser une potion")
		}
		options = append(options, "🎒 Gérer l'inventaire")
		options = append(options, "Retour")
		
		if len(options) > 1 { // Plus que juste "Retour"
			ui.AfficherMenu("Actions disponibles", options)
			choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
			
			// Vérification de sécurité pour le choix
			if choix < 1 || choix > len(options) {
				fmt.Printf("Choix invalide (%d). Retour automatique.\n", choix)
				return
			}
			
			currentIndex := 0
			
			// Utiliser potion
			if joueur.Inventaire.Potions > 0 {
				currentIndex++
				if choix == currentIndex {
					joueur.UtiliserPotion()
					fmt.Println("\nAppuyez sur Entrée pour continuer...")
					fmt.Scanln()
					continue
				}
			}
			
			// Gérer inventaire
			currentIndex++
			if choix == currentIndex {
				gererInventaire(joueur)
				continue
			}
			
			// Retour (dernière option)
			currentIndex++
			if choix == currentIndex {
				return
			}
			
			// Si aucune condition n'est remplie, sortir par sécurité
			fmt.Printf("⚠️  Erreur de logique avec le choix %d. Retour automatique.\n", choix)
			return
		} else {
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
			return
		}
	}
	
	// Vérifier si on a atteint la limite d'actions du menu statut
	if statutActionCount >= maxStatutActions {
		fmt.Printf("\n⚠️  Trop d'actions dans le menu de statut (%d). Retour automatique.\n", maxStatutActions)
		fmt.Println("Appuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// gererInventaire permet de gérer l'inventaire avec équipement
func gererInventaire(joueur *character.Character) {
	inventaireActionCount := 0
	maxInventaireActions := 30 // Limite les actions dans la gestion d'inventaire
	
	for inventaireActionCount < maxInventaireActions {
		inventaireActionCount++
		
		fmt.Println("\n🎒 === GESTION DE L'INVENTAIRE === 🎒")
		
		if len(joueur.Inventaire.Items) == 0 {
			fmt.Println("Votre inventaire est vide.")
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
			return
		}
		
		joueur.Inventaire.Afficher()
		
		// Vérifier s'il y a des objets équipables
		armesDisponibles := []int{}
		armuresDisponibles := []int{}
		
		for i, itemObj := range joueur.Inventaire.Items {
			switch itemObj.Type {
			case item.TypeArme:
				armesDisponibles = append(armesDisponibles, i)
			case item.TypeCasque, item.TypeTorse, item.TypeJambiere:
				armuresDisponibles = append(armuresDisponibles, i)
			}
		}
		
		options := []string{}
		if len(armesDisponibles) > 0 {
			options = append(options, "⚔️  Équiper une arme")
		}
		if len(armuresDisponibles) > 0 {
			options = append(options, "🛡️  Équiper une armure")
		}
		options = append(options, "Retour")
		
		if len(options) == 1 { // Seulement "Retour"
			fmt.Println("\nAucun objet équipable dans votre inventaire.")
			fmt.Println("Appuyez sur Entrée pour continuer...")
			fmt.Scanln()
			return
		}
		
		ui.AfficherMenu("Actions disponibles", options)
		choix := utils.ScanChoice("Que voulez-vous faire ? ", options)
		
		// Vérification de sécurité pour le choix
		if choix < 1 || choix > len(options) {
			fmt.Printf("Choix invalide (%d). Retour automatique.\n", choix)
			return
		}
		
		currentIndex := 0
		
		// Équiper arme
		if len(armesDisponibles) > 0 {
			currentIndex++
			if choix == currentIndex {
				equiperArme(joueur, armesDisponibles)
				continue
			}
		}
		
		// Équiper armure
		if len(armuresDisponibles) > 0 {
			currentIndex++
			if choix == currentIndex {
				equiperArmure(joueur, armuresDisponibles)
				continue
			}
		}
		
		// Retour (dernière option)
		currentIndex++
		if choix == currentIndex {
			return
		}
		
		// Si aucune condition n'est remplie, sortir par sécurité
		fmt.Printf("⚠️  Erreur de logique avec le choix %d. Retour automatique.\n", choix)
		return
	}
	
	// Vérifier si on a atteint la limite d'actions de l'inventaire
	if inventaireActionCount >= maxInventaireActions {
		fmt.Printf("\n⚠️  Trop d'actions dans l'inventaire (%d). Retour automatique.\n", maxInventaireActions)
		fmt.Println("Appuyez sur Entrée pour continuer...")
		fmt.Scanln()
	}
}

// equiperArme gère l'équipement d'armes
func equiperArme(joueur *character.Character, armesDisponibles []int) {
	fmt.Println("\n⚔️  === ÉQUIPER UNE ARME === ⚔️")
	
	options := []string{}
	for _, index := range armesDisponibles {
		arme := joueur.Inventaire.Items[index]
		options = append(options, fmt.Sprintf("%s - %s", arme.Nom, arme.Effet))
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir une arme", options)
	choix := utils.ScanChoice("Quelle arme voulez-vous équiper ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	// Équiper l'arme choisie
	indexArme := armesDisponibles[choix-1]
	arme := joueur.Inventaire.Items[indexArme]
	
	// Retirer l'arme de l'inventaire
	nouvelInventaire := make([]item.Item, 0)
	for i, item := range joueur.Inventaire.Items {
		if i != indexArme {
			nouvelInventaire = append(nouvelInventaire, item)
		}
	}
	joueur.Inventaire.Items = nouvelInventaire
	
	// Équiper
	joueur.EquiperArme(arme)
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// equiperArmure gère l'équipement d'armures
func equiperArmure(joueur *character.Character, armuresDisponibles []int) {
	fmt.Println("\n🛡️  === ÉQUIPER UNE ARMURE === 🛡️")
	
	options := []string{}
	for _, index := range armuresDisponibles {
		armure := joueur.Inventaire.Items[index]
		options = append(options, fmt.Sprintf("%s - %s", armure.Nom, armure.Effet))
	}
	options = append(options, "Retour")
	
	ui.AfficherMenu("Choisir une armure", options)
	choix := utils.ScanChoice("Quelle armure voulez-vous équiper ? ", options)
	
	if choix == len(options) {
		return // Retour
	}
	
	// Équiper l'armure choisie
	indexArmure := armuresDisponibles[choix-1]
	armure := joueur.Inventaire.Items[indexArmure]
	
	// Retirer l'armure de l'inventaire
	nouvelInventaire := make([]item.Item, 0)
	for i, item := range joueur.Inventaire.Items {
		if i != indexArmure {
			nouvelInventaire = append(nouvelInventaire, item)
		}
	}
	joueur.Inventaire.Items = nouvelInventaire
	
	// Déterminer le type d'armure et équiper
	switch armure.Type {
	case item.TypeCasque:
		joueur.EquiperCasque(armure)
	case item.TypeTorse:
		joueur.EquiperTorse(armure)
	case item.TypeJambiere:
		joueur.EquiperJambiere(armure)
	default:
		// Ancien système pour compatibilité - équiper sur le torse
		joueur.EquiperTorse(armure)
	}
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// gererNouvelleQuete gère l'attribution des nouvelles quêtes de combat spécifiques
func gererNouvelleQuete(joueur *character.Character, pnj world.PNJ) {
	// Définir les quêtes spéciales selon le nom du PNJ
	switch pnj.Nom {
	case "Gawr Gura":
		objectifs := []character.ObjectifCombat{
			{NomMonstre: "Crabe Hijacob", QuantiteRequise: 5, QuantiteActuelle: 0},
			{NomMonstre: "Moumoule", QuantiteRequise: 3, QuantiteActuelle: 0},
		}
		joueur.AjouterQueteCombat(pnj.Quete, pnj.Nom, objectifs, 300, 3, 3)
		
	case "Houshou Marine":
		objectifs := []character.ObjectifCombat{
			{NomMonstre: "Moutmout", QuantiteRequise: 5, QuantiteActuelle: 0},
			{NomMonstre: "Retourneur de panneaux", QuantiteRequise: 3, QuantiteActuelle: 0},
		}
		joueur.AjouterQueteCombat(pnj.Quete, pnj.Nom, objectifs, 300, 3, 3)
		
	case "Fillian":
		objectifs := []character.ObjectifCombat{
			{NomMonstre: "Kairis", QuantiteRequise: 8, QuantiteActuelle: 0},
		}
		joueur.AjouterQueteCombat(pnj.Quete, pnj.Nom, objectifs, 300, 3, 3)
		
	case "Shxtou":
		objectifs := []character.ObjectifCombat{
			{NomMonstre: "Ecumouilles", QuantiteRequise: 8, QuantiteActuelle: 0},
		}
		joueur.AjouterQueteCombat(pnj.Quete, pnj.Nom, objectifs, 300, 3, 3)
		
	default:
		// Quête classique
		joueur.ProposerEtAjouterQuete(pnj.Quete, pnj.Recompense)
	}
}

// afficherStatutComplet affiche le statut détaillé du personnage de façon compacte
func afficherStatutComplet(joueur *character.Character) {
	fmt.Printf("Nom : %s | Classe : %s | Niveau : %d\n", joueur.Nom, joueur.Classe.Nom, joueur.Niveau)
	fmt.Printf("PV : %d/%d | Mana : %d/%d | XP : %d/%d\n", 
		joueur.Pdv, joueur.PdvMax, joueur.Mana, joueur.ManaMax, joueur.Experience, joueur.CalculerXPRequis())
	fmt.Printf("💰 Argent : %d | 🧆 Potions : %d | 🗺️ Zones : %d/25\n", 
		joueur.Argent, joueur.Inventaire.Potions, joueur.ObtenirNombreZonesDecouvertes())
	
	// Équipement compact
	equipements := []string{}
	if joueur.ArmeEquipee != nil {
		equipements = append(equipements, "⚔️ "+joueur.ArmeEquipee.Nom)
	}
	if joueur.CasqueEquipe != nil {
		equipements = append(equipements, "🪖 "+joueur.CasqueEquipe.Nom)
	}
	if joueur.TorseEquipe != nil {
		equipements = append(equipements, "👕 "+joueur.TorseEquipe.Nom)
	}
	if joueur.JambiereEquipee != nil {
		equipements = append(equipements, "👖 "+joueur.JambiereEquipee.Nom)
	}
	
	if len(equipements) > 0 {
		fmt.Println("Équipement :", strings.Join(equipements, ", "))
		bonusAttaque := joueur.CalculerAttaqueBonus()
		bonusDefense := joueur.CalculerDefenseBonus()
		if bonusAttaque > 0 || bonusDefense > 0 {
			fmt.Printf("Bonus : +%d Attaque, +%d Défense\n", bonusAttaque, bonusDefense)
		}
	} else {
		fmt.Println("Aucun équipement")
	}
	
	// Inventaire compact
	if len(joueur.Inventaire.Items) > 0 {
		fmt.Printf("Inventaire : %d objets", len(joueur.Inventaire.Items))
		if len(joueur.Inventaire.Items) <= 3 {
			for _, item := range joueur.Inventaire.Items {
				fmt.Printf(", %s", item.Nom)
			}
		}
		fmt.Println()
	}
	
	joueur.AfficherQuetes()
}
