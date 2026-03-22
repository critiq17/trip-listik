import * as universal from '../entries/pages/_page.ts.js';

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/+page.ts";
export const imports = ["_app/immutable/nodes/2.B2M-dUSn.js","_app/immutable/chunks/BD95CvEF.js","_app/immutable/chunks/CoJjOMgx.js","_app/immutable/chunks/BXYGMuBD.js","_app/immutable/chunks/DiKKL0P5.js","_app/immutable/chunks/DWTmh51L.js","_app/immutable/chunks/lpPEBbVQ.js","_app/immutable/chunks/2FgVoid_.js"];
export const stylesheets = ["_app/immutable/assets/2.DQuN0ne_.css"];
export const fonts = [];
