import { Moon, Sun, SunMoon } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';
import { useTheme } from './ThemeProvider';

export function ThemeModeToggle() {
  const { theme, setTheme } = useTheme();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        {theme === 'dark' ? (
          <Moon aria-label="dark theme" />
        ) : theme === 'light' ? (
          <Sun aria-label="light theme" />
        ) : theme === 'system' ? (
          <SunMoon aria-label="system theme" />
        ) : null}
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem
          onClick={() => setTheme('light')}
          className="space-x-4"
        >
          <Sun aria-label="light theme" />
          <span>Light</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => setTheme('system')}
          className="space-x-4"
        >
          <SunMoon aria-label="system theme" />
          <span>System</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => setTheme('dark')}
          className="space-x-4"
        >
          <Moon aria-label="dark theme" /> <span>Dark</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
