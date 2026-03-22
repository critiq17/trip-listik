import { h as head, a as attr } from "../../chunks/index.js";
const favicon = "data:image/svg+xml,%3csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%20100%20100'%20width='100'%20height='100'%3e%3c!--%20background%20--%3e%3crect%20width='100'%20height='100'%20rx='22'%20fill='%232c313a'/%3e%3c!--%20map%20pin%20--%3e%3ccircle%20cx='50'%20cy='26'%20r='8.5'%20fill='none'%20stroke='%235cba7d'%20stroke-width='4.5'/%3e%3cpath%20d='M50%2033%20Q44%2042%2050%2050%20Q56%2042%2050%2033Z'%20fill='%235cba7d'/%3e%3c!--%20mountains%20--%3e%3cpolyline%20points='10,75%2034,45%2050,62%2066,45%2090,75'%20fill='none'%20stroke='%235cba7d'%20stroke-width='5'%20stroke-linecap='round'%20stroke-linejoin='round'%20/%3e%3c/svg%3e";
function _layout($$renderer, $$props) {
  let { children } = $$props;
  head("12qhfyh", $$renderer, ($$renderer2) => {
    $$renderer2.push(`<link rel="icon"${attr("href", favicon)}/>`);
  });
  children($$renderer);
  $$renderer.push(`<!---->`);
}
export {
  _layout as default
};
