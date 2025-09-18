// Package character gère la création, la gestion et les actions des personnages joueurs
// Inclut : création, sauvegarde/chargement, système d'expérience, équipement et quêtes
package character

import (
	"encoding/json"
	"fmt"
	"os"

	"world_of_milousques/classe"
	"world_of_milousques/inventory"
	"world_of_milousques/item"
	"world_of_milousques/ui"
	"world_of_milousques/utils"
)

// ObjectifCombat représente un objectif de combat spécifique
type ObjectifCombat struct {
	NomMonstre string `json:"nom_monstre"`
	QuantiteRequise int `json:"quantite_requise"`
	QuantiteActuelle int `json:"quantite_actuelle"`
}

// ZoneState représente l'état d'une zone (monstres tués, ressources récoltées)
type ZoneState struct {
	Visitee bool `json:"visitee"`
	RessourcesRestantes []string `json:"ressources_restantes"` // Noms des ressources encore présentes
	MonstresRestants []MonstreState `json:"monstres_restants"` // Monstres encore vivants
}

// MonstreState représente l'état d'un monstre
type MonstreState struct {
	Nom string `json:"nom"`
	Pv int `json:"pv"`
	Attaque int `json:"attaque"`
}

// MapState représente l'état complet de la map
type MapState struct {
	Zones [5][5]ZoneState `json:"zones"`
}

// Quete représente une quête avec objectifs de combat
type Quete struct {
	Nom string `json:"nom"`
	Accomplie bool `json:"accomplie"`
	Recompense string `json:"recompense"`
	DonneurPNJ string `json:"donneur_pnj"`
	Rendue bool `json:"rendue"`
	// Nouveau système d'objectifs
	ObjectifsCombat []ObjectifCombat `json:"objectifs_combat,omitempty"`
	RecompenseOr int `json:"recompense_or,omitempty"`
	RecompensePotionsVie int `json:"recompense_potions_vie,omitempty"`
	RecompensePotionsMana int `json:"recompense_potions_mana,omitempty"`
}

type Character struct {
	Nom        string               `json:"nom"`
	Niveau     int                  `json:"niveau"`
	Pdv        int                  `json:"pdv"`
	Mana       int                  `json:"mana"`
	PdvMax     int                  `json:"pdv_max"`
	ManaMax    int                  `json:"mana_max"`
	Experience int                  `json:"experience"`
	Argent     int                  `json:"argent"`
	// Nouveau système d'équipement
	ArmeEquipee     *item.Item       `json:"arme_equipee,omitempty"`
	CasqueEquipe    *item.Item       `json:"casque_equipe,omitempty"`
	TorseEquipe     *item.Item       `json:"torse_equipe,omitempty"`
	JambiereEquipee *item.Item       `json:"jambiere_equipee,omitempty"`
	Classe     classe.Classe        `json:"classe"`
	Inventaire inventory.Inventaire `json:"inventaire"`
	Quetes     []Quete              `json:"quetes"`
	// Nouveaux champs pour le système de sauvegarde avancé
	IntroEffectuee bool              `json:"intro_effectuee"`
	PositionX      int               `json:"position_x"`
	PositionY      int               `json:"position_y"`
	EtatMap        MapState          `json:"etat_map"`
	ZonesDecouvertes [5][5]bool     `json:"zones_decouvertes"`
	// Champs simples pour le suivi des zones vidées
	ZonesRessourcesRecoltees [5][5]bool `json:"zones_ressources_recoltees"`
	ZonesMonstresVaincus [5][5]bool     `json:"zones_monstres_vaincus"`
}

