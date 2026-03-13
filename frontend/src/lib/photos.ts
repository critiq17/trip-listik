import { getPublicPhotoURL } from '$lib/api';

export const resolvePhotoUrl = (url?: string | null) => {
	if (!url) return '';
	if (url.startsWith('http://') || url.startsWith('https://')) return url;
	const publicUrl = getPublicPhotoURL(url);
	return publicUrl || url;
};
