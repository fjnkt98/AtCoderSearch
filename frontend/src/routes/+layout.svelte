<script>
  import "../app.css";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiFillGithub from "svelte-icons-pack/ai/AiFillGithub";
  import AiFillTwitterCircle from "svelte-icons-pack/ai/AiFillTwitterCircle";
</script>

<nav class="flex w-full items-center justify-start border">
  <a href="/" class="mx-4 my-2 font-roboto text-3xl">AtCoder Search</a>
</nav>

<slot />

<footer class="flex h-12 w-full flex-row items-center justify-end p-2">
  <a class="mx-1 px-1" href="https://github.com/fjnkt98/AtCoderSearch" target="_blank" rel="noreferrer">
    <Icon src={AiFillGithub} size={"2em"} />
  </a>
  <a class="mx-1 px-1" href="https://twitter.com/fjnkt98_jp" target="_blank" rel="noreferrer">
    <Icon src={AiFillTwitterCircle} size={"2em"} />
  </a>
</footer>

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
  }
</style>
