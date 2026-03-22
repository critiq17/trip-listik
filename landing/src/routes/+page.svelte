<script lang="ts">
	import { onMount } from 'svelte';
	import { inview } from '$lib/actions/inview';

	const TG_LINK = 'https://t.me/tripListikBot';
	const GH_LINK = 'https://github.com/critiq17/trip-listik';

	let heroReady = $state(false);
	let mouseX = $state(50);
	let mouseY = $state(50);

	const words = ['Plan trips.', 'Travel together.'];

	onMount(() => {
		setTimeout(() => { heroReady = true; }, 120);

		const onMove = (e: MouseEvent) => {
			mouseX = (e.clientX / window.innerWidth) * 100;
			mouseY = (e.clientY / window.innerHeight) * 100;
		};
		window.addEventListener('mousemove', onMove, { passive: true });
		return () => window.removeEventListener('mousemove', onMove);
	});

	function onEnter(e: Event) {
		(e.currentTarget as HTMLElement).classList.add('visible');
	}

	const features = [
		{
			label: '01',
			title: 'Trip Feed',
			desc: 'Browse curated trips from the community. Filter by friends, popular destinations, or what\'s trending right now.'
		},
		{
			label: '02',
			title: 'Trip Planning',
			desc: 'Create trips with destinations, dates, and details. Everything organised in one focused, distraction-free interface.'
		},
		{
			label: '03',
			title: 'Friend Invites',
			desc: 'Send invites directly in Telegram. Friends accept or decline — you\'re notified instantly, right in chat.'
		},
		{
			label: '04',
			title: 'Smart Notifications',
			desc: 'Invite updates land in your Telegram. Accept, decline, or cancel — everything syncs in real time.'
		},
		{
			label: '05',
			title: 'Open Source',
			desc: 'Built in the open. Go backend, SvelteKit frontend. Read the code, fork it, make it yours.'
		}
	];
</script>

<svelte:head>
	<title>TripListik — Plan trips. Travel together.</title>
	<meta name="description" content="TripListik is a Telegram Mini App for collaborative trip planning. Browse trips, invite friends, travel together." />
</svelte:head>

<!-- ── Cursor glow (desktop) ─────────────────────────────── -->
<div
	class="cursor-glow"
	style="--mx: {mouseX}%; --my: {mouseY}%"
	aria-hidden="true"
></div>

<!-- ── Nav ──────────────────────────────────────────────────── -->
<nav class:nav-ready={heroReady}>
	<div class="wrap nav-inner">
		<span class="nav-logo">TripListik</span>
		<div class="nav-links">
			<a href={GH_LINK} target="_blank" rel="noopener noreferrer" class="nav-link">GitHub</a>
			<a href={TG_LINK} target="_blank" rel="noopener noreferrer" class="btn-ghost btn-sm">
				Open in Telegram
			</a>
		</div>
	</div>
</nav>

