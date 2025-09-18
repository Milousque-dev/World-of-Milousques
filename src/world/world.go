// Package world gère la carte du jeu, les zones, les PNJs et la génération de contenu
// Structure la carte 5x5 avec différents biomes : champs, forêts, mines, rivières et routes
package world

import (
	"fmt"
	"world_of_milousques/character"
	"world_of_milousques/fight"
	"world_of_milousques/item"
)

// PNJ représente un personnage non-joueur
type PNJ struct {
	Nom       string
	Dialogue  string
	Quete     string
	Recompense string
}

// Zone représente une sous-zone de la map
type Zone struct {
	Nom         string
	Description string
	Ressources  []item.Item
	Monstres    []fight.Ennemi
	PNJs        []PNJ
	Visitee     bool
}

// Position du joueur sur la map
type Position struct {
	X, Y int
}

// Map représente la grille 5x5 du monde
type Map struct {
	Zones    [5][5]Zone
	Position Position
}

// NewMap crée une nouvelle map avec des zones génériques
func NewMap() *Map {
	m := &Map{
		Position: Position{X: 2, Y: 2}, // Position centrale au départ
	}
	
	// Initialiser toutes les zones avec du contenu générique
	m.initializeZones()
	
	return m
}


// RestaurerPosition met à jour la position de la map depuis un personnage
func (m *Map) RestaurerPosition(x, y int) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		m.Position = Position{X: x, Y: y}
	}
}

// RestaurerEtatDecouverte met à jour l'état de découverte des zones depuis un personnage
func (m *Map) RestaurerEtatDecouverte(zonesDecouvertes [5][5]bool) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if zonesDecouvertes[y][x] {
				m.Zones[y][x].Visitee = true
			}
		}
	}
}

// RestaurerEtatRessources met à jour l'état des ressources ET monstres des zones depuis un personnage
func (m *Map) RestaurerEtatRessources(characterInterface interface{}) {
	// Essayer d'utiliser le type Character directement
	if char, ok := characterInterface.(*character.Character); ok {
		m.restaurerAvecCharacter(char)
		return
	}
	
	// Fallback pour l'ancienne méthode
	if charInterface, ok := characterInterface.(interface {
		ZoneRessourcesRecoltees(int, int) bool
		ZoneMonstresVaincus(int, int) bool
	}); ok {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				// Méthode binaire classique (ancienne)
				if charInterface.ZoneRessourcesRecoltees(x, y) {
					m.Zones[y][x].Ressources = []item.Item{}
				}
				if charInterface.ZoneMonstresVaincus(x, y) {
					m.Zones[y][x].Monstres = []fight.Ennemi{}
				}
			}
		}
	}
}

// restaurerAvecCharacter utilise l'état détaillé sauvegardé pour chaque zone
func (m *Map) restaurerAvecCharacter(char *character.Character) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			// Obtenir l'état sauvegardé détaillé de la zone
			if etat, existe := char.ObtenirEtatZone(x, y); existe && etat.Visitee {
				// Restaurer les ressources selon l'état sauvegardé
				m.restaurerRessourcesZone(&m.Zones[y][x], etat.RessourcesRestantes)
				
				// Restaurer les monstres selon l'état sauvegardé
				m.restaurerMonstresZone(&m.Zones[y][x], etat.MonstresRestants)
			}
		}
	}
}

// restaurerRessourcesZone restaure les ressources d'une zone selon l'état sauvegardé
func (m *Map) restaurerRessourcesZone(zone *Zone, ressourcesRestantes []string) {
	// Créer une nouvelle liste de ressources basée sur l'état sauvegardé
	nouvellesRessources := []item.Item{}
	for _, nomRessource := range ressourcesRestantes {
		nouvellesRessources = append(nouvellesRessources, item.NewItem(nomRessource))
	}
	zone.Ressources = nouvellesRessources
}

// restaurerMonstresZone restaure les monstres d'une zone selon l'état sauvegardé
func (m *Map) restaurerMonstresZone(zone *Zone, monstresRestants []character.MonstreState) {
	// Créer une nouvelle liste de monstres basée sur l'état sauvegardé
	nouveauxMonstres := []fight.Ennemi{}
	for _, monstreState := range monstresRestants {
		nouveauxMonstres = append(nouveauxMonstres, fight.Ennemi{
			Nom: monstreState.Nom,
			Pv: monstreState.Pv,
			Attaque: monstreState.Attaque,
		})
	}
	zone.Monstres = nouveauxMonstres
}

