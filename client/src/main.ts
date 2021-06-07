import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
		title: 'Prayer Calendar',
	}
});

export default app;