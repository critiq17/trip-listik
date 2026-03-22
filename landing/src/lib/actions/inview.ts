export function inview(node: HTMLElement, options = { threshold: 0.15 }) {
	const observer = new IntersectionObserver(
		(entries) => {
			entries.forEach((e) => {
				if (e.isIntersecting) {
					node.dispatchEvent(new CustomEvent('enter'));
					observer.unobserve(node);
				}
			});
		},
		options
	);
	observer.observe(node);
	return { destroy: () => observer.disconnect() };
}
