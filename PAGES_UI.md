# Wandr Mini App — Svelte Migration & Completion Guide

## Overview

This document maps the existing HTML design files to their responsibilities, defines the full set of screens needed for the Wandr Telegram Mini App, and provides a step-by-step guide for converting everything into Svelte components with production-grade animations.

---

## Existing HTML Files

### 1. `feed.html` — Explore / Feed Screen
**Responsibility:** The app's home screen. Shows a greeting header, horizontal pill filter tabs (All / Friends / Popular / Nearby), and a scrollable vertical list of public trip cards. Each card is a full-width image with a gradient overlay, destination name, date + member count, and a teal vote pill.

**Key UI patterns:**
- Sticky header with avatar + search icon
- Horizontal scrollable filter chips (hide-scrollbar)
- Full-bleed image trip cards with dark gradient overlay
- Bottom navigation bar (Feed · Trips · Create · Explore · Profile)

---

### 2. `trip-detail.html` — Trip Detail Screen
**Responsibility:** Full detail view of a single trip. Opens with a full-bleed hero image (parallax on scroll), floating back/share buttons, then sticky tab navigation (Details · Members · Photos · Discussion). The Details tab shows status badge, dates, description, join button, a map snippet, and quick stats grid.

**Key UI patterns:**
- Full-bleed hero 45vh with gradient overlay
- Frosted glass floating action buttons (backdrop-filter: blur)
- Sticky tab bar with active underline indicator
- Full-width gradient CTA button ("Join Trip")
- Route map image placeholder
- Stats grid (2-col)

---

### 3. `create-trip.html` — Create Trip Screen (Step 1 of 3)
**Responsibility:** Multi-step trip creation form. Step 1 covers trip name (minimal underline input), destination search (with flag emoji + autocomplete), and cover photo upload (dashed upload zone with camera icon and background preview).

**Key UI patterns:**
- Progress bar header (step N of 3 + percentage)
- Underline-only text inputs (no box)
- Destination field with emoji flag prefix
- Dashed upload area with hover animation
- Fixed footer with gradient "Next Step" CTA

---

### 4. `profile.html` — Profile Screen
**Responsibility:** User's personal profile. Shows avatar with gradient ring border, username, edit button, World Explored card (SVG map + animated pulsing dots + 3.5% counter), stats grid (Trips · Countries · Cities · Friends), and scrollable travel history list.

**Key UI patterns:**
- Avatar with linear-gradient ring
- World Explored card with glowing pulse dots on map
- 4-col stats grid with teal numbers
- Travel history as vertical item list (thumbnail + title + location + chevron)
- Settings icon in header

---

## Missing Screens (Must Be Built)

| Screen | Route | Description |
|---|---|---|
| **My Trips** | `/trips` | List of trips the user owns or joined. Tabs: Upcoming / Past / Drafts. |
| **Trip Detail — Members Tab** | `/trips/:id/members` | Member avatar list, owner crown badge, pending join requests (owner view). |
| **Trip Detail — Photos Tab** | `/trips/:id/photos` | Masonry/2-col photo grid, tap-to-fullscreen lightbox, upload FAB. |
| **Trip Detail — Discussion Tab** | `/trips/:id/discussion` | Comment thread, chat-bubble style, pinned input bar at bottom. |
| **Create Trip Step 2** | `/create/step-2` | Date range picker (calendar), Public/Private sliding toggle. |
| **Create Trip Step 3** | `/create/step-3` | Friend invite search with avatar chips, trip description textarea. |
| **Notifications / Inbox** | `/inbox` | Join requests, approvals, new comments, trip invites. |
| **Explore / Discover** | `/explore` | Discovery feed: popular/trending trips, search, country filter. |
| **Fullscreen Photo Viewer** | (modal overlay) | Pinch-to-zoom, swipe between photos, close gesture. |
| **Onboarding / Auth** | `/auth` | Telegram WebApp auth splash, animated logo, loading state. |
| **Empty States** | (inline) | No trips yet, no comments, no photos — illustrated empty states. |

---

## Project Setup

