import { animate } from 'motion';

export function fadeUp(node: HTMLElement, { delay = 0, duration = 0.35 } = {}) {
	animate(
		node,
		{ y: ['18px', '0px'], opacity: [0, 1] } as any,
		{ duration, delay, easing: 'ease-out' } as any
	);
}

export function scalePress(node: HTMLElement) {
	const down = () => animate(node, { scale: [1, 0.97] } as any, { duration: 0.1 });
	const up = () =>
		animate(
			node,
			{ scale: [0.97, 1] } as any,
			{ duration: 0.15, easing: [0.34, 1.56, 0.64, 1] } as any
		);
	node.addEventListener('pointerdown', down);
	node.addEventListener('pointerup', up);
	node.addEventListener('pointercancel', up);
	node.addEventListener('pointerleave', up);
	return {
		destroy() {
			node.removeEventListener('pointerdown', down);
			node.removeEventListener('pointerup', up);
			node.removeEventListener('pointercancel', up);
			node.removeEventListener('pointerleave', up);
		}
	};
}

export function staggerList(node: HTMLElement) {
	const items = node.querySelectorAll<HTMLElement>('[data-item]');
	animate(
		[...items],
		{ y: ['16px', '0px'], opacity: [0, 1] } as any,
		{ duration: 0.32, delay: (i: number) => i * 0.055, easing: 'ease-out' } as any
	);
}
