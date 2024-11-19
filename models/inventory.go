package models

import (
	"github.com/jinzhu/gorm"
)

type Inventory struct {
	gorm.Model

	ArticalId string  `json:"artical_id" gorm:"column:artical_id;unique"`
	Name      string  `json:"name" gorm:"column:name"`
	Quantity  int32   `json:"quantity" gorm:"column:quantity"`
	Price     float64 `json:"price" gorm:"column:price"`
}


type Req struct{

	ArticalId 	string 	`json:"artical_id"`
	Quantity 		int 		`json:"quantity"`
	Operation 	string 	`json:"operation"`
}
