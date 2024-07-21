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
    <div className="space-y-4 divide-y-2">
      {characters.map((char) => (
        <div key={char.id}>
          <Character character={char} />
        </div>
      ))}
    </div>
  );
}