// GetCurrentZone retourne la zone actuelle du joueur
func (m *Map) GetCurrentZone() *Zone {
	return &m.Zones[m.Position.Y][m.Position.X]
}

// GetZoneAt retourne la zone à la position spécifiée (pour les tests)
func (m *Map) GetZoneAt(x, y int) *Zone {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return &m.Zones[y][x]
	}
	return nil
}

// CanMoveTo vérifie si le joueur peut se déplacer vers une direction
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
		return false
	}
	
	return newX >= 0 && newX < 5 && newY >= 0 && newY < 5
}

// MoveTo déplace le joueur dans une direction
func (m *Map) MoveTo(direction string) bool {
	return m.MoveToWithCharacter(direction, nil)
}

// MoveToWithCharacter déplace le joueur et met à jour la sauvegarde du personnage
func (m *Map) MoveToWithCharacter(direction string, character interface{}) bool {
	if !m.CanMoveTo(direction) {
		return false
	}
	
	// Déplacer selon la direction
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
	
	m.GetCurrentZone().Visitee = true
	
	// Sauvegarder l'état si un personnage est fourni
	if char, ok := character.(interface {
		SauvegarderPositionMap(int, int)
		MarquerZoneDecouverte(int, int)
	}); ok {
		char.SauvegarderPositionMap(m.Position.X, m.Position.Y)
		char.MarquerZoneDecouverte(m.Position.X, m.Position.Y)
	}
	
	return true
}