```bash
npm create svelte@latest wandr-mini-app
cd wandr-mini-app
npm install

# Animation & utility libraries
npm install motion         # Motion One — performant animations
npm install @sveltejs/kit
npm install clsx           # Conditional classnames

# Telegram WebApp SDK
npm install @twa-dev/sdk
```

### `app.html` — Telegram WebApp Init

```html
<!-- src/app.html -->
<script src="https://telegram.org/js/telegram-web-app.js"></script>
```

### `src/lib/telegram.js` — Auth Helper

```javascript
import WebApp from '@twa-dev/sdk';

export function getTelegramInitData() {
  return WebApp.initData;
}

export async function authenticateWithBackend() {
  const res = await fetch('/v1/auth/telegram', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ initData: WebApp.initData })
  });
  const { token } = await res.json();
  return token;
}
```

---

## CSS Design Tokens

Create `src/lib/styles/tokens.css` and import globally:

```css
:root {
  --bg-primary:     #042444;
  --bg-secondary:   #0d3460;
  --surface:        #657b9f;
  --accent:         #2092ba;
  --accent-glow:    #7aeaf4;
  --accent-grad:    linear-gradient(135deg, #2092ba, #7aeaf4);
  --text-primary:   #ffffff;
  --text-secondary: #b8c9d5;
  --text-muted:     #657b9f;
  --danger:         #e05c7a;
  --success:        #4ecdc4;

  --radius-sm:   8px;
  --radius-md:   12px;
  --radius-lg:   16px;
  --radius-xl:   24px;
  --radius-pill: 9999px;

  --font-display: 'Plus Jakarta Sans', sans-serif;
}
```

---

## Animation Architecture

Use **Motion One** (`motion` package) for all JS-driven animations. Use CSS `@keyframes` for looping effects (pulse, shimmer, glow).

### Page Transitions — `src/lib/transitions.js`

```javascript
import { animate, stagger } from 'motion';

// Slide pages in from the right (iOS-style)
export function pageEnter(node) {
  animate(node,
    { x: ['24px', '0px'], opacity: [0, 1] },
    { duration: 0.28, easing: [0.25, 0.46, 0.45, 0.94] }
  );
}

// Staggered list items fade up on mount
export function listEnter(containerNode) {
  const items = containerNode.querySelectorAll('[data-list-item]');
  animate(items,
    { y: ['16px', '0px'], opacity: [0, 1] },
    { duration: 0.35, delay: stagger(0.05), easing: 'ease-out' }
  );
}

// Count-up number animation
export function countUp(node, { target, duration = 1.2 }) {
  let start = null;
  const step = (timestamp) => {
    if (!start) start = timestamp;
    const progress = Math.min((timestamp - start) / (duration * 1000), 1);
    node.textContent = Math.floor(progress * target);
    if (progress < 1) requestAnimationFrame(step);
  };
  requestAnimationFrame(step);
}
```

### Reusable Svelte Animation Actions

```javascript
// src/lib/actions/animate.js
import { animate } from 'motion';

// use:fadeUp on any element
export function fadeUp(node, { delay = 0, duration = 0.4 } = {}) {
  animate(node,
    { y: ['20px', '0px'], opacity: [0, 1] },
    { duration, delay, easing: 'ease-out' }
  );
}

// use:scalePress — tap/click scale feedback
export function scalePress(node) {
  const down = () => animate(node, { scale: 0.96 }, { duration: 0.1 });
  const up   = () => animate(node, { scale: 1.0  }, { duration: 0.15 });
  node.addEventListener('pointerdown', down);
  node.addEventListener('pointerup', up);
  node.addEventListener('pointerleave', up);
  return {
    destroy() {
      node.removeEventListener('pointerdown', down);
      node.removeEventListener('pointerup', up);
      node.removeEventListener('pointerleave', up);
    }
  };
}
```

---

## Component Breakdown

### File Structure

