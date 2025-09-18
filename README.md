## World of Milousques - Guide Technique Détaillé

## Sommaire

1. Introduction et présentation du projet
2. Structure du Projet
5. Explication Détaillée de Chaque Fichier

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

### 🚀 `main.go` - Point d'Entrée du Programme

```go
// Package principal - tous les programmes Go commencent ici
package main

import (
    "fmt"                           // Package pour afficher du texte
    "math/rand"                     // Package pour générer des nombres aléatoires
    "os"                           // Package pour interagir avec le système d'exploitation
    "strings"                      // Package pour manipuler les chaînes de caractères
    "time"                         // Package pour gérer le temps
    
    // Import des packages personnalisés du jeu
    "world_of_milousques/character"  // Gestion des personnages
    "world_of_milousques/classe"     // Classes de personnages
    "world_of_milousques/exploration"// Système d'exploration
    "world_of_milousques/fight"      // Système de combat
    "world_of_milousques/places"     // Lieux spéciaux
    "world_of_milousques/ui"         // Interface utilisateur
    "world_of_milousques/utils"      // Fonctions utilitaires
)

// Fonction principale - première fonction exécutée
func main() {
    // Initialiser le générateur de nombres aléatoires avec l'heure actuelle
    rand.Seed(time.Now().UnixNano())  // Garantit des nombres vraiment aléatoires
    
    // Gérer le menu principal et créer/charger un personnage
    c := gererMenuPrincipal()         // c = character (personnage)
    if c == nil {                     // Si aucun personnage n'a été créé/chargé
        return                        // Quitter le programme
    }

    // Initialiser l'état de la carte du monde
    c.InitialiserEtatMap()

    // Gérer l'introduction pour un nouveau joueur ou reprendre l'aventure
    if !executerIntroductionOuReprise(c) {  // Si le tutoriel échoue
        return                               // Quitter le programme
    }
    
    // Sauvegarder le personnage avant de commencer l'exploration
    sauvegarderPersonnageAvecMessage(c, "avant de commencer l'exploration")
    
    // Lancer le système d'exploration principal
    exploration.ExplorerMap(c)      // Le joueur peut maintenant explorer le monde
}
```

### 👤 `character/character.go` - Gestion des Personnages

Ce fichier contient la **structure principale** du personnage et toutes les fonctions pour le gérer.

