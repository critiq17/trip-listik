import { env } from '$env/dynamic/public';

const baseUrl = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

// ── In-memory token store (NOT localStorage — Telegram Mini App security model) ──
let _token = '';
let _refreshToken = '';

export const getToken = () => _token;
export const setToken = (t: string) => { _token = t; };
export const getRefreshToken = () => _refreshToken;
export const setRefreshToken = (t: string) => { _refreshToken = t; };
export const clearTokens = () => { _token = ''; _refreshToken = ''; };

const refreshSession = async () => {
	if (!_refreshToken) throw new Error('Missing refresh token');

	const res = await fetch(`${baseUrl}/v1/auth/refresh`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ refresh_token: _refreshToken })
	});

	if (!res.ok) {
		clearTokens();
		throw new Error('Session expired');
	}

	const data = await res.json();
	if (data?.token) setToken(data.token);
	if (data?.refresh_token) setRefreshToken(data.refresh_token);
	return data;
};

export async function apiFetch<T>(
	path: string,
	options: RequestInit = {},
	retry = true
): Promise<T> {
	const headers = new Headers(options.headers ?? {});
	if (!headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}
	if (_token) {
		headers.set('Authorization', `Bearer ${_token}`);
	}

	const res = await fetch(`${baseUrl}${path}`, {
		...options,
		headers
	});

	if (res.status === 401 && retry && _refreshToken) {
		await refreshSession();
		return apiFetch<T>(path, options, false);
	}

	if (!res.ok) {
		const err = await res.json().catch(() => ({}));
		const message = err?.error?.message ?? 'Request failed';
		throw new Error(message);
	}

	if (res.status === 204) {
		return {} as T;
	}
	return res.json();
}

export async function getMe() {
	return apiFetch<{ user: { id: string } }>('/v1/me');
}

export type PresignResponse = {
	signed_url: string;
	token: string;
	path: string;
};

export const getPublicPhotoURL = (path: string) => {
	const supabaseUrl = env.PUBLIC_SUPABASE_URL;
	const bucket = env.PUBLIC_SUPABASE_BUCKET;
	if (!supabaseUrl || !bucket || !path) return '';
	const cleanPath = path.startsWith(`${bucket}/`) ? path.slice(bucket.length + 1) : path;
	return `${supabaseUrl}/storage/v1/object/public/${bucket}/${cleanPath}`;
};

export async function presignTripPhoto(tripId: string, fileName: string, contentType: string) {
	return apiFetch<PresignResponse>(`/v1/trips/${tripId}/photos/presign`, {
		method: 'POST',
		body: JSON.stringify({ file_name: fileName, content_type: contentType })
	});
}

export async function uploadSignedPhoto(signedUrl: string, token: string, file: File) {
	const res = await fetch(signedUrl, {
		method: 'PUT',
		headers: {
			'Content-Type': file.type,
			Authorization: `Bearer ${token}`
		},
		body: file
	});
	if (!res.ok) {
		throw new Error('Upload failed');
	}
}