func InitCharacter(nom string, c classe.Classe, niveau int, pdv int, pdvmax int) Character {
	return Character{
		Nom:            nom,
		Niveau:         niveau,
		Pdv:            pdv,
		Mana:           c.ManaMax,
		PdvMax:         c.Pvmax,
		ManaMax:        c.ManaMax,
		Experience:     0,
		Argent:         100, // Argent de départ
		Classe:         c,
		Inventaire:     inventory.Inventaire{},
		Quetes:         []Quete{},
		// Initialiser les nouveaux champs
		IntroEffectuee: false,
		PositionX:      2, // Position centrale
		PositionY:      2, // Position centrale
		EtatMap:        MapState{}, // État de map vide (sera initialisée plus tard)
		ZonesDecouvertes: [5][5]bool{}, // Aucune zone découverte au début
	}
}

func (c *Character) Sauvegarder() error {
	filename := "saves/" + c.Nom + ".json"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	fmt.Println("Personnage sauvegardé dans", filename)
	return nil
}

func Charger(nom string) (*Character, error) {
	filename := "saves/" + nom + ".json"

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c Character
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	// Donner de l'argent de départ aux anciens personnages
	if c.Argent == 0 {
		c.Argent = 100
		fmt.Println("💰 Vous recevez 100 pièces d'or de départ !")
	}

	fmt.Println("Personnage chargé depuis", filename)
	return &c, nil
}

func (c *Character) ProposerEtAjouterQuete(nom string, recompense string) {
	c.Quetes = append(c.Quetes, Quete{Nom: nom, Accomplie: false, Recompense: recompense})
}

func (c *Character) CompleterQuete(nom string) {
	for i := range c.Quetes {
		if c.Quetes[i].Nom == nom {
			c.Quetes[i].Accomplie = true
			fmt.Println("Quête complétée :", nom)
			fmt.Println("Récompense :", c.Quetes[i].Recompense)
			if c.Quetes[i].Recompense == "1 potion" {
				c.Inventaire.Potions++
				fmt.Println("Vous recevez 1 potion !")
			}
			break
		}
	}
}

func (c *Character) AfficherQuetes() {
	quetesActives := []Quete{}
	for _, q := range c.Quetes {
		if !q.Rendue {
			quetesActives = append(quetesActives, q)
		}
	}
	
	if len(quetesActives) == 0 {
		fmt.Println("Aucune quête active.")
		return
	}
	fmt.Println("Quêtes actives :")
	for _, q := range quetesActives {
		status := "En cours"
		if q.Accomplie && !q.Rendue {
			status = "Prête à rendre"
		}
		fmt.Printf("- %s : %s | Récompense : %s\n", q.Nom, status, q.Recompense)
		
		// Afficher les objectifs de combat s'il y en a
		if len(q.ObjectifsCombat) > 0 {
			for _, obj := range q.ObjectifsCombat {
				statutObj := "✅"
				if obj.QuantiteActuelle < obj.QuantiteRequise {
					statutObj = "⏳"
				}
				fmt.Printf("  %s %s : %d/%d\n", statutObj, obj.NomMonstre, obj.QuantiteActuelle, obj.QuantiteRequise)
			}
		}
	}
}

// === SYSTÈME D'EXPÉRIENCE ===

// GagnerExperience fait gagner de l'expérience au personnage
func (c *Character) GagnerExperience(xp int) {
	c.Experience += xp
	fmt.Printf("\n✨ Vous gagnez %d points d'expérience !\n", xp)
	
	// Vérifier si montée de niveau
	xpRequis := c.CalculerXPRequis()
	if c.Experience >= xpRequis {
		c.MonterDeNiveau()
	}
}

// CalculerXPRequis calcule l'XP nécessaire pour le prochain niveau
func (c *Character) CalculerXPRequis() int {
	return c.Niveau * 100 // 100 XP pour niveau 1->2, 200 pour 2->3, etc.
}

