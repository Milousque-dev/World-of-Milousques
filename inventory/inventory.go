package inventory

type Inventaire struct {
	Potions int `json:"potions"`
	Items []string `json:"items"`
}
