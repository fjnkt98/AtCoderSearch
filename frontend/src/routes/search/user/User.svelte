<script lang="ts">
  import type { User } from "$lib/search";
  import { textColorStyles } from "$lib/colors";

  export let user: User;
</script>

<div class="my-2 w-full max-w-5xl rounded-2xl bg-white px-4 py-4 shadow-md shadow-gray-300">
  <div class="flex flex-row items-center">
    <p class="mx-3 text-xl">{user.rank}</p>
    <div class="flex-grow">
      <a href={user.user_url} target="_blank" rel="noreferrer" class={`mx-2 text-xl ${textColorStyles.get(user.color)}`}>{user.user_name}</a>
      <div class="flex flex-row items-center">
        {#if user.affiliation != null}
          <p class="mx-2 break-all text-sm text-gray-500">{user.affiliation} {user.country}</p>
        {/if}
        {#if user.birth_year != null}
          <span class="text-sm">{user.birth_year}</span>
        {/if}
      </div>
    </div>
    <div class="flex flex-col items-center justify-center">
      <p>Rating</p>
      <p>{user.rating}</p>
    </div>
  </div>
</div>
