import { useAtomValue } from 'jotai';
import { Edit } from 'lucide-react';
import { useState } from 'react';
import { ActionForm } from '~/AdminView/ActionForm';
import { AssignActionButton } from '~/components/AssignActionButton';
import { Button } from '~/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '~/components/ui/dialog';
import { isAdminAtom } from '~/state';
import type { Action } from '~/types';

export function Action({ action }: { action: Action }) {
  const isAdmin = useAtomValue(isAdminAtom);
  const [isEditing, setIsEditing] = useState(false);
  return (
    <div>
      {isAdmin ? (
        <div className="flex items-center justify-start space-x-4">
          <AssignActionButton id={action.id} player_id={action.player_id} />
          <Dialog open={isEditing} onOpenChange={setIsEditing}>
            <DialogTrigger asChild>
              <Button size="icon">
                <Edit />
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Action</DialogTitle>
              </DialogHeader>
              <ActionForm
                id={action.id}
                character_id={action.character_id}
                defaultValues={action}
                onClose={() => setIsEditing(false)}
              />
            </DialogContent>
          </Dialog>
        </div>
      ) : null}
      <p>{action.type}</p>
      <p>{action.direction}</p>
      <p>{action.content}</p>
    </div>
  );
}
