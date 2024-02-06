<script>
  import { LaunchGame, InitializeApp } from '$wails/go/main/App.js';
  import { Button } from '$lib/components/ui/button';
  import * as Select from '$lib/components/ui/select';
  import { Loader2 } from 'lucide-svelte';
  import { setMode } from 'mode-watcher';
  import LightDarkButton from '$lib/components/LightDarkButton.svelte';
  import AlertDialog from '$lib/components/AlertDialog.svelte';
  import { Quit, EventsOn, LogPrint } from '$wails/runtime/runtime.js';

  setMode('dark');

  let builds = [];
  let build;

  const runtimes = [
    { value: 'flashplayer.exe', label: 'Flash Player' },
    // { value: 'ruffle', label: 'Ruffle Player (Experimental)' }, TODO" implement this
  ];

  let runtime = runtimes[0];

  let disabled = false;
  let showError = false;
  let errorCode = '';

  let debugLogs = [];

  let version;

  EventsOn('infoLog', (event) => {
    debugLogs = [...debugLogs, event];
  });

  function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }

  EventsOn('initialLoad', (event) => {
    LogPrint(JSON.stringify(event));

    builds = Object.keys(event.manifest.builds).map((buildName) => ({
      value: buildName,
      label: capitalizeFirstLetter(buildName),
    }));
    build = builds[0];
    version = event.manifest.currentGameVersion;

    debugLogs = [
      ...debugLogs,

      `Latest SWF version: ${event.manifest.currentGameVersion}`,
      `Latest Launcher version: ${event.manifest.currentLauncherVersion}`,
    ];
  });

  InitializeApp().then(() => {
    debugLogs = [...debugLogs, 'Launcher initialized'];
  });

  function launch() {
    disabled = true;
    LaunchGame(build.value, version, runtime.value)
      .then(() => {
        showError = false;
        Quit();
      })
      .catch((err) => {
        errorCode = err;
        showError = true;
      });
    disabled = false;
  }
</script>

<main
  class="h-screen bg-background text-foreground flex flex-col gap-4 p-4 antialiased select-none font-sans"
>
  <div class="flex justify-between">
    <h1 class="text-foreground font-bold text-2xl">
      Backyard Monsters <br />
      Refitted
    </h1>
    <LightDarkButton />
  </div>
  <div class="grow bg-secondary font-mono p-3">
    <h5>Welcome to Bym Refitted!</h5>
    {#each debugLogs as log}
      <p><small>{log}</small></p>
    {/each}
  </div>

  {#if !version}
    <div role="status">
      <svg
        aria-hidden="true"
        class="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600"
        viewBox="0 0 100 101"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
          fill="currentColor"
        />
        <path
          d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
          fill="currentFill"
        />
      </svg>
      <span class="sr-only">Loading...</span>
    </div>
  {:else}
    <div class="mt-auto w-full flex justify-between">
      <label>SWF Build</label>
      <Select.Root bind:selected={build} portal={null}>
        <Select.Trigger class="w-[180px] rounded">
          <Select.Value placeholder="Select Build" />
        </Select.Trigger>
        <Select.Content>
          <Select.Group>
            {#each builds as build}
              <Select.Item value={build.value} label={build.label}>{build.label}</Select.Item>
            {/each}
          </Select.Group>
        </Select.Content>
        <Select.Input name="build" />
      </Select.Root>
    </div>
    <div class="mt-auto w-full flex justify-between">
      <label>Flash Runtime</label>
      <Select.Root bind:selected={runtime} portal={null}>
        <Select.Trigger class="w-[180px] rounded">
          <Select.Value placeholder="Select Build" />
        </Select.Trigger>
        <Select.Content>
          <Select.Group>
            {#each runtimes as runtime}
              <Select.Item value={runtime.value} label={runtime.label}>{runtime.label}</Select.Item>
            {/each}
          </Select.Group>
        </Select.Content>
        <Select.Input name="runtime" />
      </Select.Root>
    </div>

    <div class="mt-auto w-full flex justify-between">
      <Button variant="outline" class="p-4 rounded w-32" on:click={launch} {disabled}>
        {#if disabled}
          <Loader2 />
        {:else}
          Launch Game
        {/if}
      </Button>
    </div>
  {/if}
  <AlertDialog bind:open={showError} error={errorCode}></AlertDialog>
</main>
