import {Link} from 'react-router-dom';

type Props = {
  isBig: boolean;
};

export function Logo({ isBig }: Props) {
  return (
    <h1 className={`w-full text-${isBig ? 4 : 2}xl`}>
      <Link to="/">AtCoder Search</Link>
    </h1>
  );
}
