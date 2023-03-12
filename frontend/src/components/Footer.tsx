import { AiFillGithub, AiFillTwitterCircle } from "react-icons/ai";

export function Footer() {
  return (
    <footer className="flex w-full flex-row justify-end py-2 px-2">
      <a
        href="https://github.com/fjnkt98/AtCoderSearch"
        target="_blank"
        rel="noreferrer"
        className="mx-1"
      >
        <AiFillGithub size="1.5rem" />
      </a>
      <a
        href="https://twitter.com/fjnkt98_jp"
        target="_blank"
        rel="noreferrer"
        className="mx-1"
      >
        <AiFillTwitterCircle size="1.5rem" />
      </a>
    </footer>
  );
}
