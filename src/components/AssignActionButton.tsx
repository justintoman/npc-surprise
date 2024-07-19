import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';

type Props = {
  playerId?: number;
  id: number;
};

export function AssignActionButton({ id, playerId }: Props) {
  if (!playerId) {
    return null;
  }
  return (
    <Button
      size="sm"
      variant="secondary"
      onClick={() => NpcSurpriseApi.assign('action', id, playerId)}
    >
      Reveal
    </Button>
  );
}
