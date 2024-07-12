import { useAtomValue } from 'jotai';
import { loadable } from 'jotai/utils';
import { AdminView } from '~/AdminView/AdminView';
import { PlayerLogin } from '~/PlayerLogin/PlayerLogin';
import { PlayerView } from '~/PlayerView/PlayerView';
import { statusAtom } from '~/state';
import { ThemeModeToggle } from '~/ThemeModeToggle';

export function App() {
  const statusLoadable = useAtomValue(loadable(statusAtom));

  const status =
    statusLoadable.state === 'hasData' ? statusLoadable.data : null;

  return (
    <main className="h-full w-full bg-background">
      <div className="mx-auto max-w-md">
        <header className="flex justify-between p-4">
          <h1 className="text-xl font-bold">NPC Surprise 🧙‍♂️🪄</h1>
          <ThemeModeToggle />
        </header>
        <div className="mx-auto mt-8">
          {status?.is_admin ? (
            <AdminView />
          ) : status?.player_id ? (
            <PlayerView />
          ) : (
            <PlayerLogin name={status?.player_name ?? ''} />
          )}
        </div>
      </div>
    </main>
  );
}
