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

// Telegram accumulates MainButton onClick handlers; registering a new one on
// every reactive update makes a single tap fire the callback multiple times
// (duplicate trips, double submits). Track the active handler and detach it
// before attaching the next one.
let activeMainButtonHandler: (() => void) | null = null;

export const setupMainButton = (
	text: string,
	onClick: () => void,
	isVisible = true,
	isActive = true
) => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.MainButton) {
		tg.MainButton.text = text;
		if (activeMainButtonHandler) {
			tg.MainButton.offClick?.(activeMainButtonHandler);
		}
		tg.MainButton.onClick(onClick);
		activeMainButtonHandler = onClick;
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
	if (tg?.MainButton) {
		if (activeMainButtonHandler) {
			tg.MainButton.offClick?.(activeMainButtonHandler);
			activeMainButtonHandler = null;
		}
		tg.MainButton.hide();
	}
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

let activeBackButtonHandler: (() => void) | null = null;

export const setupBackButton = (onClick: () => void) => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.BackButton) {
		if (activeBackButtonHandler) {
			tg.BackButton.offClick?.(activeBackButtonHandler);
		}
		tg.BackButton.show();
		tg.BackButton.onClick(onClick);
		activeBackButtonHandler = onClick;
	}
};

export const hideBackButton = () => {
	if (typeof window === 'undefined') return;
	const tg = (window as any).Telegram?.WebApp;
	if (tg?.BackButton) {
		if (activeBackButtonHandler) {
			tg.BackButton.offClick?.(activeBackButtonHandler);
			activeBackButtonHandler = null;
		}
		tg.BackButton.hide();
	}
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
