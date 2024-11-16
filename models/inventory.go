package models

import (
	"github.com/jinzhu/gorm"
)

type Inventory struct {
	gorm.Model

	ArticalId string  `json:"artical_id" gorm:"unique"`
	Name      string  `json:"name"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}


type Req struct{

	ArticalId 	string 	`json:"artical_id"`
	Quantity 		int 		`json:"quantity"`
	Operation 	string 	`json:"operation"`
}