```go
package character

// Structure principale qui représente un personnage joueur
type Character struct {
    // Informations de base
    Nom        string               `json:"nom"`         // Nom du personnage
    Niveau     int                  `json:"niveau"`      // Niveau du personnage (1, 2, 3...)
    Pdv        int                  `json:"pdv"`         // Points de vie actuels
    Mana       int                  `json:"mana"`        // Points de mana actuels
    PdvMax     int                  `json:"pdv_max"`     // Points de vie maximum
    ManaMax    int                  `json:"mana_max"`    // Points de mana maximum
    Experience int                  `json:"experience"`  // Points d'expérience
    Argent     int                  `json:"argent"`      // Pièces d'or possédées
    
    // Équipement du personnage
    ArmeEquipee     *item.Item       `json:"arme_equipee,omitempty"`     // Arme équipée (peut être vide)
    CasqueEquipe    *item.Item       `json:"casque_equipe,omitempty"`    // Casque équipé
    TorseEquipe     *item.Item       `json:"torse_equipe,omitempty"`     // Armure de torse
    JambiereEquipee *item.Item       `json:"jambiere_equipee,omitempty"` // Jambières
    
    // Autres systèmes
    Classe     classe.Classe        `json:"classe"`      // Classe du personnage (Mage, Guerrier...)
    Inventaire inventory.Inventaire `json:"inventaire"`  // Sac du personnage avec ses objets
    Quetes     []Quete              `json:"quetes"`      // Liste des quêtes
    
    // Système de sauvegarde avancé
    IntroEffectuee bool              `json:"intro_effectuee"`      // Le tutoriel a-t-il été fait ?
    PositionX      int               `json:"position_x"`           // Position X sur la carte
    PositionY      int               `json:"position_y"`           // Position Y sur la carte
    EtatMap        MapState          `json:"etat_map"`            // État de toutes les zones visitées
    ZonesDecouvertes [5][5]bool     `json:"zones_decouvertes"`   // Grille des zones découvertes
    
    // Suivi des zones vidées
    ZonesRessourcesRecoltees [5][5]bool `json:"zones_ressources_recoltees"` // Zones sans ressources
    ZonesMonstresVaincus [5][5]bool     `json:"zones_monstres_vaincus"`     // Zones sans monstres
}

// Fonction pour créer un nouveau personnage
func InitCharacter(nom string, c classe.Classe, niveau int, pdv int, pdvmax int) Character {
    return Character{
        // Initialiser toutes les valeurs de base
        Nom:            nom,
        Niveau:         niveau,
        Pdv:            pdv,
        Mana:           c.ManaMax,     // Commencer avec le mana au maximum
        PdvMax:         c.Pvmax,       // PV maximum selon la classe
        ManaMax:        c.ManaMax,     // Mana maximum selon la classe
        Experience:     0,             // Commencer sans expérience
        Argent:         100,           // Commencer avec 100 pièces d'or
        Classe:         c,             // Assigner la classe choisie
        Inventaire:     inventory.Inventaire{}, // Inventaire vide au début
        Quetes:         []Quete{},     // Aucune quête au début
        
        // Initialiser les nouveaux champs
        IntroEffectuee: false,         // Tutoriel pas encore fait
        PositionX:      2,             // Position centrale sur la carte (milieu)
        PositionY:      2,             // Position centrale sur la carte (milieu)
        EtatMap:        MapState{},    // État de carte vide
        ZonesDecouvertes: [5][5]bool{}, // Aucune zone découverte au début
    }
}

// Fonction pour sauvegarder un personnage dans un fichier
func (c *Character) Sauvegarder() error {
    // Créer le nom du fichier de sauvegarde
    filename := "saves/" + c.Nom + ".json"  // Par exemple: "saves/Milousque.json"

    // Créer le fichier sur le disque dur
    file, err := os.Create(filename)
    if err != nil {          // Si erreur lors de la création du fichier
        return err           // Retourner l'erreur
    }
    defer file.Close()       // Fermer le fichier automatiquement à la fin

    // Encoder les données du personnage en format JSON
    encoder := json.NewEncoder(file)     // Créer un encodeur JSON
    encoder.SetIndent("", "  ")          // Formater joliment le JSON
    err = encoder.Encode(c)              // Convertir le personnage en JSON
    if err != nil {                      // Si erreur lors de l'encodage
        return err                       // Retourner l'erreur
    }

    fmt.Println("Personnage sauvegardé dans", filename)  // Confirmer la sauvegarde
    return nil                                           // Pas d'erreur = succès
}

// Fonction pour charger un personnage depuis un fichier
func Charger(nom string) (*Character, error) {
    // Construire le nom du fichier
    filename := "saves/" + nom + ".json"

    // Ouvrir le fichier de sauvegarde
    file, err := os.Open(filename)
    if err != nil {          // Si le fichier n'existe pas ou erreur
        return nil, err      // Retourner une erreur
    }
    defer file.Close()       // Fermer le fichier automatiquement

    // Créer une variable pour stocker les données chargées
    var c Character
    decoder := json.NewDecoder(file)     // Créer un décodeur JSON
    err = decoder.Decode(&c)             // Lire le JSON et remplir la structure
    if err != nil {                      // Si erreur lors du décodage
        return nil, err                  // Retourner l'erreur
    }

    // Donner de l'argent de départ aux anciens personnages (compatibilité)
    if c.Argent == 0 {                   // Si le personnage n'a pas d'argent
        c.Argent = 100                   // Lui donner 100 pièces
        fmt.Println("💰 Vous recevez 100 pièces d'or de départ !")
    }

    fmt.Println("Personnage chargé depuis", filename)  // Confirmer le chargement
    return &c, nil                                     // Retourner le personnage chargé
}
```

### ⚔️ `fight/fight.go` - Système de Combat