```
src/
  routes/
    +layout.svelte              <- BottomNav + page wrapper
    +page.svelte                <- Feed (index)
    trips/
      +page.svelte              <- My Trips
      [id]/
        +page.svelte            <- Trip Detail (tabs)
    create/
      +page.svelte              <- Create Step 1
      step-2/+page.svelte
      step-3/+page.svelte
    profile/+page.svelte
    explore/+page.svelte
    inbox/+page.svelte
    auth/+page.svelte
  lib/
    components/
      TripCard.svelte
      BottomNav.svelte
      HeroImage.svelte
      TabBar.svelte
      StatsGrid.svelte
      WorldMap.svelte
      CommentBubble.svelte
      PhotoGrid.svelte
      AvatarStack.svelte
      ProgressBar.svelte
      FilterChips.svelte
      SkeletonCard.svelte
      EmptyState.svelte
      BottomSheet.svelte
    actions/
      animate.js
    transitions.js
    telegram.js
    styles/
      tokens.css
```

---

## Component Specs

### `TripCard.svelte`

```svelte
<script>
  import { scalePress } from '$lib/actions/animate';
  export let trip; // { title, coverUrl, dateRange, memberCount, votes }
</script>

<div class="card" use:scalePress on:click>
  <img src={trip.coverUrl} alt={trip.title} loading="lazy" />
  <div class="overlay" />
  <footer>
    <div class="info">
      <h3>{trip.title}</h3>
      <p>{trip.dateRange} · {trip.memberCount} members</p>
    </div>
    <span class="votes">👍 {trip.votes}</span>
  </footer>
</div>

<style>
  .card {
    position: relative;
    border-radius: var(--radius-lg);
    overflow: hidden;
    aspect-ratio: 16/9;
    cursor: pointer;
  }
  img { width: 100%; height: 100%; object-fit: cover; }
  .overlay {
    position: absolute; inset: 0;
    background: linear-gradient(to top, #042444ee 0%, transparent 60%);
  }
  footer {
    position: absolute; bottom: 0; left: 0; right: 0;
    padding: 1.25rem;
    display: flex; justify-content: space-between; align-items: flex-end;
  }
  h3 { color: white; font-size: 1.5rem; font-weight: 700; }
  p  { color: #b8c9d5; font-size: 0.875rem; }
  .votes {
    background: linear-gradient(135deg, #2092ba, #1a7a9c);
    color: white; font-weight: 700; font-size: 0.875rem;
    padding: 0.4rem 1rem; border-radius: var(--radius-pill);
  }
</style>
```

---

### `BottomNav.svelte`

```svelte
<script>
  import { page } from '$app/stores';
  import { animate } from 'motion';

  const tabs = [
    { href: '/',        icon: 'home',           label: 'Feed'    },
    { href: '/trips',   icon: 'map',            label: 'Trips'   },
    { href: '/create',  icon: 'add',            label: 'Create', fab: true },
    { href: '/explore', icon: 'explore',        label: 'Explore' },
    { href: '/profile', icon: 'account_circle', label: 'Profile' },
  ];

  function onTabClick(e, tab) {
    const icon = e.currentTarget.querySelector('.icon');
    animate(icon, { y: ['0px', '-4px', '0px'] }, { duration: 0.3, easing: 'ease-out' });
  }
</script>

<nav>
  {#each tabs as tab}
    {@const active = $page.url.pathname === tab.href}
    <a href={tab.href} class:active class:fab={tab.fab}
       on:click={(e) => onTabClick(e, tab)}>
      {#if tab.fab}
        <div class="fab-circle">
          <span class="material-symbols-outlined">{tab.icon}</span>
        </div>
        <span class="fab-label">Create</span>
      {:else}
        <span class="material-symbols-outlined icon"
              style="font-variation-settings: 'FILL' {active ? 1 : 0}">
          {tab.icon}
        </span>
        <span class="label" class:active>{tab.label}</span>
      {/if}
    </a>
  {/each}
</nav>

<style>
  nav {
    position: fixed; bottom: 0; left: 0; right: 0; z-index: 100;
    display: flex; align-items: center; justify-content: space-around;
    padding: 0.75rem 1rem env(safe-area-inset-bottom);
    background: rgba(4, 36, 68, 0.92);
    backdrop-filter: blur(16px);
    border-top: 1px solid rgba(255,255,255,0.06);
  }
  a {
    display: flex; flex-direction: column; align-items: center; gap: 2px;
    color: var(--text-muted); text-decoration: none; flex: 1;
  }
  a.active { color: var(--accent); }
  .label { font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; }
  .fab-circle {
    width: 48px; height: 48px; border-radius: 50%;
    background: var(--accent-grad);
    display: flex; align-items: center; justify-content: center;
    color: white; margin-top: -24px;
    box-shadow: 0 4px 20px rgba(32, 146, 186, 0.4);
  }
  .fab-label { font-size: 10px; color: var(--text-muted); font-weight: 700; text-transform: uppercase; }
</style>
```

