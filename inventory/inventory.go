package inventory

import "world_of_milousques/item"

type Inventaire struct {
	Potions int `json:"potions"`
	Item item.Item `json:"items"`
}