// MonterDeNiveau gère la montée de niveau
func (c *Character) MonterDeNiveau() {
	c.Niveau++
	c.Experience = 0 // Reset XP
	
	fmt.Printf("\n🎉 === MONTÉE DE NIVEAU === 🎉\n")
	fmt.Printf("Vous êtes maintenant niveau %d !\n", c.Niveau)
	
	// Choix d'amélioration
	options := []string{"+ 10 PV maximum", "+ 10 Mana maximum"}
	ui.AfficherMenu("Choisissez votre amélioration", options)
	choix := utils.ScanChoice("Votre choix : ", options)
	
	if choix == 1 {
		c.PdvMax += 10
		fmt.Println("💙 Vos PV maximum augmentent de 10 !")
	} else {
		c.ManaMax += 10
		fmt.Println("🔮 Votre Mana maximum augmente de 10 !")
	}
	
	// Restaurer complètement PV et Mana
	c.Pdv = c.PdvMax
	c.Mana = c.ManaMax
	fmt.Println("❤️  Vos PV et Mana sont complètement restaurés !")
	
	fmt.Println("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// === NOUVEAU SYSTÈME D'ÉQUIPEMENT ===

// EquiperArme équipe une arme
func (c *Character) EquiperArme(arme item.Item) {
	// Vérifier la classe requise
	if arme.ClasseRequise != "" && arme.ClasseRequise != c.Classe.Nom {
		fmt.Printf("❌ Vous ne pouvez pas équiper %s (classe requise : %s)\n", arme.Nom, arme.ClasseRequise)
		return
	}
	
	if c.ArmeEquipee != nil {
		// Remettre l'ancienne arme dans l'inventaire
		c.Inventaire.Items = append(c.Inventaire.Items, *c.ArmeEquipee)
	}
	c.ArmeEquipee = &arme
	fmt.Printf("⚔️  Vous équipez : %s (+%d attaque)\n", arme.Nom, arme.Attaque)
}

// EquiperCasque équipe un casque
func (c *Character) EquiperCasque(casque item.Item) {
	c.equiperArmure(&c.CasqueEquipe, casque, "🪖")
}

// EquiperTorse équipe un torse
func (c *Character) EquiperTorse(torse item.Item) {
	c.equiperArmure(&c.TorseEquipe, torse, "👕")
}

// EquiperJambiere équipe des jambières
func (c *Character) EquiperJambiere(jambiere item.Item) {
	c.equiperArmure(&c.JambiereEquipee, jambiere, "👖")
}

// equiperArmure fonction utilitaire pour équiper les armures
func (c *Character) equiperArmure(emplacementActuel **item.Item, nouvelItem item.Item, emoji string) {
	if *emplacementActuel != nil {
		c.Inventaire.Items = append(c.Inventaire.Items, **emplacementActuel)
	}
	*emplacementActuel = &nouvelItem
	fmt.Printf("%s Vous équipez : %s (+%d défense)\n", emoji, nouvelItem.Nom, nouvelItem.Defense)
}

// CalculerAttaqueBonus calcule le bonus d'attaque de l'équipement
func (c *Character) CalculerAttaqueBonus() int {
	bonus := 0
	if c.ArmeEquipee != nil {
		bonus += c.ArmeEquipee.Attaque
	}
	return bonus
}

// CalculerDefenseBonus calcule le bonus de défense de l'équipement
func (c *Character) CalculerDefenseBonus() int {
	bonus := 0
	
	// Nouveau système d'équipement
	if c.CasqueEquipe != nil {
		bonus += c.CasqueEquipe.Defense
	}
	if c.TorseEquipe != nil {
		bonus += c.TorseEquipe.Defense
	}
	if c.JambiereEquipee != nil {
		bonus += c.JambiereEquipee.Defense
	}
	
	return bonus
}

// === SYSTÈME DE QUÊTES AMÉLIORÉ ===

// ProposerEtAjouterQueteAvecPNJ ajoute une quête avec le PNJ donneur
func (c *Character) ProposerEtAjouterQueteAvecPNJ(nom string, recompense string, donneurPNJ string) {
	c.Quetes = append(c.Quetes, Quete{
		Nom: nom,
		Accomplie: false,
		Recompense: recompense,
		DonneurPNJ: donneurPNJ,
		Rendue: false,
	})
}

// AjouterQueteCombat ajoute une quête avec objectifs de combat spécifiques
func (c *Character) AjouterQueteCombat(nom string, donneurPNJ string, objectifs []ObjectifCombat, or int, potionsVie int, potionsMana int) {
	recompenseText := fmt.Sprintf("%d or", or)
	if potionsVie > 0 {
		recompenseText += fmt.Sprintf(", %d potions de vie", potionsVie)
	}
	if potionsMana > 0 {
		recompenseText += fmt.Sprintf(", %d potions de mana", potionsMana)
	}
	
	c.Quetes = append(c.Quetes, Quete{
		Nom: nom,
		Accomplie: false,
		Recompense: recompenseText,
		DonneurPNJ: donneurPNJ,
		Rendue: false,
		ObjectifsCombat: objectifs,
		RecompenseOr: or,
		RecompensePotionsVie: potionsVie,
		RecompensePotionsMana: potionsMana,
	})
}

// MettreAJourProgresQuete met à jour le progrès d'une quête lors d'un combat
func (c *Character) MettreAJourProgresQuete(nomMonstre string) {
	for i := range c.Quetes {
		quete := &c.Quetes[i]
		if quete.Accomplie || quete.Rendue {
			continue
		}
		
		// Vérifier les objectifs de combat
		for j := range quete.ObjectifsCombat {
			if quete.ObjectifsCombat[j].NomMonstre == nomMonstre {
				if quete.ObjectifsCombat[j].QuantiteActuelle < quete.ObjectifsCombat[j].QuantiteRequise {
					quete.ObjectifsCombat[j].QuantiteActuelle++
					fmt.Printf("ð¯ Progrès quête '%s': %s %d/%d\n", 
						quete.Nom, nomMonstre,
						quete.ObjectifsCombat[j].QuantiteActuelle,
						quete.ObjectifsCombat[j].QuantiteRequise)
					
					// Vérifier si la quête est complète
					c.verifierCompletionQuete(quete)
				}
				break
			}
		}
		}
}

// verifierCompletionQuete vérifie si tous les objectifs d'une quête sont accomplis
func (c *Character) verifierCompletionQuete(quete *Quete) {
	if len(quete.ObjectifsCombat) == 0 {
		return
	}
	
	tousAccomplis := true
	for _, objectif := range quete.ObjectifsCombat {
		if objectif.QuantiteActuelle < objectif.QuantiteRequise {
			tousAccomplis = false
			break
		}
	}
	
	if tousAccomplis && !quete.Accomplie {
		quete.Accomplie = true
		fmt.Printf("ð Quête complétée : %s !\n", quete.Nom)
		fmt.Printf("ð Retournez voir %s pour réclamer votre récompense !\n", quete.DonneurPNJ)
	}
}

// RendreQuete rend une quête à son PNJ
func (c *Character) RendreQuete(nomQuete string) bool {
	for i := range c.Quetes {
		quete := &c.Quetes[i]
		if quete.Nom == nomQuete && quete.Accomplie && !quete.Rendue {
			quete.Rendue = true
			fmt.Printf("✅ Quête rendue : %s\n", nomQuete)
			
			// Donner les récompenses
			if quete.RecompenseOr > 0 {
				c.Argent += quete.RecompenseOr
				fmt.Printf("💰 Vous recevez %d pièces d'or !\n", quete.RecompenseOr)
			}
			if quete.RecompensePotionsVie > 0 {
				c.Inventaire.Potions += quete.RecompensePotionsVie
				fmt.Printf("🧪 Vous recevez %d potions de vie !\n", quete.RecompensePotionsVie)
			}
			if quete.RecompensePotionsMana > 0 {
				c.Inventaire.PotionsMana += quete.RecompensePotionsMana
				fmt.Printf("🧿 Vous recevez %d potions de mana !\n", quete.RecompensePotionsMana)
			}
			
			// Donner récompense XP (plus généreuse pour les quêtes complexes)
			xpBonus := 50
			if len(quete.ObjectifsCombat) > 0 {
				xpBonus = 100 // Plus d'XP pour les quêtes de combat complexes
			}
			c.GagnerExperience(xpBonus)
			
			// Gérer les anciennes quêtes simples
			if quete.Recompense == "1 potion" {
				c.Inventaire.Potions++
				fmt.Println("Vous recevez 1 potion !")
			}
			
			return true
		}
	}
	return false
}

// === GESTION DE L'ÉTAT DE LA MAP ===

// SauvegarderPositionMap sauvegarde la position actuelle du joueur
func (c *Character) SauvegarderPositionMap(x, y int) {
	c.PositionX = x
	c.PositionY = y
}

// ObtenirPosition retourne la position actuelle du joueur
func (c *Character) ObtenirPosition() (int, int) {
	return c.PositionX, c.PositionY
}

// MarquerIntroEffectuee marque l'introduction comme effectuée
func (c *Character) MarquerIntroEffectuee() {
	c.IntroEffectuee = true
}

// AIntroEffectuee retourne true si l'intro a été faite
func (c *Character) AIntroEffectuee() bool {
	return c.IntroEffectuee
}

// SauvegarderEtatZone sauvegarde l'état d'une zone après modification
func (c *Character) SauvegarderEtatZone(x, y int, visitee bool, ressourcesRestantes []string, monstresRestants []MonstreState) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
	c.EtatMap.Zones[y][x] = ZoneState{
		Visitee: visitee,
		RessourcesRestantes: ressourcesRestantes,
		MonstresRestants: monstresRestants,
	}
	}
}