---

### `WorldMap.svelte`

```svelte
<script>
  import { onMount } from 'svelte';
  import { animate, stagger } from 'motion';

  export let visitedCountries = []; // array of ISO alpha-2 codes e.g. ['JP', 'IS', 'IT']
  export let worldPercent = 3.5;

  let countEl;

  onMount(() => {
    // Count-up animation for the percentage
    let start = performance.now();
    const duration = 1400;
    const tick = (now) => {
      const p = Math.min((now - start) / duration, 1);
      const eased = 1 - Math.pow(1 - p, 3); // ease-out cubic
      countEl.textContent = (eased * worldPercent).toFixed(1) + '%';
      if (p < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);

    // Light up visited country SVG paths
    const allPaths = document.querySelectorAll('[data-country]');
    const visited = [...allPaths].filter(p =>
      visitedCountries.includes(p.dataset.country));
    animate(visited,
      { fill: ['#0d3460', '#2092ba'] },
      { duration: 0.6, delay: stagger(0.08), easing: 'ease-out' }
    );
  });
</script>

<div class="map-card">
  <div class="map-header">
    <div>
      <h3>World Explored</h3>
      <p>You've seen more than 92% of users</p>
    </div>
    <span class="percent" bind:this={countEl}>0.0%</span>
  </div>
  <div class="map-container">
    <!-- Drop in an inline SVG world map here.
         Each country path needs data-country="XX" (ISO alpha-2).
         Free source: https://simplemaps.com/static/svg/world.svg -->
    <slot />
  </div>
</div>

<style>
  .map-card {
    background: rgba(13, 52, 96, 0.5);
    border-radius: var(--radius-xl);
    padding: 1.5rem;
  }
  .map-header {
    display: flex; justify-content: space-between; align-items: flex-start;
    margin-bottom: 1rem;
  }
  h3 { color: white; font-size: 1.1rem; font-weight: 700; }
  p  { color: var(--text-secondary); font-size: 0.75rem; margin-top: 2px; }
  .percent { color: var(--accent); font-size: 2.25rem; font-weight: 900; }
  .map-container { width: 100%; height: 8rem; }
  :global([data-country]) { fill: #0a1a2e; transition: fill 0.3s; }
  :global([data-country].visited) {
    fill: var(--accent);
    filter: drop-shadow(0 0 4px rgba(32,146,186,0.5));
  }
</style>
```

---

### `SkeletonCard.svelte` — Loading Shimmer

```svelte
<div class="skeleton" />

<style>
  .skeleton {
    width: 100%;
    aspect-ratio: 16/9;
    border-radius: var(--radius-lg);
    background: linear-gradient(90deg,
      #0d3460 0%, #1a4a7a 40%, #0d3460 80%);
    background-size: 200% 100%;
    animation: shimmer 1.6s infinite;
  }
  @keyframes shimmer {
    0%   { background-position: -200% 0; }
    100% { background-position:  200% 0; }
  }
</style>
```

---

### `BottomSheet.svelte` — Replaces All Modals

