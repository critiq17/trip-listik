import { goto } from '$app/navigation';
import { clearTokens, getRefreshToken, getToken, setRefreshToken, setToken } from '$lib/api';
import { getTelegramInitData } from '$lib/telegram';
import { env } from '$env/dynamic/public';

const baseUrl = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

export const ensureAuth = async (path: string) => {
	const token = getToken();
	if (!token && path !== '/auth') {
		const refreshToken = getRefreshToken();
		if (refreshToken) {
			try {
				const res = await fetch(`${baseUrl}/v1/auth/refresh`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ refresh_token: refreshToken })
				});
				if (!res.ok) {
					clearTokens();
					await goto('/auth');
					return;
				}
				const data = await res.json();
				if (data?.token) setToken(data.token);
				if (data?.refresh_token) setRefreshToken(data.refresh_token);
				return;
			} catch {
				clearTokens();
				await goto('/auth');
				return;
			}
		}
		await goto('/auth');
	}
	if (token && path === '/auth') {
		await goto('/');
	}
};

export const authenticate = async () => {
	const initData = getTelegramInitData();
	if (!initData) {
		throw new Error('Telegram init data is missing');
	}
	const res = await fetch(`${baseUrl}/v1/auth/telegram`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ initData })
	});
	if (!res.ok) {
		throw new Error('Auth failed');
	}
	const data = await res.json();
	setToken(data.token);
	if (data.refresh_token) {
		setRefreshToken(data.refresh_token);
	}
	return data;
};
