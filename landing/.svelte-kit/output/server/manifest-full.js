export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.svg","robots.txt"]),
	mimeTypes: {".svg":"image/svg+xml",".txt":"text/plain"},
	_: {
		client: {start:"_app/immutable/entry/start.BqZ2LhMO.js",app:"_app/immutable/entry/app.DVV5x80F.js",imports:["_app/immutable/entry/start.BqZ2LhMO.js","_app/immutable/chunks/BiYGH3oD.js","_app/immutable/chunks/CoJjOMgx.js","_app/immutable/chunks/BXYGMuBD.js","_app/immutable/entry/app.DVV5x80F.js","_app/immutable/chunks/CoJjOMgx.js","_app/immutable/chunks/DiKKL0P5.js","_app/immutable/chunks/BD95CvEF.js","_app/immutable/chunks/BXYGMuBD.js","_app/immutable/chunks/DWTmh51L.js","_app/immutable/chunks/lpPEBbVQ.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
