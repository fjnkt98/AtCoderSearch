<script lang="ts">
  import type { Page } from "$lib/page";

  export let selected: Page | null = null;

  interface Tab {
    value: Page;
    href: string;
    title: string;
  }

  const tabs: Tab[] = [
    {
      value: "problem",
      href: "/problem",
      title: "Problem",
    },
    {
      value: "user",
      href: "/user",
      title: "User",
    },
    {
      value: "submission",
      href: "/submission",
      title: "Submission",
    },
  ];
</script>

<div class="flex flex-row justify-stretch">
  {#each tabs as tab}
    <div class={"block flex-1 rounded-sm p-1 text-center font-medium shadow-sm shadow-gray-500 " + (selected == tab.value ? `bg-green-600 text-slate-50` : `bg-gray-100`)}>
      <a href={tab.href}>
        <span class="cursor-pointer">{tab.title}</span>
      </a>
    </div>
  {/each}
</div>
