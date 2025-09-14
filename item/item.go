package item

type Item struct {
	Nom string
	Poids int
	Effet string
	Valeur int
}

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
			Effet: "Votre meilleur ami pour être accepter à Ynuv",
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
		return Item{
			Nom: nom,
			Poids: 10,
			Effet: "oops erreur",
			Valeur: 10,
		}
	}
}
