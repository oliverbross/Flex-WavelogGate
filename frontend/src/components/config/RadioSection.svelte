<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
  export let profile;
  export let radioEnabled = false;
</script>

<section class="bg-surface-section border border-stroke-section rounded-lg px-4 py-3">
  <div class="flex items-center justify-between {radioEnabled ? 'mb-3' : ''}">
    <div class="text-2xs text-fg-bright font-semibold uppercase tracking-wider pl-2 border-l-2 border-accent-orange">
      xCAT Radio Control
    </div>
    <label class="text-fg-label text-xs">
      <input
        type="checkbox"
        checked={radioEnabled}
        on:change={(e) => dispatch("enablechange", e.currentTarget.checked)}
      />
      Enable
    </label>
  </div>

  {#if radioEnabled}
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <label class="w-field-xs flex-shrink-0 text-fg-label text-2xs" for="xcat-host">Host</label>
        <input
          id="xcat-host"
          type="text"
          class="flex-none w-field-sm"
          value={profile.xcat_host || "127.0.0.1"}
          on:change={(e) => dispatch("fieldchange", { key: "xcat_host", value: e.currentTarget.value })}
        />
        <label class="text-fg-label text-2xs ml-1 cursor-default" for="xcat-port">RigCtlD port</label>
        <input
          id="xcat-port"
          type="text"
          class="flex-none w-field-xs"
          value={profile.xcat_port || "4532"}
          on:change={(e) => dispatch("fieldchange", { key: "xcat_port", value: e.currentTarget.value })}
        />
      </div>

      <label class="text-fg-label text-xs">
        <input
          type="checkbox"
          checked={profile.wavelog_pmode}
          on:change={(e) => dispatch("fieldchange", { key: "wavelog_pmode", value: e.currentTarget.checked })}
        />
        Set mode on QSY from Wavelog
      </label>

      <div class="bg-surface-card border border-stroke-section rounded p-2 text-2xs text-fg-muted leading-relaxed">
        Start xCAT first and enable its <span class="text-fg-bright">RigCtlD</span> protocol.
        Use port <span class="font-mono text-accent-value">4532</span>, not xCAT's CAT port 5002.
        Flex-WavelogGate reads the active slice frequency and mode and reports changes to the paired Wavelog station.
      </div>
    </div>
  {/if}
</section>
