import { useAtomValue } from 'jotai';
import { Edit } from 'lucide-react';
import { useState } from 'react';
import { ActionForm } from '~/AdminView/ActionForm';
import { RevealActionButton } from '~/components/RevealActionButton';
import { Button } from '~/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '~/components/ui/dialog';
import { isAdminAtom, statusAtom } from '~/state';
import type { Action } from '~/types';

export function Action({ action }: { action: Action }) {
  const isAdmin = useAtomValue(isAdminAtom);
  const status = useAtomValue(statusAtom);
  console.log({ isAdmin, status });
  const [isEditing, setIsEditing] = useState(false);
  return (
    <div>
      {isAdmin ? (
        <div className="flex items-center justify-start space-x-4">
          <RevealActionButton
            actionId={action.id}
            characterId={action.characterId}
            revealed={action.revealed}
          />
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
                characterId={action.characterId}
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