// AfficherMap affiche la map ASCII avec la position du joueur
func (m *Map) AfficherMap() {
	fmt.Println("\n=== CARTE DU MONDE ===")
	fmt.Println()
	
	for y := 0; y < 5; y++ {
		// Ligne du haut de chaque rangée
		for x := 0; x < 5; x++ {
			fmt.Print("+-------")
		}
		fmt.Println("+")
		
		// Ligne du milieu avec le contenu
		for x := 0; x < 5; x++ {
			zone := &m.Zones[y][x]
			symbol := " "
			
			if x == m.Position.X && y == m.Position.Y {
				symbol = "♦" // Position du joueur
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
	
	fmt.Println("\nLégende: ♦ = Vous | ○ = Visitée | ? = Inconnue")
	fmt.Printf("Position actuelle: %s (%d,%d)\n", 
		m.GetCurrentZone().Nom, m.Position.X+1, m.Position.Y+1)
}

// initializeZones remplit la map avec du contenu géographiquement cohérent
func (m *Map) initializeZones() {
	// Disposition géographique :
	// Champs (haut-gauche) | Transition | Mines (haut-droite)
	// Forêt (gauche) | Astrab (centre) | Rivière (droite)
	// Forêt (bas-gauche) | Transition | Rivière (bas-droite)
	
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
			
			// Vérifier si cette position doit être une "Route"
			// Coordonnées spécifiées : 3,1 / 3,2 / 3,4 / 3,5 / 2,3 / 1,3 / 4,3 / 5,3
			// Converties en index (x,y) : (2,0) / (2,1) / (2,3) / (2,4) / (1,2) / (0,2) / (3,2) / (4,2)
			isRoute := (x == 2 && (y == 0 || y == 1 || y == 3 || y == 4)) || (y == 2 && (x == 0 || x == 1 || x == 3 || x == 4))
			
			// Vérifier si cette position doit être un "Champs" spécial
			// Coordonnées spécifiées : 5,2 / 5,1 / 4,1 / 4,2
			// Converties en index (x,y) : (4,1) / (4,0) / (3,0) / (3,1)
			isSpecialChamps := (x == 4 && (y == 0 || y == 1)) || (x == 3 && (y == 0 || y == 1))
			
			// Vérifier si cette position doit être une "Mine" spéciale
			// Coordonnées spécifiées : 2,2 / 1,2 / 1,1 / 2,1
			// Converties en index (x,y) : (1,1) / (0,1) / (0,0) / (1,0)
			isSpecialMine := (x == 0 && (y == 0 || y == 1)) || (x == 1 && (y == 0 || y == 1))
			
			// Vérifier si cette position doit être une "Forêt" spéciale
			// Coordonnées spécifiées : 1,4 / 2,4 / 1,5 / 2,5 (bas gauche)
			// Converties en index (x,y) : (0,3) / (1,3) / (0,4) / (1,4)
		isSpecialForet := (x == 0 && (y == 3 || y == 4)) || (x == 1 && (y == 3 || y == 4))
			
			// Vérifier si cette position doit être une "Rivière" spéciale
			// Les 4 zones restantes identifiées : (3,3), (4,3), (3,4), (4,4)
			isSpecialRiviere := (x == 3 && y == 3) || (x == 4 && y == 3) || (x == 3 && y == 4) || (x == 4 && y == 4)
			
			if isRoute {
				// Configuration spéciale pour les routes
				zone.Nom = "Route"
				zone.Description = "A la croisée des chemins, on trouve tous les gros malins !"
				zone.Ressources = []item.Item{} // Aucune ressource
				zone.Monstres = []fight.Ennemi{} // Aucun monstre
				zone.PNJs = []PNJ{} // Aucun PNJ dans les zones de route
				zone.Visitee = false
			} else if isSpecialChamps {
				// Configuration spéciale pour les nouveaux champs
				zone.Nom = "Champs"
				zone.Description = "Le territoire de Mylène"
				m.setupSpecialChampsZone(zone)
				zone.Visitee = false
				// Ajouter les PNJs spéciaux aux champs spéciaux
				m.ajouterPNJsSpeciaux(zone, x, y)
			} else if isSpecialMine {
				// Configuration spéciale pour les nouvelles mines
				zone.Nom = "Mine"
				zone.Description = "Depuis les attaques de Kairis, les mineurs ont déserter l'endroit"
				m.setupSpecialMineZone(zone)
				zone.Visitee = false
				// Ajouter les PNJs spéciaux aux mines spéciales
				m.ajouterPNJsSpeciaux(zone, x, y)
			} else if isSpecialForet {
				// Configuration spéciale pour les nouvelles forêts
				zone.Nom = "Forêt"
				zone.Description = "Construire un parking pour lutter contre la forestation"
				m.setupSpecialForetZone(zone)
				zone.Visitee = false
				// Ajouter les PNJs spéciaux aux forêts spéciales
				m.ajouterPNJsSpeciaux(zone, x, y)
			} else if isSpecialRiviere {
				// Configuration spéciale pour les nouvelles rivières
				zone.Nom = "Rivière"
				zone.Description = "On aurais préférer 3 rivières"
				m.setupSpecialRiviereZone(zone)
				zone.Visitee = false
				// Ajouter les PNJs spéciaux aux rivières spéciales
				m.ajouterPNJsSpeciaux(zone, x, y)
			} else {
				// Configuration normale selon les anciens tableaux
				zoneType := zoneTypes[y][x]
				
				// Astrab a une configuration spéciale
				if zoneType == "Astrab" {
					zone.Description = "Astrab, la magnifique capitale du royaume. Ses rues pavées fourmillent de marchands, d'artisans et d'aventuriers. Au cœur de la cité se dressent la Grande Forge, le Marché Central et la Banque Royale."
					zone.Visitee = true
					// Astrab n'a pas de monstres ni de ressources mais des PNJs spéciaux
					zone.Ressources = []item.Item{}
					zone.Monstres = []fight.Ennemi{}
					zone.PNJs = []PNJ{
						{Nom: "Maître Karim le Marchand", Dialogue: "Bienvenue dans ma boutique ! J'ai tout ce dont un aventurier a besoin !", Quete: "", Recompense: ""},
						{Nom: "Maître Forgeron Hassan", Dialogue: "Ma forge est à votre disposition pour créer de merveilleux objets !", Quete: "", Recompense: ""},
						{Nom: "Banquier Salomon", Dialogue: "La Banque Royale garde vos biens précieux en sécurité !", Quete: "", Recompense: ""},
						{Nom: "Garde Royale", Dialogue: "Astrab est la cité la plus sûre du royaume, aventurier.", Quete: "", Recompense: ""},
					}
				} else {
				// Générer le contenu selon le type de zone
					m.setupZoneByType(zone, zoneType)
				
				// Ajouter les PNJs spéciaux aux positions demandées
				m.ajouterPNJsSpeciaux(zone, x, y)
				}
			}
		}
	}
}

// setupZoneByType configure une zone selon son type géographique
func (m *Map) setupZoneByType(zone *Zone, zoneType string) {
	switch zoneType {
	case "Champs":
		m.setupChampsZone(zone)
	case "Forêt":
		m.setupForetZone(zone)
	case "Mines":
		m.setupMinesZone(zone)
	case "Rivière":
		m.setupRiviereZone(zone)
	default: // Zones de transition
		m.setupTransitionZone(zone, zoneType)
	}
}

// setupChampsZone configure une zone de champs
func (m *Map) setupChampsZone(zone *Zone) {
	zone.Description = "De vastes champs s'étendent à perte de vue. L'odeur de la terre fertile et des cultures emplit l'air."
	
	// Ressources fixes des champs (15 total)
	zone.Ressources = []item.Item{}
	for i := 0; i < 8; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Laitue Vireuse"))
	}
	for i := 0; i < 7; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Blé"))
	}
	
	// Monstres des champs (4 total)
	zone.Monstres = []fight.Ennemi{
		{Nom: "Moutmout", Pv: 80, Attaque: 25},
		{Nom: "Moutmout", Pv: 80, Attaque: 25},
		{Nom: "Retourneur de panneaux", Pv: 150, Attaque: 40},
		{Nom: "Retourneur de panneaux", Pv: 150, Attaque: 40},
	}
	
	// Aucun PNJ dans les champs
	zone.PNJs = []PNJ{}
}

