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

export const setupMainButton = (text: string, onClick: () => void, isVisible = true, isActive = true) => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.MainButton) {
		tg.MainButton.text = text;
		tg.MainButton.onClick(onClick);
		if (isVisible) {
			tg.MainButton.show();
		} else {
			tg.MainButton.hide();
		}
		if (isActive) {
			tg.MainButton.enable();
		} else {
			tg.MainButton.disable();
		}
	}
};

export const hideMainButton = () => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	tg?.MainButton?.hide();
};

export const setMainButtonState = (isLoading: boolean) => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.MainButton) {
		if (isLoading) {
			tg.MainButton.showProgress();
			tg.MainButton.disable();
		} else {
			tg.MainButton.hideProgress();
			tg.MainButton.enable();
		}
	}
};

export const setupBackButton = (onClick: () => void) => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.BackButton) {
		tg.BackButton.show();
		tg.BackButton.onClick(onClick);
	}
};

export const hideBackButton = () => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	tg?.BackButton?.hide();
};

export const hapticImpact = (style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft' = 'light') => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	tg?.HapticFeedback?.impactOccurred(style);
};

export const hapticNotification = (type: 'error' | 'success' | 'warning') => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	tg?.HapticFeedback?.notificationOccurred(type);
};
