import { h as head, a2 as attr_style, a3 as attr_class, a as attr, a4 as ensure_array_like, e as escape_html, a5 as stringify } from "../../chunks/index.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    const TG_LINK = "https://t.me/tripListikBot";
    const GH_LINK = "https://github.com/critiq17/trip-listik";
    let heroReady = false;
    let mouseX = 50;
    let mouseY = 50;
    const words = ["Plan trips.", "Travel together."];
    const features = [
      {
        label: "01",
        title: "Trip Feed",
        desc: "Browse curated trips from the community. Filter by friends, popular destinations, or what's trending right now."
      },
      {
        label: "02",
        title: "Trip Planning",
        desc: "Create trips with destinations, dates, and details. Everything organised in one focused, distraction-free interface."
      },
      {
        label: "03",
        title: "Friend Invites",
        desc: "Send invites directly in Telegram. Friends accept or decline — you're notified instantly, right in chat."
      },
      {
        label: "04",
        title: "Smart Notifications",
        desc: "Invite updates land in your Telegram. Accept, decline, or cancel — everything syncs in real time."
      },
      {
        label: "05",
        title: "Open Source",
        desc: "Built in the open. Go backend, SvelteKit frontend. Read the code, fork it, make it yours."
      }
    ];
    head("1uha8ag", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>TripListik — Plan trips. Travel together.</title>`);
      });
      $$renderer3.push(`<meta name="description" content="TripListik is a Telegram Mini App for collaborative trip planning. Browse trips, invite friends, travel together." class="svelte-1uha8ag"/>`);
    });
    $$renderer2.push(`<div class="cursor-glow svelte-1uha8ag"${attr_style(`--mx: ${stringify(mouseX)}%; --my: ${stringify(mouseY)}%`)} aria-hidden="true"></div> <nav${attr_class("svelte-1uha8ag", void 0, { "nav-ready": heroReady })}><div class="wrap nav-inner svelte-1uha8ag"><span class="nav-logo svelte-1uha8ag">TripListik</span> <div class="nav-links svelte-1uha8ag"><a${attr("href", GH_LINK)} target="_blank" rel="noopener noreferrer" class="nav-link svelte-1uha8ag">GitHub</a> <a${attr("href", TG_LINK)} target="_blank" rel="noopener noreferrer" class="btn-ghost btn-sm svelte-1uha8ag">Open in Telegram</a></div></div></nav> <section class="hero svelte-1uha8ag"><div class="wrap hero-inner svelte-1uha8ag"><div class="hero-left svelte-1uha8ag"><p${attr_class("kicker svelte-1uha8ag", void 0, { "kicker-in": heroReady })}>Telegram Mini App</p> <h1 class="serif hero-h1 svelte-1uha8ag" aria-label="Plan trips. Travel together."><!--[-->`);
    const each_array = ensure_array_like(words);
    for (let wi = 0, $$length = each_array.length; wi < $$length; wi++) {
      let word = each_array[wi];
      $$renderer2.push(`<span${attr_class("word svelte-1uha8ag", void 0, { "word-in": heroReady })}${attr_style(`animation-delay: ${stringify(wi * 180 + 80)}ms`)}>${escape_html(word)}</span>`);
      if (wi < words.length - 1) {
        $$renderer2.push("<!--[0-->");
        $$renderer2.push(`<br class="svelte-1uha8ag"/>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]-->`);
    }
    $$renderer2.push(`<!--]--></h1> <p${attr_class("lead svelte-1uha8ag", void 0, { "lead-in": heroReady })}>TripListik lives inside Telegram — no install required.
				Browse community trips, create your own, and invite friends right from the chat.</p> <div${attr_class("hero-cta svelte-1uha8ag", void 0, { "cta-in": heroReady })}><a${attr("href", TG_LINK)} target="_blank" rel="noopener noreferrer" class="btn-solid svelte-1uha8ag"><svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" class="svelte-1uha8ag"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.447 1.394c-.16.16-.295.295-.605.295l.213-3.053 5.56-5.023c.242-.213-.054-.333-.373-.12L8.32 13.617l-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.537-.194 1.006.131.833.94z" class="svelte-1uha8ag"></path></svg> Open in Telegram</a> <a href="#features" class="btn-text svelte-1uha8ag">See features ↓</a></div></div> <div${attr_class("hero-phones svelte-1uha8ag", void 0, { "phones-in": heroReady })} aria-hidden="true"><div class="phone phone-b svelte-1uha8ag"><div class="pscreen svelte-1uha8ag"><div class="pbar svelte-1uha8ag"></div> <div class="pheader svelte-1uha8ag"><span class="plogo svelte-1uha8ag">TripListik</span></div> <div class="pcard svelte-1uha8ag"><div class="pimg pimg-a svelte-1uha8ag"></div> <div class="plines svelte-1uha8ag"><div class="pl w75 svelte-1uha8ag"></div> <div class="pl w50 dim svelte-1uha8ag"></div></div></div> <div class="pcard svelte-1uha8ag"><div class="pimg pimg-b svelte-1uha8ag"></div> <div class="plines svelte-1uha8ag"><div class="pl w65 svelte-1uha8ag"></div> <div class="pl w45 dim svelte-1uha8ag"></div></div></div></div></div> <div class="phone phone-f svelte-1uha8ag"><div class="pscreen svelte-1uha8ag"><div class="pbar svelte-1uha8ag"></div> <div class="ptrip-img svelte-1uha8ag"></div> <div class="ptrip-body svelte-1uha8ag"><div class="pl w70 svelte-1uha8ag" style="height:10px;border-radius:5px"></div> <div class="ptag svelte-1uha8ag">📍 Reykjavik · Jan 2026</div> <div class="pl w90 dim svelte-1uha8ag" style="margin-top:12px"></div> <div class="pl w80 dim svelte-1uha8ag"></div> <div class="pl w60 dim svelte-1uha8ag"></div> <div class="pbtn svelte-1uha8ag">Join Trip</div></div></div></div></div></div> <div${attr_class("scroll-hint svelte-1uha8ag", void 0, { "hint-in": heroReady })} aria-hidden="true"><div class="scroll-line svelte-1uha8ag"></div></div></section> <div class="rule svelte-1uha8ag"></div> <section class="section svelte-1uha8ag" id="features"><div class="wrap svelte-1uha8ag"><div class="section-head reveal svelte-1uha8ag"><p class="kicker svelte-1uha8ag">Features</p> <h2 class="serif svelte-1uha8ag">Everything you need<br class="svelte-1uha8ag"/>to plan together</h2></div> <div class="features-list svelte-1uha8ag"><!--[-->`);
    const each_array_1 = ensure_array_like(features);
    for (let i = 0, $$length = each_array_1.length; i < $$length; i++) {
      let feat = each_array_1[i];
      $$renderer2.push(`<div class="feat-row reveal-left svelte-1uha8ag"${attr_style(`transition-delay: ${stringify(i * 70)}ms`)}><span class="feat-num svelte-1uha8ag">${escape_html(feat.label)}</span> <div class="feat-body svelte-1uha8ag"><h3 class="svelte-1uha8ag">${escape_html(feat.title)}</h3> <p class="svelte-1uha8ag">${escape_html(feat.desc)}</p></div> <div class="feat-arrow svelte-1uha8ag">→</div></div>`);
    }
    $$renderer2.push(`<!--]--></div></div></section> <div class="rule svelte-1uha8ag"></div> <section class="section section-dark svelte-1uha8ag"><div class="wrap svelte-1uha8ag"><div class="section-head reveal svelte-1uha8ag"><p class="kicker svelte-1uha8ag">In Action</p> <h2 class="serif svelte-1uha8ag">Built for Telegram,<br class="svelte-1uha8ag"/>feels like home</h2></div></div> <div class="mockups-scroll svelte-1uha8ag"><div class="mockups-track svelte-1uha8ag"><div class="mwrap reveal-up svelte-1uha8ag" style="transition-delay:0ms"><div class="mphone fl-a svelte-1uha8ag"><div class="pscreen svelte-1uha8ag"><div class="pbar svelte-1uha8ag"></div> <div class="pheader svelte-1uha8ag"><span class="plogo svelte-1uha8ag">Feed</span></div> <div class="pcard svelte-1uha8ag"><div class="pimg pimg-a svelte-1uha8ag"></div> <div class="plines svelte-1uha8ag"><div class="pl w75 svelte-1uha8ag"></div> <div class="pl w50 dim svelte-1uha8ag"></div></div></div> <div class="pcard svelte-1uha8ag"><div class="pimg pimg-b svelte-1uha8ag"></div> <div class="plines svelte-1uha8ag"><div class="pl w65 svelte-1uha8ag"></div> <div class="pl w45 dim svelte-1uha8ag"></div></div></div> <div class="chip-row svelte-1uha8ag"><div class="mchip active svelte-1uha8ag">All</div> <div class="mchip svelte-1uha8ag">Friends</div> <div class="mchip svelte-1uha8ag">Popular</div></div></div></div> <p class="mcap svelte-1uha8ag">Browse the feed</p></div> <div class="mwrap reveal-up svelte-1uha8ag" style="transition-delay:100ms"><div class="mphone fl-b svelte-1uha8ag"><div class="pscreen svelte-1uha8ag"><div class="pbar svelte-1uha8ag"></div> <div class="ptrip-img pimg-c svelte-1uha8ag"></div> <div class="ptrip-body svelte-1uha8ag"><div class="pl w70 svelte-1uha8ag" style="height:10px;border-radius:5px"></div> <div class="ptag svelte-1uha8ag">📍 Iceland · 5 days</div> <div class="pl w90 dim svelte-1uha8ag" style="margin-top:12px"></div> <div class="pl w80 dim svelte-1uha8ag"></div> <div class="pl w60 dim svelte-1uha8ag"></div> <div class="pbtn svelte-1uha8ag">Join Trip</div></div></div></div> <p class="mcap svelte-1uha8ag">Trip details</p></div> <div class="mwrap reveal-up svelte-1uha8ag" style="transition-delay:200ms"><div class="mphone fl-c svelte-1uha8ag"><div class="pscreen svelte-1uha8ag"><div class="pbar svelte-1uha8ag"></div> <div class="pheader svelte-1uha8ag"><span class="plogo svelte-1uha8ag">Inbox</span></div> <div class="inv-card svelte-1uha8ag"><div class="inv-av av1 svelte-1uha8ag"></div> <div class="inv-body svelte-1uha8ag"><div class="pl w70 svelte-1uha8ag"></div> <div class="pl w55 dim svelte-1uha8ag" style="margin-top:5px"></div> <div class="inv-btns svelte-1uha8ag"><div class="ibtn ibtn-yes svelte-1uha8ag">Accept</div> <div class="ibtn ibtn-no svelte-1uha8ag">Decline</div></div></div></div> <div class="inv-card svelte-1uha8ag" style="opacity:.35;margin-top:4px"><div class="inv-av av2 svelte-1uha8ag"></div> <div class="inv-body svelte-1uha8ag"><div class="pl w60 svelte-1uha8ag"></div> <div class="pl w40 dim svelte-1uha8ag" style="margin-top:5px"></div></div></div></div></div> <p class="mcap svelte-1uha8ag">Friend invites</p></div></div></div></section> <div class="rule svelte-1uha8ag"></div> <section class="section svelte-1uha8ag"><div class="wrap svelte-1uha8ag"><div class="gh-card reveal svelte-1uha8ag"><div class="gh-left svelte-1uha8ag"><div class="gh-icon-wrap svelte-1uha8ag"><svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" class="svelte-1uha8ag"><path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z" class="svelte-1uha8ag"></path></svg></div> <div class="svelte-1uha8ag"><p class="gh-label svelte-1uha8ag">Open Source</p> <p class="gh-repo svelte-1uha8ag">critiq17 / trip-listik</p> <p class="gh-desc svelte-1uha8ag">Go backend · SvelteKit Mini App · PostgreSQL</p></div></div> <a${attr("href", GH_LINK)} target="_blank" rel="noopener noreferrer" class="btn-ghost svelte-1uha8ag">View on GitHub →</a></div></div></section> <div class="rule svelte-1uha8ag"></div> <section class="section cta-sec svelte-1uha8ag"><div class="wrap svelte-1uha8ag"><div class="cta-inner reveal svelte-1uha8ag"><p class="kicker svelte-1uha8ag">Get Started</p> <h2 class="serif cta-h2 svelte-1uha8ag">Ready to plan<br class="svelte-1uha8ag"/>your next trip?</h2> <p class="lead cta-lead svelte-1uha8ag">Open TripListik in Telegram. Browse the feed,
				create a trip, and invite your crew — it takes 30 seconds.</p> <div class="cta-btn-wrap svelte-1uha8ag"><a${attr("href", TG_LINK)} target="_blank" rel="noopener noreferrer" class="btn-solid btn-large svelte-1uha8ag"><svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true" class="svelte-1uha8ag"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.447 1.394c-.16.16-.295.295-.605.295l.213-3.053 5.56-5.023c.242-.213-.054-.333-.373-.12L8.32 13.617l-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.537-.194 1.006.131.833.94z" class="svelte-1uha8ag"></path></svg> Open TripListik</a> <div class="btn-pulse-ring svelte-1uha8ag" aria-hidden="true"></div></div> <p class="sub-note svelte-1uha8ag">Free · No install · Lives inside Telegram</p></div></div></section> <footer class="svelte-1uha8ag"><div class="wrap footer-inner svelte-1uha8ag"><span class="footer-logo svelte-1uha8ag">TripListik</span> <div class="footer-links svelte-1uha8ag"><a${attr("href", GH_LINK)} target="_blank" rel="noopener noreferrer" class="svelte-1uha8ag">GitHub</a> <a${attr("href", TG_LINK)} target="_blank" rel="noopener noreferrer" class="svelte-1uha8ag">Telegram</a></div> <p class="footer-copy svelte-1uha8ag">© 2025 TripListik</p></div></footer>`);
  });
}
export {
  _page as default
};