// setupForetZone configure une zone de forêt
func (m *Map) setupForetZone(zone *Zone) {
	zone.Description = "Une forêt dense où les arbres anciens créent une voûte verdoyante. Le bois crépite sous vos pas."
	
	// Ressources fixes de la forêt (18 total)
	zone.Ressources = []item.Item{}
	for i := 0; i < 12; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Bois"))
	}
	for i := 0; i < 6; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Laitue Vireuse"))
	}
	
	// Monstres de la forêt (5 total)
	zone.Monstres = []fight.Ennemi{
		{Nom: "Ecumouilles", Pv: 100, Attaque: 30},
		{Nom: "Ecumouilles", Pv: 100, Attaque: 30},
		{Nom: "Ecumouilles", Pv: 100, Attaque: 30},
		{Nom: "Ecumouilles", Pv: 100, Attaque: 30},
		{Nom: "Ecumouilles", Pv: 100, Attaque: 30},
	}
	
	// Aucun PNJ dans les forêts
	zone.PNJs = []PNJ{}
}

// setupMinesZone configure une zone de mines
func (m *Map) setupMinesZone(zone *Zone) {
	zone.Description = "Des galeries sombres s'enfoncent dans la montagne. L'écho de vos pas résonne dans les tunnels."
	
	// Ressources fixes des mines (12 total)
	zone.Ressources = []item.Item{}
	for i := 0; i < 12; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Fer"))
	}
	
	// Monstres des mines (3 total)
	zone.Monstres = []fight.Ennemi{
		{Nom: "Kairis", Pv: 110, Attaque: 35},
		{Nom: "Kairis", Pv: 110, Attaque: 35},
		{Nom: "Kairis", Pv: 110, Attaque: 35},
	}
	
	// Aucun PNJ dans les mines
	zone.PNJs = []PNJ{}
}

// setupRiviereZone configure une zone de rivière
func (m *Map) setupRiviereZone(zone *Zone) {
	zone.Description = "Les eaux cristallines de la rivière coulent paisiblement. Des reflets dansent à la surface."
	
	// Ressources fixes de la rivière (16 total)
	zone.Ressources = []item.Item{}
	for i := 0; i < 16; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Pichon"))
	}
	
	// Monstres de la rivière (4 total)
	zone.Monstres = []fight.Ennemi{
		{Nom: "Crabe Hijacob", Pv: 90, Attaque: 20},
		{Nom: "Crabe Hijacob", Pv: 90, Attaque: 20},
		{Nom: "Moumoule", Pv: 250, Attaque: 15},
		{Nom: "Moumoule", Pv: 250, Attaque: 15},
	}
	
	// Aucun PNJ dans les rivières
	zone.PNJs = []PNJ{}
}