```go
package fight

// Structure représentant un ennemi
type Ennemi struct {
    Nom     string    // Nom du monstre (ex: "Gobelin")
    Pv      int       // Points de vie du monstre
    Attaque int       // Force d'attaque du monstre
}

// Fonction principale de combat - prend un joueur et un ennemi
func Fight(joueur *character.Character, ennemi *Ennemi) {
    tourCount := 0           // Compteur de tours de combat
    maxTours := 100          // Limite pour éviter les combats infinis

    // Boucle de combat - continue tant que les deux sont vivants
    for joueur.Pdv > 0 && ennemi.Pv > 0 && tourCount < maxTours {
        tourCount++          // Incrémenter le compteur de tours
        
        fmt.Printf("\n=== Tour %d ===\n", tourCount)  // Afficher le numéro du tour
        
        // Afficher l'interface de combat avec les statistiques
        ui.AfficherMenuCombat(
            joueur.Nom, joueur.Pdv, joueur.Classe.Pvmax, joueur.Mana, joueur.Classe.ManaMax,
            ennemi.Nom, ennemi.Pv, joueur.Classe.Sorts, joueur.Inventaire.Potions, joueur.Inventaire.PotionsMana,
        )

        // Créer la liste des actions possibles
        options := make([]string, 0)  // Slice vide pour stocker les options
        
        // Vérifier quels sorts le joueur peut utiliser
        sortUtilisable := false
        for _, s := range joueur.Classe.Sorts {  // Pour chaque sort de la classe
            // Ajouter le sort à la liste des options
            options = append(options, fmt.Sprintf("%s (Dégâts: %d, Mana: %d)", s.Nom, s.Degats, s.Cout))
            if joueur.Mana >= s.Cout {           // Si le joueur a assez de mana
                sortUtilisable = true            // Au moins un sort est utilisable
            }
        }
        
        // Ajouter les options d'utilisation de potions et de fuite
        options = append(options, fmt.Sprintf("Utiliser une potion de vie (+50 PV) (%d disponibles)", joueur.Inventaire.Potions))
        options = append(options, fmt.Sprintf("Utiliser une potion de mana (+50 Mana) (%d disponibles)", joueur.Inventaire.PotionsMana))
        options = append(options, "Fuir le combat")  // Option de fuite toujours disponible
        
        // Demander au joueur de choisir une action
        choix := utils.ScanChoice("Choisis ton action : ", options)

        // Calculer les indices des options spéciales
        optionPotionVie := len(joueur.Classe.Sorts) + 1   // Position de l'option potion de vie
        optionPotionMana := len(joueur.Classe.Sorts) + 2  // Position de l'option potion de mana
        optionFuite := len(joueur.Classe.Sorts) + 3       // Position de l'option de fuite

        // Traiter le choix du joueur
        if choix == optionFuite {                          // Si le joueur choisit de fuir
            fmt.Println("\n🏃 Vous fuyez le combat !")
            break                                          // Sortir de la boucle de combat
            
        } else if choix == optionPotionVie {               // Si utilisation d'une potion de vie
            if joueur.Inventaire.Potions > 0 {             // Si le joueur a des potions
                anciensPV := joueur.Pdv                    // Sauvegarder les PV actuels
                joueur.Pdv += 50                           // Ajouter 50 PV
                if joueur.Pdv > joueur.Classe.Pvmax {      // Ne pas dépasser le maximum
                    joueur.Pdv = joueur.Classe.Pvmax
                }
                joueur.Inventaire.Potions--                // Consommer une potion
                pvRecuperes := joueur.Pdv - anciensPV      // Calculer les PV récupérés
                fmt.Printf("🧆 Vous utilisez une potion de vie et récupérez %d PV !\n", pvRecuperes)
            } else {
                fmt.Println("⚠️  Vous n'avez pas de potion de vie !")
                continue                                   // Recommencer le tour
            }
            
        } else if choix == optionPotionMana {              // Si utilisation d'une potion de mana
            if joueur.Inventaire.PotionsMana > 0 {         // Si le joueur a des potions de mana
                ancienMana := joueur.Mana                  // Sauvegarder le mana actuel
                joueur.Mana += 50                          // Ajouter 50 mana
                if joueur.Mana > joueur.Classe.ManaMax {   // Ne pas dépasser le maximum
                    joueur.Mana = joueur.Classe.ManaMax
                }
                joueur.Inventaire.PotionsMana--           // Consommer une potion
                manaRecupere := joueur.Mana - ancienMana   // Calculer le mana récupéré
                fmt.Printf("🧙 Vous utilisez une potion de mana et récupérez %d Mana !\n", manaRecupere)
            } else {
                fmt.Println("⚠️  Vous n'avez pas de potion de mana !")
                continue                                   // Recommencer le tour
            }
            
        } else if choix >= 1 && choix <= len(joueur.Classe.Sorts) {  // Si choix d'un sort
            s := joueur.Classe.Sorts[choix-1]              // Récupérer le sort choisi
            if joueur.Mana < s.Cout {                      // Vérifier si assez de mana
                fmt.Println("⚠️  Pas assez de mana pour lancer ce sort !")
                continue                                   // Recommencer le tour
            }
            joueur.Mana -= s.Cout                          // Consommer le mana
            
            // Appliquer les bonus d'attaque de l'équipement
            bonusAttaque := joueur.CalculerAttaqueBonus()  // Calculer les bonus d'équipement
            degatsFinaux := s.Degats + bonusAttaque        // Ajouter les bonus aux dégâts
            ennemi.Pv -= degatsFinaux                      // Infliger les dégâts à l'ennemi
            
            // Afficher le résultat de l'attaque
            if bonusAttaque > 0 {
                fmt.Printf("⚔️  Tu lances %s et infliges %d dégâts (%d base + %d bonus équipement) !\n", 
                    s.Nom, degatsFinaux, s.Degats, bonusAttaque)
            } else {
                fmt.Printf("⚔️  Tu lances %s et infliges %d dégâts !\n", s.Nom, degatsFinaux)
            }
        }

        // Vérifier si l'ennemi est vaincu
        if ennemi.Pv <= 0 {
            fmt.Printf("🏆 %s est vaincu !\n", ennemi.Nom)
            break                                          // Sortir de la boucle de combat
        }

        // Tour de l'ennemi - il attaque le joueur
        bonusDefense := joueur.CalculerDefenseBonus()      // Calculer la défense du joueur
        degatsSubis := ennemi.Attaque - bonusDefense       // Calculer les dégâts après défense
        if degatsSubis < 1 {                               // Minimum 1 dégât
            degatsSubis = 1
        }
        
        joueur.Pdv -= degatsSubis                          // Le joueur perd des PV
        
        // Afficher l'attaque de l'ennemi
        if bonusDefense > 0 {
            fmt.Printf("🔴 %s t'attaque ! Tu subis %d dégâts (%d - %d défense) !\n", 
                ennemi.Nom, degatsSubis, ennemi.Attaque, bonusDefense)
        } else {
            fmt.Printf("🔴 %s t'attaque et inflige %d dégâts !\n", ennemi.Nom, degatsSubis)
        }
        
        // Petite pause pour que le joueur puisse lire
        fmt.Println("\nAppuyez sur Entrée pour continuer...")
        fmt.Scanln()  // Attendre que le joueur appuie sur Entrée
    }
    
    // Fin du combat - déterminer le résultat
    if joueur.Pdv > 0 && ennemi.Pv <= 0 {                 // Si le joueur gagne
        // Mettre à jour le progrès des quêtes
        joueur.MettreAJourProgresQuete(ennemi.Nom)
        
        // Donner de l'expérience au joueur
        xpGagne := 25 + (ennemi.Attaque * 2)              // Formule de calcul d'XP
        joueur.GagnerExperience(xpGagne)
        
    } else if joueur.Pdv <= 0 {                           // Si le joueur perd
        fmt.Println("💀 Tu as été vaincu... Game Over.")
    } else {                                              // Si fuite ou match nul
        fmt.Println("🏃 Vous avez fui le combat avec succès !")
    }
}
```

