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
      {characters.length === 0 ? (
        <div className="space-y-6 p-4">
          <p className="mb-10 text-center text-2xl">🎉 You're connected! 🎉</p>
          <p>You don't have any characters yet 🥲</p>
          <p className="text-sm">
            Wait here and Justin will assign you a juicy NPC soon.
          </p>
          <p className="text-xs italic text-muted-foreground">
            In the mean time you can have loads of fun by changing the theme
            between dark and light mode!
          </p>
          <p className="text-[0.625rem] italic text-muted-foreground/80">
            Also it's literally the only thing you can do until you get a
            character, so try to make the best of your situation.
          </p>
        </div>
      ) : (
        characters.map((char) => (
          <div key={char.id}>
            <Character character={char} />
          </div>
        ))
      )}
    </div>
  );
}
