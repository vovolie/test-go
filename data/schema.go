package data

import (
	"strconv"
	"golang.org/x/net/context"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

var materialType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions
var materialConnection *relay.GraphQLConnectionDefinitions

var Schema graphql.Schema

func init() {

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher:func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Material" {
				return GetMaterialById(resolvedID.ID)
			}
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
		Interfaces: []*graphql.Interface{nodeDefinitions.NodeInterface},
	})

	materialConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name: "Material",
		NodeType: materialType,
	})



}
