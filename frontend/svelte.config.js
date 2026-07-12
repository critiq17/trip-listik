import adapter from '@sveltejs/adapter-auto';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter(),
		version: {
			// Telegram keeps the WebView alive for hours; a stale session keeps
			// navigating client-side and mixes old and new bundles. Polling lets
			// SvelteKit detect a new deploy and full-reload on the next navigation.
			pollInterval: 60000
		}
	}
};

export default config;
