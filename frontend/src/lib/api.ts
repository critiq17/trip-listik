import { env } from '$env/dynamic/public';

const baseUrl = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

export const getToken = () => {
	if (typeof localStorage === 'undefined') return '';
	return localStorage.getItem('wandr_token') ?? '';
};

export const setToken = (token: string) => {
	if (typeof localStorage === 'undefined') return;
	localStorage.setItem('wandr_token', token);
};

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
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

	if (!res.ok) {
		const err = await res.json().catch(() => ({}));
		const message = err?.error?.message ?? 'Request failed';
		throw new Error(message);
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
