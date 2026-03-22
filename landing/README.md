# TripListik Landing

Standalone landing page for the TripListik Telegram Mini App.
Built with SvelteKit + `adapter-static` — outputs a pure static site, no server needed.

## Dev

```bash
npm install
npm run dev
```

## Build

```bash
npm run build   # output → ./build/
npm run preview # preview the production build locally
```

## Deploy

The `build/` folder is a fully static site — deploy anywhere:

**Vercel** (recommended):
```bash
npx vercel --cwd .
# set output directory: build
# set framework: SvelteKit
```

**Netlify**:
- Build command: `npm run build`
- Publish directory: `build`

**Cloudflare Pages**:
- Build command: `npm run build`
- Output directory: `build`

**Railway / any static host**:
Just serve the `build/` folder as a static directory.
