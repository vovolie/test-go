package data

import (
	"golang.org/x/net/context"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

var materialType *graphql.Object
var queryType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions
var materialsConnection *relay.GraphQLConnectionDefinitions

var Schema graphql.Schema

func init() {

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher:func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Material" {
				return GetMaterialById(resolvedID.ID)
			}
			//if resolvedID.Type == "Category" {
			//	return
			//}
			return nil, nil
		},
		TypeResolve:func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *Material:
				return materialType
			default:
				return nil
			}
		},
	})

	materialType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Material",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Material", nil),
			"category": &graphql.Field{
				Type: graphql.Int,
			},
			"cover": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"sha": &graphql.Field{
				Type: graphql.String,
			},
			"version": &graphql.Field{
				Type: graphql.String,
			},
			"mate_info": &graphql.Field{
				Type: graphql.String,
			},
			"hidden_at": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{nodeDefinitions.NodeInterface},
	})

	materialsConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name: "Material",
		NodeType: materialType,
	})

	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Material", nil),
			"material": &graphql.Field{
				Type: materialType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return GetMaterialById(idQuery)
					}
					return Material{}, nil
				},
			},
			"materials": &graphql.Field{
				Type: materialsConnection.ConnectionType,
				Args: relay.NewConnectionArgs(graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
						DefaultValue: "1",
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idCat, _ := p.Args["id"].(string)
					args := relay.NewConnectionArguments(p.Args)
					tmpMaterials, _ := GetMaterialByCategory(idCat, 10, 0)
					ms := MaterialsToSliceInterface(tmpMaterials)
					return relay.ConnectionFromArray(ms, args), nil
				},
			},
		},
		Interfaces: []*graphql.Interface{nodeDefinitions.NodeInterface},
	})

	rootType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Root",
		Fields: graphql.Fields{
			"all": &graphql.Field{
				Type: queryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetAllMaterial(10, 0)
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootType,
	})

	if err != nil {
		panic(err)
	}

}