<!-- ── Hero ─────────────────────────────────────────────────── -->
<section class="hero">
	<div class="wrap hero-inner">
		<div class="hero-left">
			<p class="kicker" class:kicker-in={heroReady}>Telegram Mini App</p>

			<h1 class="serif hero-h1" aria-label="Plan trips. Travel together.">
				{#each words as word, wi}
					<span
						class="word"
						class:word-in={heroReady}
						style="animation-delay: {wi * 180 + 80}ms"
					>{word}</span>{#if wi < words.length - 1}<br />{/if}
				{/each}
			</h1>

			<p class="lead" class:lead-in={heroReady}>
				TripListik lives inside Telegram — no install required.
				Browse community trips, create your own, and invite friends right from the chat.
			</p>

			<div class="hero-cta" class:cta-in={heroReady}>
				<a href={TG_LINK} target="_blank" rel="noopener noreferrer" class="btn-solid">
					<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
						<path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.447 1.394c-.16.16-.295.295-.605.295l.213-3.053 5.56-5.023c.242-.213-.054-.333-.373-.12L8.32 13.617l-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.537-.194 1.006.131.833.94z"/>
					</svg>
					Open in Telegram
				</a>
				<a href="#features" class="btn-text">See features ↓</a>
			</div>
		</div>

		<div class="hero-phones" class:phones-in={heroReady} aria-hidden="true">
			<div class="phone phone-b">
				<div class="pscreen">
					<div class="pbar"></div>
					<div class="pheader"><span class="plogo">TripListik</span></div>
					<div class="pcard">
						<div class="pimg pimg-a"></div>
						<div class="plines">
							<div class="pl w75"></div>
							<div class="pl w50 dim"></div>
						</div>
					</div>
					<div class="pcard">
						<div class="pimg pimg-b"></div>
						<div class="plines">
							<div class="pl w65"></div>
							<div class="pl w45 dim"></div>
						</div>
					</div>
				</div>
			</div>
			<div class="phone phone-f">
				<div class="pscreen">
					<div class="pbar"></div>
					<div class="ptrip-img"></div>
					<div class="ptrip-body">
						<div class="pl w70" style="height:10px;border-radius:5px"></div>
						<div class="ptag">📍 Reykjavik · Jan 2026</div>
						<div class="pl w90 dim" style="margin-top:12px"></div>
						<div class="pl w80 dim"></div>
						<div class="pl w60 dim"></div>
						<div class="pbtn">Join Trip</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- scroll hint -->
	<div class="scroll-hint" class:hint-in={heroReady} aria-hidden="true">
		<div class="scroll-line"></div>
	</div>
</section>

<div class="rule"></div>

<!-- ── Features ─────────────────────────────────────────────── -->
<section class="section" id="features">
	<div class="wrap">
		<div class="section-head reveal" use:inview onenter={onEnter}>
			<p class="kicker">Features</p>
			<h2 class="serif">Everything you need<br />to plan together</h2>
		</div>

		<div class="features-list">
			{#each features as feat, i}
				<div
					class="feat-row reveal-left"
					use:inview
					onenter={onEnter}
					style="transition-delay: {i * 70}ms"
				>
					<span class="feat-num">{feat.label}</span>
					<div class="feat-body">
						<h3>{feat.title}</h3>
						<p>{feat.desc}</p>
					</div>
					<div class="feat-arrow">→</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<div class="rule"></div>

<!-- ── Mockups ───────────────────────────────────────────────── -->
<section class="section section-dark">
	<div class="wrap">
		<div class="section-head reveal" use:inview onenter={onEnter}>
			<p class="kicker">In Action</p>
			<h2 class="serif">Built for Telegram,<br />feels like home</h2>
		</div>
	</div>

	<div class="mockups-scroll">
		<div class="mockups-track">
			<div class="mwrap reveal-up" use:inview onenter={onEnter} style="transition-delay:0ms">
				<div class="mphone fl-a">
					<div class="pscreen">
						<div class="pbar"></div>
						<div class="pheader"><span class="plogo">Feed</span></div>
						<div class="pcard">
							<div class="pimg pimg-a"></div>
							<div class="plines">
								<div class="pl w75"></div>
								<div class="pl w50 dim"></div>
							</div>
						</div>
						<div class="pcard">
							<div class="pimg pimg-b"></div>
							<div class="plines">
								<div class="pl w65"></div>
								<div class="pl w45 dim"></div>
							</div>
						</div>
						<div class="chip-row">
							<div class="mchip active">All</div>
							<div class="mchip">Friends</div>
							<div class="mchip">Popular</div>
						</div>
					</div>
				</div>
				<p class="mcap">Browse the feed</p>
			</div>

			<div class="mwrap reveal-up" use:inview onenter={onEnter} style="transition-delay:100ms">
				<div class="mphone fl-b">
					<div class="pscreen">
						<div class="pbar"></div>
						<div class="ptrip-img pimg-c"></div>
						<div class="ptrip-body">
							<div class="pl w70" style="height:10px;border-radius:5px"></div>
							<div class="ptag">📍 Iceland · 5 days</div>
							<div class="pl w90 dim" style="margin-top:12px"></div>
							<div class="pl w80 dim"></div>
							<div class="pl w60 dim"></div>
							<div class="pbtn">Join Trip</div>
						</div>
					</div>
				</div>
				<p class="mcap">Trip details</p>
			</div>

			<div class="mwrap reveal-up" use:inview onenter={onEnter} style="transition-delay:200ms">
				<div class="mphone fl-c">
					<div class="pscreen">
						<div class="pbar"></div>
						<div class="pheader"><span class="plogo">Inbox</span></div>
						<div class="inv-card">
							<div class="inv-av av1"></div>
							<div class="inv-body">
								<div class="pl w70"></div>
								<div class="pl w55 dim" style="margin-top:5px"></div>
								<div class="inv-btns">
									<div class="ibtn ibtn-yes">Accept</div>
									<div class="ibtn ibtn-no">Decline</div>
								</div>
							</div>
						</div>
						<div class="inv-card" style="opacity:.35;margin-top:4px">
							<div class="inv-av av2"></div>
							<div class="inv-body">
								<div class="pl w60"></div>
								<div class="pl w40 dim" style="margin-top:5px"></div>
							</div>
						</div>
					</div>
				</div>
				<p class="mcap">Friend invites</p>
			</div>
		</div>
	</div>
</section>

<div class="rule"></div>

<!-- ── GitHub ────────────────────────────────────────────────── -->
<section class="section">
	<div class="wrap">
		<div class="gh-card reveal" use:inview onenter={onEnter}>
			<div class="gh-left">
				<div class="gh-icon-wrap">
					<svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
						<path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
					</svg>
				</div>
				<div>
					<p class="gh-label">Open Source</p>
					<p class="gh-repo">critiq17 / trip-listik</p>
					<p class="gh-desc">Go backend · SvelteKit Mini App · PostgreSQL</p>
				</div>
			</div>
			<a href={GH_LINK} target="_blank" rel="noopener noreferrer" class="btn-ghost">
				View on GitHub →
			</a>
		</div>
	</div>
</section>

<div class="rule"></div>

<!-- ── CTA ───────────────────────────────────────────────────── -->
<section class="section cta-sec">
	<div class="wrap">
		<div class="cta-inner reveal" use:inview onenter={onEnter}>
			<p class="kicker">Get Started</p>
			<h2 class="serif cta-h2">Ready to plan<br />your next trip?</h2>
			<p class="lead cta-lead">
				Open TripListik in Telegram. Browse the feed,
				create a trip, and invite your crew — it takes 30 seconds.
			</p>
			<div class="cta-btn-wrap">
				<a href={TG_LINK} target="_blank" rel="noopener noreferrer" class="btn-solid btn-large">
					<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
						<path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.447 1.394c-.16.16-.295.295-.605.295l.213-3.053 5.56-5.023c.242-.213-.054-.333-.373-.12L8.32 13.617l-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.537-.194 1.006.131.833.94z"/>
					</svg>
					Open TripListik
				</a>
				<div class="btn-pulse-ring" aria-hidden="true"></div>
			</div>
			<p class="sub-note">Free · No install · Lives inside Telegram</p>
		</div>
	</div>
</section>

<!-- ── Footer ────────────────────────────────────────────────── -->
<footer>
	<div class="wrap footer-inner">
		<span class="footer-logo">TripListik</span>
		<div class="footer-links">
			<a href={GH_LINK} target="_blank" rel="noopener noreferrer">GitHub</a>
			<a href={TG_LINK} target="_blank" rel="noopener noreferrer">Telegram</a>
		</div>
		<p class="footer-copy">© 2025 TripListik</p>
	</div>
</footer>

<style>
/* ── Tokens ───────────────────────────────────────────────── */
:root {
	--bg:      #0c0c0c;
	--surface: #111;
	--lift:    #1a1a1a;
	--border:  rgba(255,255,255,0.07);
	--text:    #f0f0f0;
	--sub:     #888;
	--muted:   #3a3a3a;
	--ease:    cubic-bezier(0.25, 0.46, 0.45, 0.94);
	--spring:  cubic-bezier(0.34, 1.56, 0.64, 1);
	--out:     cubic-bezier(0.16, 1, 0.3, 1);
}

:global(body) {
	background: var(--bg);
	color: var(--text);
	font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
	-webkit-font-smoothing: antialiased;
	overflow-x: hidden;
}

.wrap {
	max-width: 1100px;
	margin: 0 auto;
	padding: 0 32px;
}

.serif {
	font-family: 'DM Serif Display', Georgia, serif;
	font-weight: 400;
	letter-spacing: -0.01em;
	line-height: 1.08;
	color: var(--text);
}

/* ── Cursor glow ──────────────────────────────────────────── */
.cursor-glow {
	pointer-events: none;
	position: fixed;
	inset: 0;
	z-index: 0;
	background: radial-gradient(
		600px circle at var(--mx) var(--my),
		rgba(255,255,255,0.032),
		transparent 60%
	);
	transition: background 0.1s;
}

@media (hover: none) {
	.cursor-glow { display: none; }
}

/* ── Rule ─────────────────────────────────────────────────── */
.rule {
	width: 100%;
	height: 1px;
	background: var(--border);
}

/* ── Kicker ───────────────────────────────────────────────── */
.kicker {
	font-size: 11px;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.18em;
	color: var(--muted);
	margin-bottom: 20px;
}

/* ── Buttons ──────────────────────────────────────────────── */
.btn-solid {
	display: inline-flex;
	align-items: center;
	gap: 9px;
	background: var(--text);
	color: var(--bg);
	font-weight: 600;
	font-size: 15px;
	padding: 14px 28px;
	border-radius: 10px;
	border: none;
	cursor: pointer;
	text-decoration: none;
	transition: opacity 0.18s var(--ease), transform 0.22s var(--spring);
	position: relative;
	z-index: 1;
}

.btn-solid:hover { opacity: 0.87; transform: translateY(-2px) scale(1.01); }
.btn-solid:active { transform: translateY(0) scale(0.99); }
.btn-solid.btn-large { font-size: 17px; padding: 16px 36px; }

.btn-ghost {
	display: inline-flex;
	align-items: center;
	gap: 6px;
	background: transparent;
	color: var(--sub);
	font-weight: 500;
	font-size: 14px;
	padding: 11px 22px;
	border-radius: 10px;
	border: 1px solid var(--border);
	cursor: pointer;
	text-decoration: none;
	transition: color 0.18s, border-color 0.18s, transform 0.22s var(--spring);
}

.btn-ghost:hover {
	color: var(--text);
	border-color: rgba(255,255,255,0.22);
	transform: translateY(-2px);
}

.btn-ghost.btn-sm { font-size: 13px; padding: 9px 18px; }

.btn-text {
	font-size: 14px;
	font-weight: 500;
	color: var(--muted);
	text-decoration: none;
	transition: color 0.18s;
	letter-spacing: 0.01em;
}

.btn-text:hover { color: var(--sub); }

/* CTA pulse ring */
.cta-btn-wrap {
	position: relative;
	display: inline-flex;
	align-items: center;
	justify-content: center;
}

.btn-pulse-ring {
	position: absolute;
	inset: -8px;
	border-radius: 18px;
	border: 1px solid rgba(255,255,255,0.15);
	animation: pulseRing 2.4s var(--ease) infinite;
	pointer-events: none;
}

@keyframes pulseRing {
	0%   { opacity: 0; transform: scale(0.96); }
	40%  { opacity: 1; }
	100% { opacity: 0; transform: scale(1.08); }
}

/* ── Nav ──────────────────────────────────────────────────── */
nav {
	position: fixed;
	inset: 0 0 auto 0;
	z-index: 100;
	background: rgba(12,12,12,0.75);
	backdrop-filter: blur(16px);
	-webkit-backdrop-filter: blur(16px);
	border-bottom: 1px solid var(--border);
	opacity: 0;
	transform: translateY(-8px);
	transition: opacity 0.5s var(--ease), transform 0.5s var(--ease);
}

nav.nav-ready {
	opacity: 1;
	transform: translateY(0);
}

.nav-inner {
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 62px;
}

.nav-logo {
	font-size: 17px;
	font-weight: 700;
	letter-spacing: -0.025em;
}

.nav-links {
	display: flex;
	align-items: center;
	gap: 20px;
}

.nav-link {
	font-size: 14px;
	font-weight: 500;
	color: var(--sub);
	text-decoration: none;
	transition: color 0.18s;
}
.nav-link:hover { color: var(--text); }

/* ── Hero ─────────────────────────────────────────────────── */
.hero {
	position: relative;
	min-height: 100dvh;
	display: flex;
	flex-direction: column;
	justify-content: center;
	padding-top: 62px;
	overflow: hidden;
}

.hero-inner {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 80px;
	padding-top: 80px;
	padding-bottom: 80px;
}

.hero-left { flex: 1; max-width: 560px; }

/* kicker */
.kicker.kicker-in { animation: fadeUp 0.6s var(--out) forwards; }
@keyframes fadeUp {
	from { opacity: 0; transform: translateY(14px); }
	to   { opacity: 1; transform: translateY(0); }
}

/* word-by-word blur reveal */
.hero-h1 {
	font-size: clamp(52px, 6.5vw, 88px);
	margin-bottom: 28px;
	display: block;
}

.word {
	display: inline-block;
	opacity: 0;
	filter: blur(10px);
	transform: translateY(16px);
}

.word.word-in {
	animation: wordReveal 0.7s var(--out) forwards;
}

@keyframes wordReveal {
	from { opacity: 0; filter: blur(10px); transform: translateY(16px); }
	to   { opacity: 1; filter: blur(0);    transform: translateY(0); }
}

/* lead */
.lead {
	font-size: 17px;
	line-height: 1.75;
	color: var(--sub);
	margin-bottom: 40px;
	opacity: 0;
}

.lead.lead-in {
	animation: fadeUp 0.7s 0.42s var(--out) forwards;
}

/* cta */
.hero-cta {
	display: flex;
	align-items: center;
	gap: 20px;
	flex-wrap: wrap;
	opacity: 0;
}

.hero-cta.cta-in {
	animation: fadeUp 0.7s 0.56s var(--out) forwards;
}

/* scroll hint */
.scroll-hint {
	position: absolute;
	bottom: 36px;
	left: 50%;
	transform: translateX(-50%);
	opacity: 0;
	transition: opacity 0.6s 1.1s var(--ease);
}

.scroll-hint.hint-in { opacity: 1; }

.scroll-line {
	width: 1px;
	height: 48px;
	background: linear-gradient(to bottom, rgba(255,255,255,0.25), transparent);
	margin: 0 auto;
	animation: scrollDrop 1.8s var(--ease) infinite;
}

@keyframes scrollDrop {
	0%   { transform: scaleY(0); transform-origin: top; opacity: 1; }
	50%  { transform: scaleY(1); transform-origin: top; opacity: 1; }
	51%  { transform: scaleY(1); transform-origin: bottom; }
	100% { transform: scaleY(0); transform-origin: bottom; opacity: 0; }
}

/* ── Hero phones ──────────────────────────────────────────── */
.hero-phones {
	flex-shrink: 0;
	position: relative;
	width: 300px;
	height: 430px;
}

.phone {
	position: absolute;
	background: #0f0f0f;
	border-radius: 28px;
	border: 1px solid rgba(255,255,255,0.09);
	box-shadow: 0 48px 96px rgba(0,0,0,0.85), 0 0 0 0.5px rgba(255,255,255,0.04);
	overflow: hidden;
}

.phone-b {
	width: 200px;
	top: 16px;
	right: 0;
	opacity: 0;
	transform: translateX(24px) rotate(5deg) scale(0.92);
	transition: opacity 0.9s 0.6s var(--out), transform 0.9s 0.6s var(--out);
	animation: floatA 5.5s ease-in-out 1.6s infinite alternate;
}

.phone-f {
	width: 220px;
	bottom: 0;
	left: 0;
	opacity: 0;
	transform: translateX(-24px) rotate(-4deg) scale(0.92);
	transition: opacity 0.9s 0.8s var(--out), transform 0.9s 0.8s var(--out);
	animation: floatB 7s ease-in-out 1.8s infinite alternate;
}

.hero-phones.phones-in .phone-b {
	opacity: 1;
	transform: rotate(5deg) scale(1);
}

.hero-phones.phones-in .phone-f {
	opacity: 1;
	transform: rotate(-4deg) scale(1);
}

/* ── Phone internals ──────────────────────────────────────── */
.pscreen {
	padding: 12px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.pbar {
	width: 48px;
	height: 3px;
	background: rgba(255,255,255,0.1);
	border-radius: 2px;
	margin: 0 auto 2px;
}

.pheader { padding: 2px 0 4px; }

.plogo {
	font-size: 12px;
	font-weight: 700;
	color: var(--text);
}

.pcard {
	background: #1c1c1c;
	border-radius: 10px;
	overflow: hidden;
}

.pimg { height: 64px; }
.pimg-a { background: linear-gradient(135deg, #1e2e22, #283822); }
.pimg-b { background: linear-gradient(135deg, #1e1e2e, #28283a); }
.pimg-c { background: linear-gradient(135deg, #2a2020, #3a2828); }

.plines {
	padding: 8px 10px;
	display: flex;
	flex-direction: column;
	gap: 5px;
}

.pl {
	height: 7px;
	background: rgba(255,255,255,0.13);
	border-radius: 4px;
}

.pl.dim { background: rgba(255,255,255,0.05); }
.w40 { width: 40%; } .w45 { width: 45%; } .w50 { width: 50%; }
.w55 { width: 55%; } .w60 { width: 60%; } .w65 { width: 65%; }
.w70 { width: 70%; } .w75 { width: 75%; } .w80 { width: 80%; }
.w90 { width: 90%; }

.ptrip-img {
	height: 100px;
	background: linear-gradient(135deg, #1e2e22, #2d4030);
}

.ptrip-body {
	padding: 12px;
	display: flex;
	flex-direction: column;
}

.ptag {
	font-size: 9px;
	font-weight: 600;
	color: rgba(255,255,255,0.35);
	margin-top: 8px;
}

.pbtn {
	background: rgba(255,255,255,0.92);
	color: #0c0c0c;
	font-size: 10px;
	font-weight: 700;
	text-align: center;
	padding: 9px;
	border-radius: 9px;
	margin-top: 14px;
}

/* ── Floats ───────────────────────────────────────────────── */
@keyframes floatA {
	from { transform: translateY(0) rotate(5deg); }
	to   { transform: translateY(-12px) rotate(5deg); }
}
@keyframes floatB {
	from { transform: translateY(-5px) rotate(-4deg); }
	to   { transform: translateY(7px) rotate(-4deg); }
}
@keyframes floatC {
	from { transform: translateY(0); }
	to   { transform: translateY(-10px); }
}

/* ── Section ──────────────────────────────────────────────── */
.section {
	padding: 120px 0;
}

.section-dark {
	background: var(--surface);
}

.section-head {
	margin-bottom: 72px;
}

.section-head h2 {
	font-size: clamp(34px, 4.5vw, 56px);
	margin-top: 10px;
}

/* ── Features ─────────────────────────────────────────────── */
.features-list { display: flex; flex-direction: column; }

.feat-row {
	display: flex;
	align-items: center;
	gap: 40px;
	padding: 36px 0;
	border-bottom: 1px solid var(--border);
	cursor: default;
	transition:
		opacity 0.65s var(--ease),
		transform 0.65s var(--ease);
}

.feat-row:first-of-type { border-top: 1px solid var(--border); }

.feat-row:hover .feat-body h3 { color: #fff; }
.feat-row:hover .feat-arrow { opacity: 1; transform: translateX(0); }
.feat-row:hover .feat-num { color: var(--sub); }

.feat-num {
	font-size: 12px;
	font-weight: 500;
	color: var(--muted);
	font-variant-numeric: tabular-nums;
	min-width: 28px;
	flex-shrink: 0;
	transition: color 0.25s;
}

.feat-body {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.feat-body h3 {
	font-size: 19px;
	font-weight: 600;
	color: var(--text);
	letter-spacing: -0.02em;
	transition: color 0.25s;
}

.feat-body p {
	font-size: 15px;
	line-height: 1.7;
	color: var(--sub);
	max-width: 560px;
}

.feat-arrow {
	font-size: 18px;
	color: var(--muted);
	flex-shrink: 0;
	opacity: 0;
	transform: translateX(-8px);
	transition: opacity 0.25s, transform 0.3s var(--spring), color 0.25s;
}

/* ── Mockups ──────────────────────────────────────────────── */
.mockups-scroll {
	overflow-x: auto;
	-webkit-overflow-scrolling: touch;
	scrollbar-width: none;
}

.mockups-scroll::-webkit-scrollbar { display: none; }

.mockups-track {
	display: flex;
	justify-content: center;
	align-items: flex-end;
	gap: 48px;
	padding: 0 32px 4px;
	min-width: max-content;
	margin: 0 auto;
}

.mwrap {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 24px;
}

.mphone {
	background: #0f0f0f;
	border-radius: 26px;
	border: 1px solid rgba(255,255,255,0.09);
	box-shadow: 0 32px 64px rgba(0,0,0,0.75);
	overflow: hidden;
	width: 230px;
}

.fl-a { animation: floatC 5.5s ease-in-out infinite alternate; }
.fl-b { animation: floatC 7s ease-in-out 0.6s infinite alternate; }
.fl-c { animation: floatC 4.8s ease-in-out 1.2s infinite alternate; }

.mcap {
	font-size: 12px;
	font-weight: 500;
	color: var(--muted);
	text-transform: uppercase;
	letter-spacing: 0.12em;
}

.chip-row {
	display: flex;
	gap: 5px;
	margin-top: 2px;
}

.mchip {
	font-size: 9px;
	font-weight: 600;
	padding: 3px 9px;
	border-radius: 20px;
	background: rgba(255,255,255,0.06);
	color: var(--sub);
}

.mchip.active {
	background: rgba(255,255,255,0.15);
	color: var(--text);
}

.inv-card {
	display: flex;
	align-items: flex-start;
	gap: 8px;
	background: #1c1c1c;
	border-radius: 10px;
	padding: 10px;
}

.inv-av {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	flex-shrink: 0;
}

.av1 { background: linear-gradient(135deg, #3a4a3a, #5a6a5a); }
.av2 { background: linear-gradient(135deg, #3a3a4a, #5a5a6a); }

.inv-body { flex: 1; }

.inv-btns {
	display: flex;
	gap: 6px;
	margin-top: 9px;
}

.ibtn {
	flex: 1;
	font-size: 9px;
	font-weight: 700;
	text-align: center;
	padding: 6px;
	border-radius: 7px;
}

.ibtn-yes {
	background: rgba(255,255,255,0.9);
	color: #0c0c0c;
}

.ibtn-no {
	background: rgba(255,255,255,0.07);
	color: var(--sub);
}

/* ── GitHub ───────────────────────────────────────────────── */
.gh-card {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 32px;
	padding: 44px 48px;
	background: var(--surface);
	border-radius: 20px;
	border: 1px solid var(--border);
	flex-wrap: wrap;
	transition: border-color 0.25s;
}

.gh-card:hover { border-color: rgba(255,255,255,0.14); }

.gh-left {
	display: flex;
	align-items: center;
	gap: 22px;
}

.gh-icon-wrap {
	width: 52px;
	height: 52px;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 14px;
	background: var(--lift);
	color: var(--sub);
	flex-shrink: 0;
}

.gh-label {
	font-size: 11px;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.14em;
	color: var(--muted);
	margin-bottom: 5px;
}

.gh-repo {
	font-size: 18px;
	font-weight: 700;
	color: var(--text);
	letter-spacing: -0.02em;
	margin-bottom: 5px;
}

.gh-desc {
	font-size: 14px;
	color: var(--sub);
}

/* ── CTA ──────────────────────────────────────────────────── */
.cta-sec { text-align: center; padding: 140px 0; }

.cta-h2 {
	font-size: clamp(40px, 5.5vw, 64px);
	margin: 10px 0 28px;
}

.cta-lead {
	max-width: 460px;
	margin: 0 auto 44px;
}

.sub-note {
	margin-top: 24px;
	font-size: 12px;
	color: var(--muted);
	letter-spacing: 0.06em;
}

/* ── Footer ───────────────────────────────────────────────── */
footer {
	border-top: 1px solid var(--border);
	padding: 32px 0;
}

.footer-inner {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 16px;
	flex-wrap: wrap;
}

.footer-logo {
	font-size: 14px;
	font-weight: 700;
	color: var(--muted);
}

.footer-links {
	display: flex;
	gap: 24px;
}

.footer-links a {
	font-size: 13px;
	color: var(--muted);
	text-decoration: none;
	transition: color 0.18s;
}
.footer-links a:hover { color: var(--sub); }

.footer-copy {
	font-size: 12px;
	color: var(--muted);
}

/* ── Scroll reveals ───────────────────────────────────────── */
.reveal {
	opacity: 0;
	transform: translateY(28px);
	transition: opacity 0.7s var(--out), transform 0.7s var(--out);
}

.reveal-left {
	opacity: 0;
	transform: translateX(-20px);
	transition: opacity 0.65s var(--out), transform 0.65s var(--out);
}

.reveal-up {
	opacity: 0;
	transform: translateY(40px) scale(0.97);
	transition: opacity 0.75s var(--out), transform 0.75s var(--out);
}

:global(.reveal.visible),
:global(.reveal-left.visible),
:global(.reveal-up.visible) {
	opacity: 1;
	transform: none;
}

/* ── Responsive ───────────────────────────────────────────── */
@media (max-width: 900px) {
	.hero-inner {
		flex-direction: column;
		align-items: flex-start;
		gap: 64px;
		padding-top: 56px;
		padding-bottom: 72px;
	}

	.hero-left { max-width: 100%; }

	.hero-phones {
		width: 260px;
		height: 370px;
		align-self: center;
	}

	.phone-b { width: 176px; }
	.phone-f { width: 194px; }

	.gh-card { padding: 32px; }
}

@media (max-width: 640px) {
	.wrap { padding: 0 20px; }

	.hero-h1 { font-size: clamp(44px, 12vw, 64px); }

	.hero-phones { align-self: center; }

	.section { padding: 80px 0; }
	.cta-sec  { padding: 100px 0; }

	.section-head { margin-bottom: 48px; }

	.feat-row { gap: 20px; padding: 28px 0; }
	.feat-body h3 { font-size: 17px; }
	.feat-body p  { font-size: 14px; }
	.feat-arrow   { display: none; }

	.mockups-track { gap: 28px; padding: 0 20px 4px; }
	.mphone { width: 196px; }

	.gh-card {
		flex-direction: column;
		align-items: flex-start;
		padding: 24px;
		gap: 20px;
	}

	.nav-link { display: none; }

	.footer-inner {
		flex-direction: column;
		align-items: flex-start;
		gap: 10px;
	}
}

@media (max-width: 380px) {
	.hero-phones { width: 220px; height: 320px; }
	.phone-b { width: 154px; }
	.phone-f { width: 170px; }
}
</style>
