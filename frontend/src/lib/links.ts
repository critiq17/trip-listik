// Single source for bot deep links. startapp payloads are routed in
// +layout.svelte and attributed as referrals on the backend during auth.
export const BOT_URL = 'https://t.me/tripListikBot';

export const buildProfileLink = (userId: string) => `${BOT_URL}?startapp=profile_${userId}`;

export const buildJoinLink = (token: string) => `${BOT_URL}?startapp=join_${token}`;

/**
 * Open Telegram's chat picker with a prefilled message containing the link.
 * Works for recipients who never started the bot — the link itself opens
 * the mini app.
 */
export const shareLinkViaTelegram = (link: string, text: string) => {
	const shareUrl = `https://t.me/share/url?url=${encodeURIComponent(link)}&text=${encodeURIComponent(text)}`;
	const tg = (window as unknown as { Telegram?: { WebApp?: { openTelegramLink?: (url: string) => void } } })
		.Telegram?.WebApp;
	if (tg?.openTelegramLink) {
		tg.openTelegramLink(shareUrl);
	} else {
		window.open(shareUrl, '_blank');
	}
};
