package sorts

type Sorts struct {
	Nom string
	Degats int
	Cout int
}

func InitSorts(nom string, degats int, cout int) Sorts {
	return Sorts{
		Nom: nom,
		Degats: degats,
		Cout: cout,
	}
}

func GetSorts(nom string) Sorts{
	switch nom {
	case "Boule de feu":
		return Sorts{
			Nom: "Boule de feu",
			Degats: 30,
			Cout: 20,
		}
	case "Explosion":
		return Sorts{
			Nom: "Explosion",
			Degats: 50,
			Cout: 40,
		}
	case "Coup bas":
		return Sorts{
			Nom: "Coup bas",
			Degats: 25,
			Cout: 15,
		}
	case "Fourberie":
		return Sorts{
			Nom: "Fourberie",
			Degats: 10,
			Cout: 0,
		}
	case "Fracasser":
		return Sorts{
			Nom: "Fracasser",
			Degats: 20,
			Cout: 10,
		}
	case "Briser":
		return Sorts{
			Nom: "Briser",
			Degats: 40,
			Cout: 20,
		}
	default:
		return Sorts{
			Nom: nom,
			Degats: 25,
			Cout: 15,
		}
	}
}
