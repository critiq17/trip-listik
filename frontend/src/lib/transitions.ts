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
 * flipNumber: CSS-powered number flip animation.
 * Add to an element when value changes.
 */
export function flipNumber(el: HTMLElement, value: number) {
	el.setAttribute('data-prev', el.textContent ?? '');
	el.textContent = value.toString();
	el.style.animation = 'none';
	// Force reflow
	void el.offsetHeight;
	el.style.animation = 'flipIn 0.3s var(--transition-spring) both';
}
