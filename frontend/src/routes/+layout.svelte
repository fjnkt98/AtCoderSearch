<script>
  import "../app.css";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiFillGithub from "svelte-icons-pack/ai/AiFillGithub";
  import AiFillTwitterCircle from "svelte-icons-pack/ai/AiFillTwitterCircle";
</script>

<slot />

<footer class="flex h-8 w-full flex-row items-center justify-end p-0.5">
  <a class="mx-1 px-1" href="https://github.com/fjnkt98/AtCoderSearch" target="_blank" rel="noreferrer">
    <Icon src={AiFillGithub} size={"1.5em"} />
  </a>
  <a class="mx-0.5 px-0.5" href="https://twitter.com/fjnkt98_jp" target="_blank" rel="noreferrer">
    <Icon src={AiFillTwitterCircle} size={"1.5em"} />
  </a>
</footer>

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
  }
</style>
