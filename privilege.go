package auth

type Privilege struct {
	Id          string       `json:"id,omitempty" bson:"_id,omitempty" gorm:"column:id"`
	Name        string       `json:"name,omitempty" bson:"name,omitempty" gorm:"column:name"`
	ResourceKey string       `json:"resourceKey,omitempty" bson:"resourceKey,omitempty" gorm:"column:resourcekey"`
	Path        string       `json:"path,omitempty" bson:"path,omitempty" gorm:"column:path"`
	Icon        string       `json:"icon,omitempty" bson:"icon,omitempty" gorm:"column:icon"`
	Sequence    int          `json:"sequence,omitempty" bson:"sequence,omitempty" gorm:"column:sequence"`
	Children    *[]Privilege `json:"children,omitempty" bson:"children,omitempty" gorm:"column:children"`
}
