import { useAtomValue } from 'jotai';
import { Eye, EyeOff } from 'lucide-react';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import { characterRevealedFieldsAtomFamily } from '~/state';
import { CharacterRevealedFields } from '~/types';

type Props = {
  characterId: number;
  field: keyof Omit<CharacterRevealedFields, 'characterId'>;
};

export function RevealFieldButton({ characterId, field }: Props) {
  const fields = useAtomValue(characterRevealedFieldsAtomFamily(characterId));
  if (!fields) {
    console.log(
      'cannot render reveal fields button because there are no fields',
    );
    return null;
  }
  const isRevealed = fields[field];
  return (
    <Button
      size="tiny"
      variant="ghost"
      onClick={() => {
        if (isRevealed) {
          fields[field] = false;
        } else {
          fields[field] = true;
        }
        NpcSurpriseApi.updateRevealedFields(characterId, {
          ...fields,
          [field]: !isRevealed,
        });
      }}
    >
      {isRevealed ? (
        <Eye className="h-3 w-3" />
      ) : (
        <EyeOff className="h-3 w-3" />
      )}
    </Button>
  );
}
