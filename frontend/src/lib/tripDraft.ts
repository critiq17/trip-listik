const key = 'tl_trip_draft';

type TripDraft = {
	id: string;
};

export const setDraft = (id: string) => {
	if (typeof localStorage === 'undefined') return;
	const data: TripDraft = { id };
	localStorage.setItem(key, JSON.stringify(data));
};

export const getDraft = (): TripDraft | null => {
	if (typeof localStorage === 'undefined') return null;
	const raw = localStorage.getItem(key);
	if (!raw) return null;
	try {
		return JSON.parse(raw) as TripDraft;
	} catch {
		return null;
	}
};

export const clearDraft = () => {
	if (typeof localStorage === 'undefined') return;
	localStorage.removeItem(key);
};
