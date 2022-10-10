import { configureWunderGraphOperations } from '@wundergraph/sdk';
import type { OperationsConfiguration } from './generated/wundergraph.operations';

export default configureWunderGraphOperations<OperationsConfiguration>({
	operations: { 
		defaultConfig:{authentication:{required:false}},
		queries : (config) => ({
			...config,
			caching : {enable:false,staleWhileRevalidate:60,maxAge:60,public:false},
			liveQuery: {enable:false,pollingIntervalSeconds:0},
	}),
 		mutations : (config) => ({
		...config,
	}),
       	subscriptions : (config) => ({
		...config,
	}),
	 custom: {CreateOneBlog : config => ({
				...config,
				 
			}), ProtectedWeather : config => ({
				...config,
				authentication:{
		required: true ,
	}, 
			}),  }, },
});