// ObtenirEtatZone retourne l'état sauvegardé d'une zone
func (c *Character) ObtenirEtatZone(x, y int) (ZoneState, bool) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return c.EtatMap.Zones[y][x], true
	}
	return ZoneState{}, false
}

// InitialiserEtatMap initialise l'état de la map si ce n'est pas encore fait
func (c *Character) InitialiserEtatMap() {
	// Vérifier si l'état de la map est déjà initialisé
	mapVide := true
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if len(c.EtatMap.Zones[y][x].RessourcesRestantes) > 0 || len(c.EtatMap.Zones[y][x].MonstresRestants) > 0 {
				mapVide = false
				break
			}
		}
		if !mapVide {
			break
		}
	}
	
	// Si la map est vide, on la marque comme non-initialisée pour qu'elle soit générée
	if mapVide {
		c.EtatMap = MapState{} // Réinitialiser
	}
	
	// Initialiser la zone de départ comme découverte si aucune zone n'est marquée
	aucuneZoneDecouverte := true
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if c.ZonesDecouvertes[y][x] {
				aucuneZoneDecouverte = false
				break
			}
		}
		if !aucuneZoneDecouverte {
			break
		}
	}
	
	if aucuneZoneDecouverte {
		// Marquer la zone de départ (centre) comme découverte
		c.ZonesDecouvertes[2][2] = true
	}
}

