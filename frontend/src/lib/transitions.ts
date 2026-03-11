import { animate } from 'motion';

export function pageEnter(node: HTMLElement) {
	animate(
		node,
		{ x: ['20px', '0px'], opacity: [0, 1] } as any,
		{ duration: 0.25, easing: [0.25, 0.46, 0.45, 0.94] } as any
	);
}

export function countUp(node: HTMLElement, target: number, duration = 1.1) {
	let start: number | null = null;
	const tick = (ts: number) => {
		if (!start) start = ts;
		const p = Math.min((ts - start) / (duration * 1000), 1);
		const eased = 1 - Math.pow(1 - p, 3);
		node.textContent = Number.isInteger(target)
			? Math.floor(eased * target).toString()
			: `${(eased * target).toFixed(1)}%`;
		if (p < 1) requestAnimationFrame(tick);
	};
	requestAnimationFrame(tick);
}
