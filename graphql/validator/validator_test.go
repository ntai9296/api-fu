package validator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ccbrown/apifu/graphql/parser"
	"github.com/ccbrown/apifu/graphql/schema"
)

var nodeType = &schema.InterfaceType{
	Name: "Node",
	Fields: map[string]*schema.FieldDefinition{
		"id": &schema.FieldDefinition{
			Type: schema.IDType,
		},
	},
}

var unionMemberType = &schema.InterfaceType{
	Name: "UnionMember",
	Fields: map[string]*schema.FieldDefinition{
		"scalar": &schema.FieldDefinition{
			Type: schema.StringType,
		},
	},
}

var objectType = &schema.ObjectType{
	Name: "Object",
}

func init() {
	objectType.Fields = map[string]*schema.FieldDefinition{
		"node": &schema.FieldDefinition{
			Type: nodeType,
			Arguments: map[string]*schema.InputValueDefinition{
				"id": &schema.InputValueDefinition{
					Type: schema.NewNonNullType(schema.IDType),
				},
			},
		},
		"resource": &schema.FieldDefinition{
			Type: &schema.ObjectType{
				Name: "Resource",
				Fields: map[string]*schema.FieldDefinition{
					"id": &schema.FieldDefinition{
						Type: schema.IDType,
					},
				},
				ImplementedInterfaces: []*schema.InterfaceType{nodeType},
			},
		},
		"objects": &schema.FieldDefinition{
			Type: schema.NewListType(objectType),
		},
		"object": &schema.FieldDefinition{
			Type: objectType,
			Arguments: map[string]*schema.InputValueDefinition{
				"object": &schema.InputValueDefinition{
					Type: &schema.InputObjectType{
						Name: "ObjectInput",
						Fields: map[string]*schema.InputValueDefinition{
							"defaultedString": &schema.InputValueDefinition{
								Type:         schema.NewNonNullType(schema.StringType),
								DefaultValue: "foo",
							},
							"requiredString": &schema.InputValueDefinition{
								Type: schema.NewNonNullType(schema.StringType),
							},
						},
					},
				},
			},
		},
		"object2": &schema.FieldDefinition{
			Type: &schema.ObjectType{
				Name: "Object2",
				Fields: map[string]*schema.FieldDefinition{
					"scalar": &schema.FieldDefinition{
						Type: schema.StringType,
					},
				},
			},
		},
		"interface": &schema.FieldDefinition{
			Type: &schema.InterfaceType{
				Name: "Interface",
				Fields: map[string]*schema.FieldDefinition{
					"scalar": &schema.FieldDefinition{
						Type: schema.StringType,
					},
				},
			},
		},
		"union": &schema.FieldDefinition{
			Type: &schema.UnionType{
				Name: "Union",
				MemberTypes: []schema.NamedType{
					&schema.ObjectType{
						Name: "UnionObjectA",
						Fields: map[string]*schema.FieldDefinition{
							"a": &schema.FieldDefinition{
								Type: schema.StringType,
							},
							"scalar": &schema.FieldDefinition{
								Type: schema.StringType,
							},
						},
						ImplementedInterfaces: []*schema.InterfaceType{unionMemberType},
					},
					&schema.ObjectType{
						Name: "UnionObjectB",
						Fields: map[string]*schema.FieldDefinition{
							"b": &schema.FieldDefinition{
								Type: schema.StringType,
							},
							"scalar": &schema.FieldDefinition{
								Type: schema.StringType,
							},
						},
						ImplementedInterfaces: []*schema.InterfaceType{unionMemberType},
					},
				},
			},
		},
		"scalar": &schema.FieldDefinition{
			Type: schema.StringType,
		},
	}
}

func validateSource(t *testing.T, src string) []*Error {
	s, err := schema.New(&schema.SchemaDefinition{
		Query: objectType,
		DirectiveDefinitions: map[string]*schema.DirectiveDefinition{
			"include": schema.IncludeDirective,
			"skip":    schema.SkipDirective,
		},
	})
	require.NoError(t, err)

	doc, errs := parser.ParseDocument([]byte(src))
	require.Empty(t, errs)
	require.NotNil(t, doc)
	return ValidateDocument(doc, s)
}
