import { clearTokens, getToken, setRefreshToken, setToken } from '$lib/api';
import { getTelegramInitData } from '$lib/telegram';
import { env } from '$env/dynamic/public';

const baseUrl = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

/**
 * Authenticate using Telegram initData.
 * Always re-authenticates on every app mount (initData is always available in TG context).
 */
export const authenticate = async (): Promise<{
	token: string;
	refresh_token: string;
	user: Record<string, unknown>;
}> => {
	const initData = getTelegramInitData();
	if (!initData) {
		throw new Error('Telegram init data is missing — are you inside Telegram?');
	}
	// start_param carries the referral / invite deep link. The backend prefers
	// the signed copy inside initData; this is a fallback for older clients.
	const startParam =
		(window as unknown as { Telegram?: { WebApp?: { initDataUnsafe?: { start_param?: string } } } })
			.Telegram?.WebApp?.initDataUnsafe?.start_param ?? '';
	const res = await fetch(`${baseUrl}/v1/auth/telegram`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ initData, start_param: startParam })
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({}));
		throw new Error(err?.error?.message ?? 'Auth failed');
	}
	const data = await res.json();
	setToken(data.token);
	if (data.refresh_token) setRefreshToken(data.refresh_token);
	return data;
};

/**
 * ensureAuth: always re-validate initData on every mount.
 * Tokens are in-memory only — every time the app mounts (Telegram always provides initData),
 * we re-authenticate. This is fast (~50ms) and more secure than trusting persisted tokens.
 */
export const ensureAuth = async (_path: string): Promise<void> => {
	// If we already have a token in memory (e.g. navigating between routes in same session),
	// skip re-auth to avoid unnecessary API calls.
	if (getToken()) return;

	const initData = getTelegramInitData();
	if (!initData) {
		// Not in Telegram context (dev mode) — skip auth without redirect
		console.warn('[auth] No Telegram initData — running in dev mode without auth');
		return;
	}

	try {
		await authenticate();
	} catch (err) {
		console.error('[auth] Authentication failed', err);
		clearTokens();
	}
};
