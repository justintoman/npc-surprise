import { Outlet } from 'react-router-dom';
import { ThemeModeToggle } from '~/ThemeModeToggle';

export function App() {
  return (
    <main className="h-full w-full bg-background">
      <div className="mx-auto max-w-7xl">
        <header className="flex justify-between p-4">
          <h1 className="text-xl font-bold">NPC Surprise 🧙‍♂️🪄</h1>
          <ThemeModeToggle />
        </header>
        <div className="mx-auto mt-8">
          <Outlet />
        </div>
      </div>
    </main>
  );
}
