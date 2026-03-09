import type { UserStats } from '$lib/types';

export const formatShortDate = (value?: string | null) => {
	if (!value) return 'Dates pending';
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return 'Dates pending';
	return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
};

export const formatDateRange = (start?: string | null, end?: string | null) => {
	if (!start && !end) return 'Dates pending';
	if (start && end) {
		return `${formatShortDate(start)} - ${formatShortDate(end)}`;
	}
	return formatShortDate(start ?? end);
};

export const formatLongDate = (value?: string | null) => {
	if (!value) return 'Not set';
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return 'Not set';
	return date.toLocaleDateString(undefined, {
		weekday: 'short',
		month: 'short',
		day: 'numeric',
		year: 'numeric'
	});
};

export const formatRelativeDate = (value?: string | null) => {
	if (!value) return 'just now';
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return 'just now';

	const diff = date.getTime() - Date.now();
	const minutes = Math.round(diff / 60000);
	const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' });
	if (Math.abs(minutes) < 60) return rtf.format(minutes, 'minute');
	const hours = Math.round(minutes / 60);
	if (Math.abs(hours) < 24) return rtf.format(hours, 'hour');
	const days = Math.round(hours / 24);
	if (Math.abs(days) < 30) return rtf.format(days, 'day');
	const months = Math.round(days / 30);
	return rtf.format(months, 'month');
};

export const getUserInitials = (first?: string, last?: string, username?: string) => {
	const source = [first, last].filter(Boolean).join(' ').trim() || username || 'T';
	return source
		.split(/\s+/)
		.slice(0, 2)
		.map((part) => part[0]?.toUpperCase() ?? '')
		.join('');
};

export const getUserName = (user: {
	first_name?: string;
	last_name?: string;
	username?: string;
}) => {
	const fullName = [user.first_name, user.last_name].filter(Boolean).join(' ').trim();
	if (fullName) return fullName;
	if (user.username) return `@${user.username}`;
	return 'Traveler';
};

export const normalizeStats = (stats: Record<string, unknown> | null | undefined): UserStats => ({
	total_trips: Number(stats?.total_trips ?? stats?.TotalTrips ?? 0),
	countries_visited: Number(stats?.countries_visited ?? stats?.CountriesVisited ?? 0),
	cities_visited: Number(stats?.cities_visited ?? stats?.CitiesVisited ?? 0),
	trips_with_friends: Number(stats?.trips_with_friends ?? stats?.TripsWithFriends ?? 0),
	solo_trips: Number(stats?.solo_trips ?? stats?.SoloTrips ?? 0)
});

export const getTripLocation = (trip: { city?: string; country_code?: string }) => {
	return [trip.city, trip.country_code].filter(Boolean).join(', ') || 'Location pending';
};

export const getStatusLabel = (status?: string) => {
	switch (status) {
		case 'completed':
			return 'Completed';
		case 'draft':
			return 'Draft';
		case 'planned':
			return 'Planned';
		case 'canceled':
			return 'Canceled';
		default:
			return 'Planned';
	}
};

export const parseNotificationPayload = (
	payload: string | Record<string, unknown> | null | undefined
) => {
	if (!payload) return {};
	if (typeof payload === 'string') {
		try {
			return JSON.parse(payload) as Record<string, unknown>;
		} catch {
			return {};
		}
	}
	return payload;
};