// MarquerZoneDecouverte marque une zone comme découverte
func (c *Character) MarquerZoneDecouverte(x, y int) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		if !c.ZonesDecouvertes[y][x] {
			c.ZonesDecouvertes[y][x] = true
			fmt.Printf("✨ Nouvelle zone découverte ! (%d, %d)\n", x+1, y+1)
		}
	}
}

// EstZoneDecouverte vérifie si une zone a été découverte
func (c *Character) EstZoneDecouverte(x, y int) bool {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return c.ZonesDecouvertes[y][x]
	}
	return false
}

// ObtenirNombreZonesDecouvertes retourne le nombre total de zones découvertes
func (c *Character) ObtenirNombreZonesDecouvertes() int {
	compte := 0
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if c.ZonesDecouvertes[y][x] {
				compte++
			}
		}
	}
	return compte
}

// SauvegarderRessourcesRecoltees marque les ressources d'une zone comme récoltées (gardé pour compatibilité)
func (c *Character) SauvegarderRessourcesRecoltees(x, y int) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		// Marquer cette zone comme ayant eu ses ressources récoltées
		c.EtatMap.Zones[y][x].RessourcesRestantes = []string{} // Vider les ressources
		c.EtatMap.Zones[y][x].Visitee = true
	}
}

// SauvegarderEtatZoneComplete sauvegarde l'état complet d'une zone (ressources et monstres)
func (c *Character) SauvegarderEtatZoneComplete(x, y int, ressources []string, monstres []MonstreState) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		// Marquer les ressources comme récoltées si la liste est vide
		if len(ressources) == 0 {
			c.ZonesRessourcesRecoltees[y][x] = true
		}
		
		// Marquer les monstres comme vaincus si la liste est vide
		if len(monstres) == 0 {
			c.ZonesMonstresVaincus[y][x] = true
		}
		
		// Sauvegarder aussi dans l'ancien système pour compatibilité
		c.EtatMap.Zones[y][x] = ZoneState{
			Visitee: true,
			RessourcesRestantes: ressources,
			MonstresRestants: monstres,
		}
	}
}

