package mail

import (
	"github.com/core-go/auth"
	"github.com/core-go/mail"
)

type AuthMailConfig struct {
	Secret   string              `mapstructure:"secret" json:"secret,omitempty" gorm:"column:secret" bson:"secret,omitempty" dynamodbav:"secret,omitempty" firestore:"secret,omitempty"`
	Expires  int64               `mapstructure:"expires" json:"expires,omitempty" gorm:"column:expires" bson:"expires,omitempty" dynamodbav:"expires,omitempty" firestore:"expires,omitempty"`
	Schema   auth.SchemaConfig   `mapstructure:"schema" json:"schema,omitempty" gorm:"column:schema" bson:"schema,omitempty" dynamodbav:"schema,omitempty" firestore:"schema,omitempty"`
	Template mail.TemplateConfig `mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"`
}