```svelte
<script>
  import { fly, fade } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  export let open = false;
  export let onClose = () => {};
</script>

{#if open}
  <div class="backdrop" on:click={onClose}
       transition:fade={{ duration: 200 }} />
  <div class="sheet"
       transition:fly={{ y: 300, duration: 320, easing: cubicOut }}>
    <div class="handle" />
    <slot />
  </div>
{/if}

<style>
  .backdrop {
    position: fixed; inset: 0; z-index: 200;
    background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);
  }
  .sheet {
    position: fixed; bottom: 0; left: 0; right: 0; z-index: 201;
    background: #0d3460;
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    padding: 1rem 1.5rem env(safe-area-inset-bottom);
    max-height: 85vh; overflow-y: auto;
  }
  .handle {
    width: 40px; height: 4px; border-radius: 2px;
    background: rgba(255,255,255,0.2);
    margin: 0 auto 1.5rem;
  }
</style>
```

---

## Feed Page with Animations

```svelte
<!-- src/routes/+page.svelte -->
<script>
  import { onMount, tick } from 'svelte';
  import { listEnter } from '$lib/transitions';
  import TripCard from '$lib/components/TripCard.svelte';
  import FilterChips from '$lib/components/FilterChips.svelte';
  import SkeletonCard from '$lib/components/SkeletonCard.svelte';

  let trips = [];
  let loading = true;
  let listEl;

  const filters = ['All', 'Friends', 'Popular', 'Nearby'];
  let activeFilter = 'All';

  onMount(async () => {
    trips = await fetch('/v1/feed').then(r => r.json());
    loading = false;
    await tick();
    listEnter(listEl);
  });
</script>

<header>
  <img src="/avatar.jpg" alt="you" class="avatar" />
  <div class="title-block">
    <h1>Explore Trips</h1>
    <p>Hey, Alex 👋</p>
  </div>
  <button class="icon-btn">
    <span class="material-symbols-outlined">search</span>
  </button>
</header>

<FilterChips {filters} bind:active={activeFilter} />

<main bind:this={listEl} class="feed">
  {#if loading}
    {#each Array(3) as _}
      <SkeletonCard />
    {/each}
  {:else}
    {#each trips as trip}
      <div data-list-item>
        <TripCard {trip} />
      </div>
    {/each}
  {/if}
</main>
```

---

## Create Trip — Multi-Step with Animated Progress

```svelte
<!-- src/routes/create/+page.svelte -->
<script>
  import { animate } from 'motion';
  import { goto } from '$app/navigation';
  import { fly } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';

  let step = 1;
  const totalSteps = 3;
  let progressEl;

  $: percent = (step / totalSteps) * 100;

  async function nextStep() {
    await animate(progressEl, { width: `${percent}%` },
      { duration: 0.4, easing: [0.34, 1.56, 0.64, 1] }).finished;
    step++;
    if (step > totalSteps) goto('/trips');
  }
</script>

<header class="create-header">
  <button on:click={() => history.back()} class="back-btn">
    <span class="material-symbols-outlined">arrow_back</span>
  </button>
  <span class="create-label">New Journey</span>
  <div />
</header>

<div class="progress-track">
  <span class="step-label">Step {step} of {totalSteps}</span>
  <span class="pct">{Math.round(percent)}%</span>
</div>
<div class="progress-bar">
  <div class="progress-fill" bind:this={progressEl}
       style="width: {(1/totalSteps)*100}%" />
</div>

{#key step}
  <div in:fly={{ x: 32, duration: 260, easing: cubicOut }}
       out:fly={{ x: -32, duration: 200 }}>
    {#if step === 1}
      <!-- Trip name, destination, cover photo -->
    {:else if step === 2}
      <!-- Date range picker, visibility toggle -->
    {:else if step === 3}
      <!-- Friend invite, description -->
    {/if}
  </div>
{/key}

<footer class="create-footer">
  <button class="cta-btn" on:click={nextStep}>
    {step < totalSteps ? 'Next Step' : 'Create Trip'}
    <span class="material-symbols-outlined">arrow_forward</span>
  </button>
</footer>
```

---

## Auth Flow — Telegram WebApp

