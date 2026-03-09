import { getToken } from '$lib/api';
import { env } from '$env/dynamic/public';

const baseUrl = (env.PUBLIC_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '');

export const connectTripStream = (tripId: string) => {
	const token = getToken();
	const url = `${baseUrl}/v1/trips/${tripId}/stream?token=${encodeURIComponent(token)}`;
	return new EventSource(url);
};
