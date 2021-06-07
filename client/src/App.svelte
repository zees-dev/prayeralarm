
<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>

<script lang="ts">
	import Prayer from "./Prayer.svelte";

	type Adhan = "Fajr" | "Dhuhr" | "Asr" | "Maghrib" | "Isha"
	interface Timing {
		play: boolean
		time: string
		type: Adhan
	}

	export let title: string;

	let month: string;
	async function getPrayerTimings() {
		const res = await fetch("http://localhost:8080/api/timings");
		// const res = await fetch("/api/timings");
		const timings: Timing[] = await res.json();
		// TODO calculate month
		if (res.ok) {
			return timings;
		} else {
			throw new Error("failed to get data");
		}
	}
</script>

<main>
	<h1>{title}</h1>
	<h3>{month}</h3>
	{#await getPrayerTimings()}
		<p>...waiting</p>
	{:then timings}
		<div style="width:80%; margin:auto;">
			<table style="width:100%; border: 2px solid; border-radius: 0.5em;">
				<tr>
					<th>Adhan</th>
					<th>Time</th>
					<th>Status</th>
				</tr>
				{#each timings as timing}
					<Prayer {timing} />
				{/each}
			</table>
		</div>
	{:catch error}
		<p style="color: red">{error.message}</p>
	{/await}
</main>