```svelte
<!-- src/routes/auth/+page.svelte -->
<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { animate } from 'motion';
  import { authenticateWithBackend } from '$lib/telegram';

  let logoEl;

  onMount(async () => {
    await animate(logoEl,
      { scale: [0.6, 1], opacity: [0, 1] },
      { duration: 0.6, easing: [0.34, 1.56, 0.64, 1] }
    ).finished;

    const token = await authenticateWithBackend();
    localStorage.setItem('wandr_token', token);
    goto('/');
  });
</script>

<div class="auth-screen">
  <div class="logo" bind:this={logoEl}>
    <span class="logo-icon">✈️</span>
    <h1>Wandr</h1>
  </div>
  <p class="connecting">Connecting to Telegram…</p>
  <div class="spinner" />
</div>

<style>
  .auth-screen {
    min-height: 100dvh;
    display: flex; flex-direction: column;
    align-items: center; justify-content: center;
    background: var(--bg-primary);
    gap: 1rem;
  }
  .logo { text-align: center; }
  .logo-icon { font-size: 3rem; }
  h1 { font-size: 2.5rem; font-weight: 900; color: white; }
  .connecting { color: var(--text-secondary); font-size: 0.875rem; }
  .spinner {
    width: 28px; height: 28px; border-radius: 50%;
    border: 2px solid rgba(32,146,186,0.2);
    border-top-color: var(--accent);
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>
```

---

## Global CSS Animations

Add to `src/app.css`:

```css
/* Pulse glow for world map dots */
@keyframes pulseGlow {
  0%, 100% { opacity: 1; transform: scale(1); box-shadow: 0 0 6px #2092ba; }
  50%       { opacity: 0.6; transform: scale(1.4); box-shadow: 0 0 14px #7aeaf4; }
}
.pulse-dot { animation: pulseGlow 2s ease-in-out infinite; }

/* Shimmer for skeleton cards */
@keyframes shimmer {
  0%   { background-position: -200% 0; }
  100% { background-position:  200% 0; }
}

/* Accent button breathing glow */
@keyframes glowPulse {
  0%, 100% { box-shadow: 0 4px 20px rgba(32,146,186,0.3); }
  50%       { box-shadow: 0 4px 32px rgba(122,234,244,0.6); }
}
.btn-glow { animation: glowPulse 2.5s ease-in-out infinite; }

/* Respect user motion preferences */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## API Integration Map

| Screen | Endpoint |
|---|---|
| Auth | `POST /v1/auth/telegram` |
| Feed | `GET /v1/feed?cursor=&limit=` |
| My Trips | `GET /v1/trips` (filter by user membership) |
| Trip Detail | `GET /v1/trips/:id` |
| Create Trip | `POST /v1/trips` |
| Members | `GET /v1/trips/:id` → members array |
| Join Trip | `POST /v1/trips/:id/join` |
| Photos | `GET /v1/trips/:id/photos` |
| Upload Photo | `POST /v1/trips/:id/photos/presign` then `POST /v1/trips/:id/photos` |
| Comments | `GET /v1/trips/:id/comments` / `POST /v1/trips/:id/comments` |
| Vote | `POST /v1/trips/:id/votes` |
| Profile | `GET /v1/me` |
| Stats | `GET /v1/me/stats` |
| World Map | `GET /v1/me/world` |

---

## World Map SVG Resource

Use the free public domain SVG from simplemaps:

```
https://simplemaps.com/static/svg/world.svg
```

Each `<path>` element has `id="XX"` (ISO alpha-2). Fetch it, inject inline, then set `fill` on visited paths after mount using `animate()` from Motion One.

---

## Performance Checklist

- [ ] `loading="lazy"` on all `<img>` tags
- [ ] `srcset` for responsive trip cover photos
- [ ] Route-based code splitting (automatic in SvelteKit)
- [ ] Skeleton loaders everywhere — never raw spinners
- [ ] `will-change: transform` only on actively animating elements
- [ ] Throttle scroll handlers with `requestAnimationFrame`
- [ ] `env(safe-area-inset-bottom)` for bottom nav on iOS
- [ ] `prefers-reduced-motion` fallback in global CSS
- [ ] Prefetch trip detail routes on hover/focus via `<link rel="prefetch">`
- [ ] Image lazy loading + skeleton placeholder before reveal

---

*Built for Wandr · Telegram Mini App · SvelteKit + Motion One*