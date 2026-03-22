import type { Action } from 'svelte/action';

/**
 * pageEnter action: applies page enter transition (fade + slide up).
 * Used on the root div in +layout.svelte to animate route changes.
 */
export const pageEnter: Action = (node: HTMLElement) => {
	node.style.opacity = '0';
	node.style.transform = 'translateY(6px)';
	node.style.transition = 'opacity 150ms ease-out, transform 150ms ease-out';

	const frame = requestAnimationFrame(() => {
		node.style.opacity = '1';
		node.style.transform = 'translateY(0)';
	});

	return {
		destroy() {
			cancelAnimationFrame(frame);
		}
	};
};

/**
 * countUp: animates a number from 0 to target value.
 */
export function countUp(el: HTMLElement, target: number, duration = 600) {
	const start = performance.now();
	const tick = (now: number) => {
		const elapsed = now - start;
		const progress = Math.min(elapsed / duration, 1);
		// Ease out cubic
		const eased = 1 - Math.pow(1 - progress, 3);
		el.textContent = Math.round(eased * target).toString();
		if (progress < 1) requestAnimationFrame(tick);
	};
	requestAnimationFrame(tick);
}

/**
 * flipNumber action: proper Svelte action that animates number changes.
 * Use as: use:flipNumber={value}
 */
export function flipNumber(
	el: HTMLElement,
	value: number
): { update: (v: number) => void; destroy: () => void } {
	el.textContent = value.toString();
	return {
		update(newValue: number) {
			if (el.textContent === newValue.toString()) return;
			el.setAttribute('data-prev', el.textContent ?? '');
			el.textContent = newValue.toString();
			el.style.animation = 'none';
			void el.offsetHeight;
			el.style.animation =
				'flipIn 0.3s var(--transition-spring, cubic-bezier(0.34,1.56,0.64,1)) both';
		},
		destroy() {}
	};
}