// setupTransitionZone configure une zone de transition
func (m *Map) setupTransitionZone(zone *Zone, zoneType string) {
	switch zoneType {
	case "Transition Nord":
		zone.Description = "Une zone de transition entre les champs verdoyants et les mines rocailleuses."
		// Mélange de ressources
		zone.Ressources = []item.Item{
			item.NewItem("Blé"), item.NewItem("Blé"), 
			item.NewItem("Fer"), item.NewItem("Fer"),
		}
		zone.Monstres = []fight.Ennemi{{Nom: "Moutmout", Pv: 80, Attaque: 25}}
	case "Transition Ouest":
		zone.Description = "Le passage entre la forêt mystérieuse et les champs ouverts."
		zone.Ressources = []item.Item{
			item.NewItem("Bois"), item.NewItem("Bois"),
			item.NewItem("Laitue Vireuse"), item.NewItem("Laitue Vireuse"),
		}
		zone.Monstres = []fight.Ennemi{{Nom: "Ecumouilles", Pv: 100, Attaque: 30}}
	case "Transition Est":
		zone.Description = "La frontière entre les terres cultivées et les rivières sauvages."
		zone.Ressources = []item.Item{
			item.NewItem("Pichon"), item.NewItem("Pichon"),
			item.NewItem("Blé"), item.NewItem("Blé"),
		}
		zone.Monstres = []fight.Ennemi{{Nom: "Crabe Hijacob", Pv: 90, Attaque: 20}}
	case "Transition Sud":
		zone.Description = "Un carrefour où se mélangent les influences de la forêt et de la rivière."
		zone.Ressources = []item.Item{
			item.NewItem("Bois"), item.NewItem("Bois"),
			item.NewItem("Pichon"), item.NewItem("Pichon"),
		}
		zone.Monstres = []fight.Ennemi{{Nom: "Ecumouilles", Pv: 100, Attaque: 30}}
	}
	
	zone.PNJs = []PNJ{} // Aucun PNJ selon les nouvelles spécifications
}

// setupSpecialChampsZone configure les nouveaux champs spéciaux avec variation naturelle
func (m *Map) setupSpecialChampsZone(zone *Zone) {
	// Utiliser la position pour déterminer une variation fixe mais naturelle
	positionSeed := m.getZonePositionSeed(zone)
	
	// Génération de ressources (10-20 total)
	totalRessources := 10 + (positionSeed % 11) // 10 à 20
	bleRatio := 0.4 + float64((positionSeed*7)%30)/100 // 40% à 70% de blé
	nbBle := int(float64(totalRessources) * bleRatio)
	nbLaitue := totalRessources - nbBle
	
	zone.Ressources = []item.Item{}
	for i := 0; i < nbBle; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Blé"))
	}
	for i := 0; i < nbLaitue; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Laitue Vireuse"))
	}
	
	// Génération de monstres (3-5 total)
	totalMonstres := 3 + (positionSeed % 3) // 3 à 5
	moutmoutRatio := 0.3 + float64((positionSeed*11)%40)/100 // 30% à 70% de Moutmout
	nbMoutmout := int(float64(totalMonstres) * moutmoutRatio)
	nbRetourneur := totalMonstres - nbMoutmout
	
	zone.Monstres = []fight.Ennemi{}
	for i := 0; i < nbMoutmout; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Moutmout", Pv: 80, Attaque: 25})
	}
	for i := 0; i < nbRetourneur; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Retourneur de panneaux", Pv: 150, Attaque: 40})
	}
	
	// Aucun PNJ dans les nouveaux champs
	zone.PNJs = []PNJ{}
}

// setupSpecialMineZone configure les nouvelles mines spéciales avec variation naturelle
func (m *Map) setupSpecialMineZone(zone *Zone) {
	// Utiliser la position pour déterminer une variation fixe mais naturelle
	positionSeed := m.getZonePositionSeed(zone)
	
	// Génération de ressources (10-20 total) - que du fer
	totalRessources := 10 + (positionSeed % 11) // 10 à 20
	
	zone.Ressources = []item.Item{}
	for i := 0; i < totalRessources; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Fer"))
	}
	
	// Génération de monstres (3-5 total) - que des Kairis
	totalMonstres := 3 + (positionSeed % 3) // 3 à 5
	
	zone.Monstres = []fight.Ennemi{}
	for i := 0; i < totalMonstres; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Kairis", Pv: 110, Attaque: 35})
	}
	
	// Aucun PNJ dans les nouvelles mines (les mineurs ont fui)
	zone.PNJs = []PNJ{}
}

// setupSpecialForetZone configure les nouvelles forêts spéciales avec variation naturelle
func (m *Map) setupSpecialForetZone(zone *Zone) {
	// Utiliser la position pour déterminer une variation fixe mais naturelle
	positionSeed := m.getZonePositionSeed(zone)
	
	// Génération de ressources (10-20 total)
	totalRessources := 10 + (positionSeed % 11) // 10 à 20
	boisRatio := 0.6 + float64((positionSeed*13)%30)/100 // 60% à 90% de bois
	nbBois := int(float64(totalRessources) * boisRatio)
	nbLaitue := totalRessources - nbBois
	
	zone.Ressources = []item.Item{}
	for i := 0; i < nbBois; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Bois"))
	}
	for i := 0; i < nbLaitue; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Laitue Vireuse"))
	}
	
	// Génération de monstres (3-5 total) - que des Ecumouilles
	totalMonstres := 3 + (positionSeed % 3) // 3 à 5
	
	zone.Monstres = []fight.Ennemi{}
	for i := 0; i < totalMonstres; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Ecumouilles", Pv: 100, Attaque: 30})
	}
	
	// Aucun PNJ dans les nouvelles forêts
	zone.PNJs = []PNJ{}
}

