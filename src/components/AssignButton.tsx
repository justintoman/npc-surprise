import { useAtomValue } from 'jotai';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';
import { playersAtom } from '~/state';

type Props = {
  id: number;
  type: 'action' | 'character';
};

export function AssignButton({ id, type }: Props) {
  const players = useAtomValue(playersAtom);
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="secondary">
          {type === 'character' ? 'Assign' : 'Reveal'}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        {players.map((player) => (
          <DropdownMenuItem
            key={player.id}
            onClick={() => NpcSurpriseApi.assign(type, id, player.id)}
          >
            {player.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