// ZoneRessourcesRecoltees vérifie si les ressources d'une zone ont déjà été récoltées
func (c *Character) ZoneRessourcesRecoltees(x, y int) bool {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return c.ZonesRessourcesRecoltees[y][x]
	}
	return false
}

// ZoneMonstresVaincus vérifie si les monstres d'une zone ont déjà été vaincus
func (c *Character) ZoneMonstresVaincus(x, y int) bool {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		return c.ZonesMonstresVaincus[y][x]
	}
	return false
}

// InitialiserEtatZoneSiNecessaire initialise l'état d'une zone si elle n'a pas encore été visitée
func (c *Character) InitialiserEtatZoneSiNecessaire(x, y int, ressources []string, monstres []MonstreState) {
	if x >= 0 && x < 5 && y >= 0 && y < 5 {
		état := c.EtatMap.Zones[y][x]
		// Si la zone n'a pas encore été initialisée, l'initialiser avec le contenu par défaut
		if !état.Visitee && len(état.RessourcesRestantes) == 0 && len(état.MonstresRestants) == 0 {
			c.EtatMap.Zones[y][x] = ZoneState{
				Visitee: false, // Pas encore visitée, juste initialisée
				RessourcesRestantes: ressources,
				MonstresRestants: monstres,
			}
		}
	}
}

// === UTILISATION DE POTIONS ===

// UtiliserPotion utilise une potion de vie hors combat
func (c *Character) UtiliserPotion() {
	if c.Inventaire.Potions == 0 {
		fmt.Println("❌ Vous n'avez pas de potions de vie !")
		return
	}
	
	if c.Pdv >= c.PdvMax {
		fmt.Println("❤️  Vos PV sont déjà au maximum !")
		return
	}
	
	anciensPV := c.Pdv
	c.Pdv += 50
	if c.Pdv > c.PdvMax {
		c.Pdv = c.PdvMax
	}
	c.Inventaire.Potions--
	
	pvGagnes := c.Pdv - anciensPV
	fmt.Printf("🧪 Vous utilisez une potion de vie et récupérez %d PV !\n", pvGagnes)
	fmt.Printf("PV actuels : %d/%d\n", c.Pdv, c.PdvMax)
}

// UtiliserPotionMana utilise une potion de mana hors combat
func (c *Character) UtiliserPotionMana() {
	if c.Inventaire.PotionsMana == 0 {
		fmt.Println("❌ Vous n'avez pas de potions de mana !")
		return
	}
	
	if c.Mana >= c.ManaMax {
		fmt.Println("🔮 Votre Mana est déjà au maximum !")
		return
	}
	
	ancienMana := c.Mana
	c.Mana += 50
	if c.Mana > c.ManaMax {
		c.Mana = c.ManaMax
	}
	c.Inventaire.PotionsMana--
	
	manaGagne := c.Mana - ancienMana
	fmt.Printf("🧿 Vous utilisez une potion de mana et récupérez %d Mana !\n", manaGagne)
	fmt.Printf("Mana actuel : %d/%d\n", c.Mana, c.ManaMax)
}
