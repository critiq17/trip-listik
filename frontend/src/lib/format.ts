import type { UserStats } from '$lib/types';

/**
 * Parse a calendar date without timezone shifts.
 * Backend sends dates as "2026-07-15" or "2026-07-15T00:00:00Z"; passing those
 * to `new Date(...)` interprets them as UTC midnight, which renders as the
 * previous day for viewers west of UTC. Extract Y-M-D and build a local Date.
 */
export const parseCalendarDate = (value?: string | null): Date | null => {
	if (!value) return null;
	const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(value);
	if (match) {
		return new Date(Number(match[1]), Number(match[2]) - 1, Number(match[3]));
	}
	const fallback = new Date(value);
	return Number.isNaN(fallback.getTime()) ? null : fallback;
};

/**
 * Format a Date as "YYYY-MM-DD" using local calendar fields.
 * `toISOString()` must not be used here: it converts to UTC first, so for
 * UTC+2/+3 users a locally selected July 15 becomes July 14.
 */
export const toCalendarDateString = (d?: Date | null): string => {
	if (!d) return '';
	const year = d.getFullYear();
	const month = String(d.getMonth() + 1).padStart(2, '0');
	const day = String(d.getDate()).padStart(2, '0');
	return `${year}-${month}-${day}`;
};

export const formatShortDate = (value?: string | null) => {
	const date = parseCalendarDate(value);
	if (!date) return 'Dates pending';
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
	const date = parseCalendarDate(value);
	if (!date) return 'Not set';
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
