export const getTelegramInitData = () => {
	if (typeof window === 'undefined') return '';
	const tg = (window as any).Telegram?.WebApp;
	return tg?.initData ?? '';
};

export const expandTelegram = () => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	tg?.expand?.();
};
