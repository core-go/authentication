package auth

import (
	"context"
	"sort"
)

type PrivilegesLoader interface {
	Load(ctx context.Context, id string) ([]Privilege, error)
}

type Module struct {
	Id          string  `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" sql:"id"`
	Name        string  `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty" sql:"name"`
	Resource    *string `mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource_key" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty" sql:"resource"`
	Path        *string `mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty" sql:"path"`
	Icon        *string `mapstructure:"icon" json:"icon,omitempty" gorm:"column:icon" bson:"icon,omitempty" dynamodbav:"icon,omitempty" firestore:"icon,omitempty" sql:"icon"`
	Parent      *string `mapstructure:"parent" json:"parent" gorm:"column:parent" bson:"parent" dynamodbav:"parent,omitempty" firestore:"parent,omitempty" sql:"parent"`
	Sequence    int     `mapstructure:"sequence" json:"sequence" gorm:"column:sequence" bson:"sequence" dynamodbav:"sequence,omitempty" firestore:"sequence,omitempty" sql:"sequence"`
	Permissions int32   `mapstructure:"permissions" json:"permissions" gorm:"column:permissions" bson:"permissions" dynamodbav:"permissions,omitempty" firestore:"permissions,omitempty" sql:"permissions"`
}

func OrPermissions(modules []Module) []Module {
	if modules == nil || len(modules) <= 1 {
		return modules
	}
	ms := make([]Module, 0)
	SortModulesById(modules)
	l1 := len(modules) - 1
	l := len(modules)
	for i := 0; i < l1; {
		for j := i + 1; j < l; j++ {
			if modules[i].Id == modules[j].Id {
				modules[i].Permissions = modules[i].Permissions | modules[j].Permissions
				if j == l1 {
					ms = append(ms, modules[i])
					i = l1 + 3
					break
				}
			} else {
				ms = append(ms, modules[i])
				i = j
			}
		}
	}
	return ms
}
func ToPrivileges(modules []Module) []Privilege {
	var menuModule []Privilege
	SortModulesById(modules) // sort by id
	root := FindRootModules(modules)
	for _, v := range root {
		par := Privilege{
			Id:          v.Id,
			Name:        v.Name,
			Sequence:    v.Sequence,
			Permissions: v.Permissions,
		}
		if v.Resource != nil {
			par.Resource = *v.Resource
		}
		if v.Path != nil {
			par.Path = *v.Path
		}
		if v.Icon != nil {
			par.Icon = *v.Icon
		}
		var child []Privilege
		for i := 0; i < len(modules); i++ {
			if modules[i].Parent != nil && v.Id == *modules[i].Parent {
				item := modules[i]
				sp := Privilege{
					Id:          item.Id,
					Name:        item.Name,
					Sequence:    item.Sequence,
					Permissions: item.Permissions,
				}
				if item.Resource != nil {
					sp.Resource = *item.Resource
				}
				if item.Path != nil {
					sp.Path = *item.Path
				}
				if item.Icon != nil {
					sp.Icon = *item.Icon
				}
				child = append(child, sp)
			}
		}
		par.Children = &child
		menuModule = append(menuModule, par)
	}
	SortPrivileges(menuModule)
	return menuModule
}

func ToPrivilegesWithNoSequence(modules []Module) []Privilege {
	var menuModule []Privilege
	SortModulesById(modules) // sort by id
	root := FindRootModules(modules)
	for _, v := range root {
		par := Privilege{
			Id:          v.Id,
			Name:        v.Name,
			Sequence:    v.Sequence,
			Permissions: v.Permissions,
		}
		if v.Resource != nil {
			par.Resource = *v.Resource
		}
		if v.Path != nil {
			par.Path = *v.Path
		}
		if v.Icon != nil {
			par.Icon = *v.Icon
		}
		var child []Privilege
		for i := 0; i < len(modules); i++ {
			if modules[i].Parent != nil && v.Id == *modules[i].Parent {
				item := modules[i]
				sp := Privilege{
					Id:          item.Id,
					Name:        item.Name,
					Sequence:    item.Sequence,
					Permissions: item.Permissions,
				}
				if item.Resource != nil {
					sp.Resource = *item.Resource
				}
				if item.Path != nil {
					sp.Path = *item.Path
				}
				if item.Icon != nil {
					sp.Icon = *item.Icon
				}
				child = append(child, sp)
			}
		}
		par.Children = &child
		menuModule = append(menuModule, par)
	}
	SortPrivileges(menuModule)
	for j := 0; j < len(menuModule); j++ {
		menuModule[j].Sequence = 0
		child := *menuModule[j].Children
		if child != nil {
			for x, _ := range child {
				child[x].Sequence = 0
			}
		}
	}
	return menuModule
}

func FindRootModules(sortModules []Module) []Module {
	var root []Module
	for _, module := range sortModules {
		if *module.Parent == "" {
			root = append(root, module)
		}
	}
	return root
}

func SortPrivileges(menu []Privilege) {
	sort.Slice(menu, func(i, j int) bool { return menu[i].Sequence < menu[j].Sequence })
	for _, v := range menu {
		sort.Slice(*v.Children, func(i, j int) bool { return (*v.Children)[i].Sequence < (*v.Children)[j].Sequence })
	}
}
func SortPrivilegesById(menu []Privilege) {
	sort.Slice(menu, func(i, j int) bool { return menu[i].Id < menu[j].Id })
}

func SortModulesById(modulePath []Module) {
	sort.Slice(modulePath, func(i, j int) bool { return modulePath[i].Id < modulePath[j].Id })
}
