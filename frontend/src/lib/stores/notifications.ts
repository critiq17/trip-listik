import { writable } from 'svelte/store';

/** Unread notification count. Polled every 30 seconds. */
export const unreadCount = writable(0);
