import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';

type Props = {
  player_id?: number;
  id: number;
};

export function AssignActionButton({ id, player_id }: Props) {
  if (!player_id) {
    return null;
  }
  return (
    <Button
      size="sm"
      variant="secondary"
      onClick={() => NpcSurpriseApi.assign('action', id, player_id)}
    >
      Reveal
    </Button>
  );
}
