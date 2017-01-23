package data

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

var materialType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions
var materialConnection *relay.GraphQLConnectionDefinitions

var Schema graphql.Schema

func init() {
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Material" {
				return GetMaterialById(strconv.Atoi(resolvedID.ID))
			}
			return nil, nil
		},
	})
}
