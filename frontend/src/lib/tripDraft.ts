export type TripDraft = {
	title: string;
	destination?: string; // e.g. city
	country_code?: string;
	lat?: number;
	lon?: number;
	cover_photo_id?: string;
	cover_photo_url?: string;
	visibility: string;
	start_date?: string;
	end_date?: string;
};

const key = 'tl_trip_draft_v2';

export const setDraft = (data: Partial<TripDraft>) => {
	if (typeof localStorage === 'undefined') return;
	const existing = getDraft() || {};
	const merged = { ...existing, ...data };
	localStorage.setItem(key, JSON.stringify(merged));
};

export const getDraft = (): Partial<TripDraft> | null => {
	if (typeof localStorage === 'undefined') return null;
	const raw = localStorage.getItem(key);
	if (!raw) return null;
	try {
		return JSON.parse(raw) as Partial<TripDraft>;
	} catch {
		return null;
	}
};

export const clearDraft = () => {
	if (typeof localStorage === 'undefined') return;
	localStorage.removeItem(key);
};
