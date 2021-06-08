<script lang="ts">
    // import { MONTHS, DAYS_OF_WEEK } from "./DateUtils";
    import type { PrayerCall } from "./models";

    export let prayer: PrayerCall;
    export let nextPrayerIndex: number;

    function getDisplayDate(dateStr: string): string {
        const date = new Date(dateStr);
        // const dayOfWeek = DAYS_OF_WEEK[date.getDay()].substring(0, 3);
        // const month = MONTHS[date.getMonth()].substring(0, 3);
        // return `${dayOfWeek} ${date.getDate()} ${month}`;
        return date.toDateString();
    }

    function getDisplayTime(dateStr: string): string {
        const date = new Date(dateStr);
        const timeString = date.toLocaleTimeString("en-US", {
            hour: "numeric",
            minute: "numeric",
            hour12: true,
        });
        return timeString;
    }

    async function handleToggleAdhan(index: number) {
        const res = await fetch(`/api/timings/toggle/${index}`, {
            method: "POST",
        });
        const updatedPrayer = await res.json();
        if (res.ok) {
            // mutate prayer -> re-render
            prayer = updatedPrayer;
            return updatedPrayer;
        } else {
            throw new Error("failed to toggle adhan");
        }
    }
</script>

<tr class:isnext={prayer.index === nextPrayerIndex}>
    <td>{getDisplayDate(prayer.time)}</td>
    <td>{prayer.type}</td>
    <td>{getDisplayTime(prayer.time)}</td>
    <td
        class="clickable "
        class:on={prayer.play}
        class:off={!prayer.play}
        on:click={() => handleToggleAdhan(prayer.index)}
    >
        {prayer.play ? "ON" : "OFF"}
    </td>
</tr>

<style>
    tr:nth-child(odd) {
        background-color: #dadada;
    }
    tr:nth-child(even) {
        background-color: #eaeaea;
    }

    .isnext {
        background-color: yellow !important;
    }
    .clickable {
        cursor: pointer;
        border-radius: 0.25em;
        text-align: center;
    }
    .on {
        background: #75d37e;
        border: 1px solid grey;
        color: #fff;
    }
    .off {
        background: rgb(197, 39, 39);
        border: 1px solid grey;
        color: #fff;
    }
</style>
