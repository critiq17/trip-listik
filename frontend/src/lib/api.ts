import { env } from '$env/dynamic/public';

const baseUrl = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

export const getToken = () => {
	if (typeof localStorage === 'undefined') return '';
	return localStorage.getItem('tl_token') ?? '';
};

export const setToken = (token: string) => {
	if (typeof localStorage === 'undefined') return;
	localStorage.setItem('tl_token', token);
};

export const getRefreshToken = () => {
	if (typeof localStorage === 'undefined') return '';
	return localStorage.getItem('tl_refresh') ?? '';
};

export const setRefreshToken = (token: string) => {
	if (typeof localStorage === 'undefined') return;
	localStorage.setItem('tl_refresh', token);
};

export const clearTokens = () => {
	if (typeof localStorage === 'undefined') return;
	localStorage.removeItem('tl_token');
	localStorage.removeItem('tl_refresh');
};

const refreshSession = async () => {
	const refreshToken = getRefreshToken();
	if (!refreshToken) {
		throw new Error('Missing refresh token');
	}

	const res = await fetch(`${baseUrl}/v1/auth/refresh`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ refresh_token: refreshToken })
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
	const token = getToken();
	const headers = new Headers(options.headers ?? {});
	if (!headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	const res = await fetch(`${baseUrl}${path}`, {
		...options,
		headers
	});

	if (res.status === 401 && retry && getRefreshToken()) {
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
	if (!env.PUBLIC_SUPABASE_URL || !env.PUBLIC_SUPABASE_BUCKET) return '';
	return `${env.PUBLIC_SUPABASE_URL}/storage/v1/object/public/${env.PUBLIC_SUPABASE_BUCKET}/${path}`;
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