### 🌍 `world/world.go` - Génération du Monde

```go
package world

// Structure représentant un PNJ (Personnage Non-Joueur)
type PNJ struct {
    Nom       string    // Nom du PNJ (ex: "Marchand Bob")
    Dialogue  string    // Ce que dit le PNJ quand on lui parle
    Quete     string    // Nom de la quête qu'il peut donner (peut être vide)
    Recompense string   // Récompense de la quête (peut être vide)
}

// Structure représentant une zone du monde
type Zone struct {
    Nom         string          // Nom de la zone (ex: "Forêt Mystérieuse")
    Description string          // Description détaillée de la zone
    Ressources  []item.Item     // Liste des ressources qu'on peut récolter
    Monstres    []fight.Ennemi  // Liste des monstres présents
    PNJs        []PNJ          // Liste des PNJs présents
    Visitee     bool           // Est-ce que le joueur a déjà visité cette zone ?
}

// Structure représentant la position du joueur
type Position struct {
    X, Y int    // Coordonnées X et Y sur la carte
}

// Structure représentant la carte complète du monde (5x5 = 25 zones)
type Map struct {
    Zones    [5][5]Zone    // Grille de 5x5 zones
    Position Position      // Position actuelle du joueur
}

// Fonction pour créer une nouvelle carte du monde
func NewMap() *Map {
    m := &Map{
        Position: Position{X: 2, Y: 2},  // Commencer au centre de la carte
    }
    
    // Remplir toutes les zones avec du contenu
    m.initializeZones()  // Appeler la fonction d'initialisation
    
    return m  // Retourner la carte créée
}

// Fonction pour initialiser toutes les zones de la carte
func (m *Map) initializeZones() {
    // Tableaux définissant les types de chaque zone
    zoneTypes := [5][5]string{
        {"Champs", "Champs", "Transition Nord", "Mines", "Mines"},
        {"Champs", "Forêt", "Transition Ouest", "Transition Est", "Mines"},
        {"Forêt", "Forêt", "Astrab", "Rivière", "Rivière"},
        {"Forêt", "Forêt", "Transition Sud", "Rivière", "Rivière"},
        {"Forêt", "Forêt", "Rivière", "Rivière", "Rivière"},
    }
    
    // Noms spécifiques de chaque zone
    zoneNames := [5][5]string{
        {"Champs du Nord", "Grandes Cultures", "Carrefour des Vents", "Entrée des Mines", "Puits Profonds"},
        {"Terres Fertiles", "Orlière de la Forêt", "Route d'Astrab Ouest", "Route d'Astrab Est", "Mines Actives"},
        {"Cœur de la Forêt", "Clairière Sacrée", "Astrab - Capitale", "Berges Paisibles", "Confluent des Eaux"},
        {"Forêt Profonde", "Sentier des Chasseurs", "Route d'Astrab Sud", "Gues de la Rivière", "Delta Sauvage"},
        {"Bois Anciens", "Refuge des Bûcherons", "Embouchure", "Rapides Tumultueux", "Estuaire Mystérieux"},
    }
    
    // Parcourir chaque position de la grille 5x5
    for y := 0; y < 5; y++ {        // Pour chaque ligne
        for x := 0; x < 5; x++ {    // Pour chaque colonne
            zone := &m.Zones[y][x]  // Obtenir une référence à la zone actuelle
            
            // Assigner le nom et le type de la zone
            zone.Nom = zoneNames[y][x]
            zoneType := zoneTypes[y][x]
            
            // Configuration spéciale pour Astrab (la capitale)
            if zoneType == "Astrab" {
                zone.Description = "Astrab, la magnifique capitale du royaume."
                zone.Visitee = true        // Astrab est découverte dès le début
                zone.Ressources = []item.Item{}    // Pas de ressources à Astrab
                zone.Monstres = []fight.Ennemi{}   // Pas de monstres à Astrab
                zone.PNJs = []PNJ{         // PNJs spéciaux d'Astrab
                    {Nom: "Maître Karim le Marchand", Dialogue: "Bienvenue dans ma boutique !"},
                    {Nom: "Maître Forgeron Hassan", Dialogue: "Ma forge est à votre disposition !"},
                    {Nom: "Banquier Salomon", Dialogue: "La banque garde vos biens en sécurité !"},
                    {Nom: "Garde Royale", Dialogue: "Astrab est la cité la plus sûre du royaume."},
                }
            } else {
                // Pour toutes les autres zones, générer le contenu selon le type
                m.setupZoneByType(zone, zoneType)
                
                // Ajouter les PNJs spéciaux avec des quêtes aux positions correctes
                m.ajouterPNJsSpeciaux(zone, x, y)
            }
        }
    }
}

// Fonction pour configurer une zone selon son type
func (m *Map) setupZoneByType(zone *Zone, zoneType string) {
    switch zoneType {    // Selon le type de zone
    case "Champs":
        m.setupChampsZone(zone)      // Configurer comme une zone de champs
    case "Forêt":
        m.setupForetZone(zone)       // Configurer comme une zone de forêt
    case "Mines":
        m.setupMinesZone(zone)       // Configurer comme une zone de mines
    case "Rivière":
        m.setupRiviereZone(zone)     // Configurer comme une zone de rivière
    default:
        m.setupTransitionZone(zone, zoneType)  // Zone de transition
    }
}

// Fonction pour configurer une zone de champs
func (m *Map) setupChampsZone(zone *Zone) {
    zone.Description = "De vastes champs s'étendent à perte de vue."
    
    // Ajouter des ressources typiques des champs
    zone.Ressources = []item.Item{}
    for i := 0; i < 8; i++ {    // Ajouter 8 laitues
        zone.Ressources = append(zone.Ressources, item.NewItem("Laitue Vireuse"))
    }
    for i := 0; i < 7; i++ {    // Ajouter 7 blés
        zone.Ressources = append(zone.Ressources, item.NewItem("Blé"))
    }
    
    // Ajouter des monstres typiques des champs
    zone.Monstres = []fight.Ennemi{
        {Nom: "Moutmout", Pv: 80, Attaque: 25},                    // Monstre faible
        {Nom: "Moutmout", Pv: 80, Attaque: 25},                    // Un deuxième
        {Nom: "Retourneur de panneaux", Pv: 150, Attaque: 40},     // Monstre plus fort
        {Nom: "Retourneur de panneaux", Pv: 150, Attaque: 40},     // Un deuxième
    }
    
    // Pas de PNJ par défaut dans les champs
    zone.PNJs = []PNJ{}
}

// Fonction pour ajouter les PNJs spéciaux avec des quêtes
func (m *Map) ajouterPNJsSpeciaux(zone *Zone, x int, y int) {
    // PNJ à la position (3,3) - Gawr Gura dans les rivières
    if x == 3 && y == 3 {
        pnjGura := PNJ{
            Nom: "Gawr Gura",
            Dialogue: "Shaaaark ! La danse des crabe hijacob est insupportable",
            Quete: "Nettoyage des Rivières",
            Recompense: "300 or, 3 potions de vie, 3 potions de mana",
        }
        zone.PNJs = append(zone.PNJs, pnjGura)  // Ajouter le PNJ à la zone
    }
    
    // PNJ à la position (3,1) - Houshou Marine dans les champs
    if x == 3 && y == 1 {
        pnjMarine := PNJ{
            Nom: "Houshou Marine",
            Dialogue: "Les champs sont envahi, la BRUV N est dépasser ! Va apporter la démocratie",
            Quete: "Raid des Champs",
            Recompense: "300 or, 3 potions de vie, 3 potions de mana",
        }
        zone.PNJs = append(zone.PNJs, pnjMarine)
    }
    
    // PNJ à la position (1,1) - Fillian dans les mines
    if x == 1 && y == 1 {
        pnjFillian := PNJ{
            Nom: "Fillian",
            Dialogue: "Ces Kairis ont envahi mes mines ! Il faut les stopper !",
            Quete: "Répression des Kairis",
            Recompense: "300 or, 3 potions de vie, 3 potions de mana",
        }
        zone.PNJs = append(zone.PNJs, pnjFillian)
    }
    
    // PNJ à la position (1,3) - Shxtou dans la forêt
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
```

