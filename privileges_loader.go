package auth

import (
	"context"
	"sort"
)

type PrivilegesLoader interface {
	Load(ctx context.Context, id string) ([]Privilege, error)
}

type Module struct {
	Id          string `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Name        string `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Resource    string `json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Path        string `json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Icon        string `json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty"`
	Permissions int32  `json:"permissions" gorm:"column:permissions" bson:"permissions" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty"`
	Sequence    int    `json:"sequence" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty"`
	Level       int    `json:"level" gorm:"column:level" bson:"level" dynamodbav:"level,omitempty" firestore:"level,omitempty"`
	Parent      string `json:"parent" gorm:"column:parent" bson:"parent" dynamodbav:"parent,omitempty" firestore:"parent,omitempty"`
}

func ToPrivileges(modules []Module) []Privilege {
	var menuModule []Privilege
	for _, v := range modules {
		if v.Level == 0 {
			child := make([]Privilege, 0)
			menuModule = append(menuModule,
				Privilege{
					Id:       v.Id,
					Name:     v.Name,
					Resource: v.Resource,
					Path:     v.Path,
					Icon:     v.Icon,
					Sequence: v.Sequence,
					Children: &child,
				})
		} else {
			index := findIndex(menuModule, v.Parent)
			if index != -1 {
				if findMenuModule(modules, v.Id) {
					*menuModule[index].Children = append(*menuModule[index].Children, Privilege{
						Id:       v.Id,
						Name:     v.Name,
						Resource: v.Resource,
						Path:     v.Path,
						Icon:     v.Icon,
						Sequence: v.Sequence,
					})
				}
			}
		}
	}
	return menuModule
}

func findIndex(menuModule []Privilege, key string) int {
	for i, v := range menuModule {
		if v.Id == key {
			return i
		}
	}
	return -1
}

func findMenuModule(accessModules []Module, key string) bool {
	for _, v := range accessModules {
		if v.Id == key {
			return true
		}
	}
	return false
}

func sortMenu(menu []Privilege) {
	sort.Slice(menu, func(i, j int) bool { return menu[i].Sequence < menu[j].Sequence })
	for _, v := range menu {
		sort.Slice(*v.Children, func(i, j int) bool { return (*v.Children)[i].Sequence < (*v.Children)[j].Sequence })
	}
}

func sortPrivilege(menu []Privilege) {
	sort.Slice(menu, func(i, j int) bool { return menu[i].Id < menu[j].Id })
}
