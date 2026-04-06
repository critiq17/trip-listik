export type TripCardData = {
	id: string;
	owner_id?: string;
	title: string;
	description?: string;
	cover_photo_url?: string;
	start_date?: string | null;
	end_date?: string | null;
	member_count?: number;
	vote_count?: number;
	vote_average?: number;
	country_code?: string;
	city?: string;
	status?: string;
	visibility?: string;
	created_at?: string;
	updated_at?: string;
};

export type User = {
	id: string;
	username?: string;
	first_name?: string;
	last_name?: string;
	photo_url?: string;
	bio?: string;
	is_public?: boolean;
};

export type UserStats = {
	total_trips: number;
	countries_visited: number;
	cities_visited: number;
	trips_with_friends: number;
	solo_trips: number;
};

export type CountryVisit = {
	code: string;
	visit_count: number;
	last_visit: string;
};

export type Member = {
	trip_id: string;
	user_id: string;
	role: string;
	joined_at: string;
	username?: string;
	first_name?: string;
	last_name?: string;
	photo_url?: string;
};

export type JoinRequest = {
	trip_id: string;
	user_id: string;
	status: string;
	created_at: string;
	username?: string;
	first_name?: string;
	last_name?: string;
	photo_url?: string;
};

export type Photo = {
	id: string;
	trip_id: string;
	user_id: string;
	url: string;
	storage_path?: string;
	created_at: string;
};

export type NotificationItem = {
	id: string;
	type: string;
	trip_id?: string | null;
	actor_id?: string | null;
	payload?: string | Record<string, unknown> | null;
	read_at?: string | null;
	created_at: string;
};

export type UserSummary = {
	id: string;
	username?: string;
	first_name?: string;
	last_name?: string;
	photo_url?: string;
};

export type InviteItem = {
	id: string;
	trip_id: string;
	invited_user_id: string;
	invited_by_user_id: string;
	status: string;
	comment?: string | null;
	created_at: string;
	updated_at: string;
	trip_title?: string;
	trip_cover_photo_url?: string;
	inviter_username?: string;
	inviter_first_name?: string;
	inviter_last_name?: string;
	inviter_photo_url?: string;
};

export type VoteItemData = {
	id: string;
	trip_id: string;
	category: string;
	title: string;
	description?: string;
	image_url?: string;
	added_by: string;
	vote_count: number;
	has_voted: boolean;
	created_at: string;
};
