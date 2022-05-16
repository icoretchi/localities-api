package models

// swagger:parameters localities newLocality
type Locality struct {
	//swagger:ignore
	Code            int    `json:"code" bson:"code"`
	StatisticalCode int    `json:"statisticalCode" bson:"statisticalcode"`
	Name            string `json:"name" bson:"name"`
	Status          int    `json:"status" bson:"status"`
	ParentCode      int    `json:"parentCode" bson:"parentcode"`
}
