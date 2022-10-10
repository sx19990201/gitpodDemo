import {
    Application,
    authProviders,
    configureWunderGraphApplication,
    cors,
    introspect,
    templates,
} from '@wundergraph/sdk';
import server from './wundergraph.server';
import operations from './wundergraph.operations';

const weather = introspect.graphql({
    apiNamespace: "weather",
    url: "https://graphql-weather-api.herokuapp.com/",
});

const db_oauthDB_0 = introspect.mysql({
    databaseURL: "mysql://root:shaoxiong123456@8.142.115.204:3306/blog", apiNamespace: "oauthDB"
});
const db_authingApi_1 = introspect.openApi({
    source: {kind: "file", filePath: "../../static/upload/oas/7ec84a92-d2ea-46bc-ac4f-618cc162799b"},
    apiNamespace: "authingApi",
    statusCodeUnions: true
});


const myApplication = new Application({
    name: 'app',
    apis: [weather, db_oauthDB_0, db_authingApi_1],
});

// configureWunderGraph emits the configuration
configureWunderGraphApplication({
    application: myApplication,
    server,
    operations,
    authentication: {
        cookieBased: {
            providers: [authProviders.openIdConnect({
                id: 'authing',
                issuer: 'https://anson.authing.cn/oidc',
                clientId: '6291fbed1bf630cb07d8d7b7',
                clientSecret: '99989891974681e5f5f255009f26522b'
            }),], authorizedRedirectUris: []
        }
    },
    authorization: {roles: ["superadmin", "admin", "user"]},
    s3UploadProvider: [],
    cors: {...cors.allowAll, allowedOrigins: ["http://localhost:3000"], allowedHeaders: [""]},
    security: {enableGraphQLEndpoint: true},
    dotGraphQLConfig: {
        hasDotWunderGraphDirectory: false,
    },
    codeGenerators: [
        {
            templates: [
                ...templates.typescript.all,
                templates.typescript.operations,
                templates.typescript.linkBuilder,
            ],
        },
        {
            templates: [
                ...templates.typescript.nextjs,
            ],
            // create-react-app expects all code to be inside /src
            path: "../../static/sdk/src",
        },
    ],
});