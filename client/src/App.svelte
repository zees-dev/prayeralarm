
<style>
	main { max-width: 240px; padding: 0.5em; margin: 0 auto; }

	h1 { text-align: center; color: #ff3e00; text-transform: uppercase; font-size: 4em; font-weight: 100; }

	.calendar-subtitle { display: grid; grid-gap: 0.5em; grid-template: 1fr 1fr / repeat(6, 1fr); padding: 0.25em; }

	.subtitle { margin: 0em; text-align: center; grid-column: 3/5; }

	table { width:100%; border: 2px solid; border-radius: 0.5em; padding: 0.5em; }

	.button { height: min-content; margin: 0em; padding: 0.5rem; color: #FFF; grid-row: 2; }
	.on { background: lime; border-color: green; grid-column: 4;}
	.off { background: red; border-color: brown; grid-column: 5;}

	@media (min-width: 640px) {
		main {
			max-width: 800px;
		}
	}
</style>

<script lang="ts">
	import { MONTHS } from "./DateUtils";
	import type { Timing } from "./models"
	import Prayer from "./Prayer.svelte";

	export let title: string;
	let prayerPromise: Promise<Timing[]> = getPrayerTimings();
	let calendarTitle: string = "";
	let nextPrayerIndex: number = -1;

	async function getPrayerTimings() {
		// const res = await fetch("http://localhost:8080/api/timings");
		const res = await fetch("/api/timings");
		const timings: Timing[] = await res.json();
		if (res.ok) {
			calendarTitle = updateMonthName(timings);
			nextPrayerIndex = getNextPrayerIndex(timings);
			return timings;
		} else {
			throw new Error("failed to get data");
		}
	}

	async function setAllPrayerCalls(on: boolean) {
		const status = on?"on":"off";
		// const res = await fetch(`http://localhost:8080/api/timings/${status}`, { method: "POST" });
		const res = await fetch(`/api/timings/${status}`, { method: "POST" });
		const timings: Timing[] = await res.json();
		if (res.ok) {
			calendarTitle = updateMonthName(timings);
			nextPrayerIndex = getNextPrayerIndex(timings);
			return timings;
		} else {
			throw new Error("failed to get data");
		}
	}

	function updateMonthName(timings: Timing[]): string {
		const dateStr = timings[0].date;
		const monthVal = MONTHS[new Date(dateStr).getMonth()];
		const yearVal = new Date(dateStr).getFullYear();
		return `${monthVal} - ${yearVal}`;
	}

	function getNextPrayerIndex(timings: Timing[]): null | number {
        const currentDate = new Date();
		const nextPrayer = timings.flatMap(t => t.prayers).find(p => new Date(p.time) > currentDate);
		if (nextPrayer) {
			return nextPrayer.index;
		}
		return null;
    }

</script>

<main>
	<h1>{title}</h1>
	<div class="calendar-subtitle">
		<h3 class="subtitle">{calendarTitle}</h3>
		<button class="button" class:on={true} on:click={() => prayerPromise=setAllPrayerCalls(true)}>ON</button>
		<button class="button" class:off={true} on:click={() => prayerPromise=setAllPrayerCalls(false)}>OFF</button>
	</div>
	{#await prayerPromise}
		<p>...waiting</p>
	{:then timings}
		<div style="width:80%; margin:auto;">
			<table>
				<tr>
					<th>Date</th>
					<th>Adhan</th>
					<th>Time</th>
					<th>Status</th>
				</tr>
				{#each timings as timing}
					<td colspan="4"><hr /></td>
					{#each timing.prayers as prayer}
						<Prayer {prayer} nextPrayerIndex={nextPrayerIndex}/>
					{/each}
				{/each}
			</table>
		</div>
	{:catch error}
		<p style="color: red">{error.message}</p>
	{/await}
</main>