// setupSpecialRiviereZone configure les nouvelles rivières spéciales avec variation naturelle
func (m *Map) setupSpecialRiviereZone(zone *Zone) {
	// Utiliser la position pour déterminer une variation fixe mais naturelle
	positionSeed := m.getZonePositionSeed(zone)
	
	// Génération de ressources (10-20 total) - que des Pichons
	totalRessources := 10 + (positionSeed % 11) // 10 à 20
	
	zone.Ressources = []item.Item{}
	for i := 0; i < totalRessources; i++ {
		zone.Ressources = append(zone.Ressources, item.NewItem("Pichon"))
	}
	
	// Génération de monstres (3-5 total)
	totalMonstres := 3 + (positionSeed % 3) // 3 à 5
	crabeRatio := 0.4 + float64((positionSeed*17)%40)/100 // 40% à 80% de crabes
	nbCrabe := int(float64(totalMonstres) * crabeRatio)
	nbMoumoule := totalMonstres - nbCrabe
	
	zone.Monstres = []fight.Ennemi{}
	for i := 0; i < nbCrabe; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Crabe Hijacob", Pv: 90, Attaque: 20})
	}
	for i := 0; i < nbMoumoule; i++ {
		zone.Monstres = append(zone.Monstres, fight.Ennemi{Nom: "Moumoule", Pv: 250, Attaque: 15})
	}
	
	// Aucun PNJ dans les nouvelles rivières
	zone.PNJs = []PNJ{}
}

// getZonePositionSeed génère un seed fixé basé sur la position de la zone
func (m *Map) getZonePositionSeed(zone *Zone) int {
	// Trouver la position de la zone dans la carte
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if &m.Zones[y][x] == zone {
				// Créer un seed unique basé sur la position
				// Utiliser des nombres premiers pour éviter les patterns
				return (x*17 + y*23 + x*y*11) % 1000
			}
		}
	}
	return 42 // Valeur par défaut
}

// ajouterPNJsSpeciaux ajoute les PNJs avec quêtes spéciales aux positions demandées
func (m *Map) ajouterPNJsSpeciaux(zone *Zone, x int, y int) {
	// Position 4,4 -> index (3,3)
	if x == 3 && y == 3 {
		pnjGura := PNJ{
			Nom: "Gawr Gura",
			Dialogue: "Shaaaark ! La danse des crabe hijacob est insupportable",
			Quete: "Nettoyage des Rivières",
			Recompense: "300 or, 3 potions de vie, 3 potions de mana",
		}
		zone.PNJs = append(zone.PNJs, pnjGura)
	}
	
	// Position 4,2 -> index (3,1)
	if x == 3 && y == 1 {
		pnjMarine := PNJ{
			Nom: "Houshou Marine",
			Dialogue: "Les champs sont envahi, la BRUV N est dépasser ! Va apporter la démocratie",
			Quete: "Raid des Champs",
			Recompense: "300 or, 3 potions de vie, 3 potions de mana",
		}
		zone.PNJs = append(zone.PNJs, pnjMarine)
	}
	
	// Position 2,2 -> index (1,1)
	if x == 1 && y == 1 {
		pnjFillian := PNJ{
			Nom: "Fillian",
			Dialogue: "Ces Kairis ont envahi mes mines ! Ils détournent les mineurs ! Il faut les stopper de toute urgence !",
			Quete: "Répression des Kairis",
			Recompense: "300 or, 3 potions de vie, 3 potions de mana",
		}
		zone.PNJs = append(zone.PNJs, pnjFillian)
	}
	
	// Position 2,4 -> index (1,3)
	if x == 1 && y == 3 {
		pnjShxtou := PNJ{
			Nom: "Shxtou",
			Dialogue: "J'en peut plus des ecumouilles, va faire un petit massacre pitié !",
			Quete: "Nettoyage de Forêt",
			Recompense: "300 or, 3 potions de vie, 3 potions de mana",
		}
		zone.PNJs = append(zone.PNJs, pnjShxtou)
	}
}

