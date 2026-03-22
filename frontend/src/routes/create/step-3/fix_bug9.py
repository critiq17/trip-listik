import re

content = open('+page.svelte').read()

# 1. Imports
imports = """
	import type { User } from '$lib/types';
	import { getUserName, getUserInitials } from '$lib/format';
"""
content = content.replace("import { setupMainButton", imports.strip() + "\n\timport { setupMainButton")

# 2. State and logic
mock = """
	// Mock friend chips (visual only — full friends API to be integrated later)
	const mockFriends = [
		{ id: '1', name: 'Mia K.',    avatar: 'https://i.pravatar.cc/40?img=1' },
		{ id: '2', name: 'Sam R.',    avatar: 'https://i.pravatar.cc/40?img=2' },
	];
	let selectedFriends = $state<typeof mockFriends>([...mockFriends]);
	let friendSearch = $state('');

	const removeFriend = (id: string) => {
		selectedFriends = selectedFriends.filter((f) => f.id !== id);
	};
"""

state = """
	let selectedFriends = $state<User[]>([]);
	let friendSearch = $state('');
	let searchResults = $state<User[]>([]);
	let searchTimeout: ReturnType<typeof setTimeout>;

	const removeFriend = (id: string) => {
		selectedFriends = selectedFriends.filter((f) => f.id !== id);
	};

	const onSearchInput = () => {
		clearTimeout(searchTimeout);
		if (!friendSearch.trim()) {
			searchResults = [];
			return;
		}
		searchTimeout = setTimeout(async () => {
			try {
				const res = await apiFetch<{ users: User[] }>(`/v1/users/search?q=${encodeURIComponent(friendSearch.trim())}`);
				// Filter out already selected
				searchResults = (res.users || []).filter(u => !selectedFriends.some(sf => sf.id === u.id));
			} catch (e) {
				console.error('Search failed', e);
			}
		}, 300);
	};

	const addFriend = (user: User) => {
		selectedFriends = [...selectedFriends, user];
		friendSearch = '';
		searchResults = [];
	};
"""
content = content.replace(mock.strip(), state.strip())

# 3. Submit logic
submit_loc = "			// 3. Clear draft and navigate"
invites = """			// 2.5 Send Invites
			for (const friend of selectedFriends) {
				try {
					await apiFetch(`/v1/trips/${trip.id}/invite`, {
						method: 'POST',
						body: JSON.stringify({ user_id: friend.id })
					});
				} catch (e) {
					console.warn('Failed to invite friend:', e);
				}
			}

			// 3. Clear draft and navigate"""
content = content.replace(submit_loc, invites)

# 4. Markup changes
old_chips = """
				{#each selectedFriends as friend (friend.id)}
					<div class="chip">
						<img class="chip-avatar" src={friend.avatar} alt={friend.name} loading="lazy" />
						<span class="chip-name">{friend.name}</span>
						<button
							class="chip-remove"
							aria-label="Remove {friend.name}"
							onclick={() => removeFriend(friend.id)}
						>
"""

new_chips = """
				{#each selectedFriends as friend (friend.id)}
					<div class="chip">
						{#if friend.photo_url}
							<img class="chip-avatar" src={friend.photo_url} alt={getUserName(friend)} loading="lazy" />
						{:else}
							<div class="chip-avatar fallback">{getUserInitials(friend.first_name, friend.last_name, friend.username)}</div>
						{/if}
						<span class="chip-name">{getUserName(friend)}</span>
						<button
							class="chip-remove"
							aria-label="Remove {getUserName(friend)}"
							onclick={() => removeFriend(friend.id)}
						>
"""
content = content.replace(old_chips.strip(), new_chips.strip())

old_search = """
		<div class="friend-search">
			<input
				type="text"
				class="search-input"
				placeholder="Search friends..."
				bind:value={friendSearch}
				autocomplete="off"
			/>
			<span class="material-symbols-outlined search-icon">search</span>
		</div>
"""

new_search = """
		<div class="friend-search">
			<input
				type="text"
				class="search-input"
				placeholder="Search friends..."
				bind:value={friendSearch}
				oninput={onSearchInput}
				autocomplete="off"
			/>
			<span class="material-symbols-outlined search-icon">search</span>
			
			{#if searchResults.length > 0}
				<div class="search-dropdown">
					{#each searchResults as user (user.id)}
						<button class="search-item" onclick={() => addFriend(user)}>
							{#if user.photo_url}
								<img src={user.photo_url} alt={getUserName(user)} loading="lazy" />
							{:else}
								<div class="fallback-avatar">{getUserInitials(user.first_name, user.last_name, user.username)}</div>
							{/if}
							<span>{getUserName(user)}</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>
"""
content = content.replace(old_search.strip(), new_search.strip())

# 5. Add CSS
css = """
	margin-bottom: 24px;
}

.search-dropdown {
	position: absolute;
	top: calc(100% + 4px);
	left: 0;
	right: 0;
	background: #1a2420;
	border: 1px solid rgba(77, 157, 109, 0.1);
	border-radius: 12px;
	overflow: hidden;
	z-index: 10;
	box-shadow: 0 4px 12px rgba(0,0,0,0.5);
}

.search-item {
	display: flex;
	align-items: center;
	gap: 12px;
	width: 100%;
	padding: 12px 16px;
	background: none;
	border: none;
	color: #f1f5f9;
	font-family: 'Inter', sans-serif;
	font-size: 15px;
	text-align: left;
	cursor: pointer;
	border-bottom: 1px solid rgba(255,255,255,0.05);
}
.search-item:last-child {
	border-bottom: none;
}
.search-item img, .fallback-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	flex-shrink: 0;
	object-fit: cover;
}
.fallback-avatar, .chip-avatar.fallback {
	display: flex;
	align-items: center;
	justify-content: center;
	background: #4d9d6d;
	color: white;
	font-weight: 600;
	font-size: 13px;
}
.chip-avatar.fallback {
	width: 24px;
	height: 24px;
	font-size: 10px;
}

"""
content = content.replace("\n\tmargin-bottom: 24px;\n}", css)


open('+page.svelte', 'w').write(content)

print("Done")
