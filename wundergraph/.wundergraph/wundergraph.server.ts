import {GraphQLObjectType, GraphQLSchema, GraphQLString} from 'graphql';
import {configureWunderGraphServer} from '@wundergraph/sdk';
import type {AuthenticationResponse, HooksConfig} from './generated/wundergraph.hooks';
import type {InternalClient} from './generated/wundergraph.internal.client';

   

export default configureWunderGraphServer<HooksConfig, InternalClient>(() => ({
   hooks:{  authentication: {
				
            },  }, graphqlServers: [{
			apiNamespace: "gql",
			serverName: "gql",
			schema: new GraphQLSchema({
				query: new GraphQLObjectType({
					name: 'RootQueryType',
					fields: {
						hello: {
							type: GraphQLString,
							resolve() {
								return 'world';
							},
						},
					},
				}),
			}),
		},] 

}));
