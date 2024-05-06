<script lang="ts">
  import Tab from "$lib/Tab.svelte";
</script>

<div class="flex h-dvh flex-col">
  <nav class="flex w-full flex-grow-0 items-center justify-start sm:py-2">
    <a href="/" class="font-roboto mx-4 my-2 text-xl font-medium text-gray-700">AtCoder Search</a>
  </nav>

  <Tab selected={"user"} />
</div>
