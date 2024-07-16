import { useAtomValue } from 'jotai';
import { useEffect } from 'react';
import { Character } from '~/components/Character';
import { charactersAtom, initStream } from '~/state';

export function PlayerView() {
  const characters = useAtomValue(charactersAtom);
  useEffect(() => {
    return initStream();
  }, []);

  return (
    <div>
      <header>
        <h1>Player View</h1>
      </header>

      <h2 className="text-lg font-bold">Assigned Characters</h2>

      <ul>
        {characters.map((char) => (
          <li key={char.id}>
            <Character character={char} />
          </li>
        ))}
      </ul>
    </div>
  );
}