### 🎒 `inventory/inventory.go` - Système d'Inventaire

```go
package inventory

// Structure représentant l'inventaire d'un personnage
type Inventaire struct {
    Potions     int         `json:"potions"`      // Nombre de potions de vie
    PotionsMana int         `json:"potions_mana"` // Nombre de potions de mana
    Items       []item.Item `json:"items"`        // Liste des objets possédés
}

// Fonction pour ajouter des objets à l'inventaire
func (inv *Inventaire) AddItem(it item.Item, quantity int) bool {
    // Vérifier si l'inventaire a assez de place (limite de 100 objets)
    espaceDisponible := 100 - len(inv.Items)  // Calculer l'espace libre
    if quantity > espaceDisponible {           // Si on veut ajouter trop d'objets
        fmt.Printf("❌ Inventaire plein ! Vous ne pouvez ajouter que %d objets sur les %d demandés.\n", 
            espaceDisponible, quantity)
        quantity = espaceDisponible            // Limiter à l'espace disponible
    }
    
    if quantity <= 0 {    // Si plus de place ou quantité invalide
        return false      // Échec de l'ajout
    }
    
    // Ajouter les objets un par un
    for i := 0; i < quantity; i++ {
        inv.Items = append(inv.Items, it)      // Ajouter l'objet à la liste
    }
    return true          // Succès de l'ajout
}

// Fonction pour récolter des ressources dans une zone
func (inv *Inventaire) Recolter(ressources []item.Item) {
    if len(ressources) == 0 {    // Si aucune ressource à récolter
        fmt.Println("Aucune ressource à récolter ici.")
        return
    }

    // Vérifier l'espace disponible dans l'inventaire
    espaceDisponible := 100 - len(inv.Items)
    if len(ressources) > espaceDisponible {    // Si trop de ressources pour l'inventaire
        fmt.Printf("⚠️  Votre inventaire ne peut contenir que %d objets supplémentaires.\n", espaceDisponible)
        fmt.Printf("Vous ne pouvez récolter que les %d premiers objets.\n", espaceDisponible)
        ressources = ressources[:espaceDisponible]  // Limiter la liste des ressources
    }

    if len(ressources) == 0 {    // Si finalement aucune ressource ne peut être ajoutée
        fmt.Println("❌ Inventaire plein ! Impossible de récolter quoi que ce soit.")
        return
    }

    // Récolter chaque ressource
    fmt.Println("Vous récoltez :")
    for _, it := range ressources {    // Pour chaque ressource
        fmt.Printf("- %s\n", it.Nom)           // Afficher le nom de la ressource
        inv.AddItem(it, 1)                     // Ajouter 1 exemplaire à l'inventaire
    }
    fmt.Printf("✅ Votre inventaire contient maintenant %d/100 objets.\n", len(inv.Items))
}

// Fonction pour afficher le contenu de l'inventaire
func (inv *Inventaire) Afficher() {
    if len(inv.Items) == 0 {    // Si l'inventaire est vide
        fmt.Printf("🎒 Votre inventaire est vide (0/100 objets).\n")
        return
    }
    
    // Afficher l'en-tête
    fmt.Printf("🎒 === INVENTAIRE (%d/100 objets) === 🎒\n", len(inv.Items))
    
    // Afficher chaque objet avec ses propriétés
    for i, it := range inv.Items {
        fmt.Printf("%d) %s | Poids: %d | Effet: %s | Valeur: %d\n", 
            i+1, it.Nom, it.Poids, it.Effet, it.Valeur)
    }
    
    // Avertissement si l'inventaire est presque plein
    if len(inv.Items) >= 90 {
        fmt.Printf("⚠️  Attention ! Votre inventaire est presque plein (%d/100).\n", len(inv.Items))
    } else {
        fmt.Printf("💼 Espace disponible : %d objets\n", 100-len(inv.Items))
    }
}
