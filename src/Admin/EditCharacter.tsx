import { useAtomValue } from 'jotai';
import { useNavigate, useParams } from 'react-router-dom';
import { NpcSurpriseApi } from '~/api';
import { characterAtomFamily } from '~/state';
import { CharacterForm } from './CharacterForm';

export function EditCharacter() {
  const { characterId } = useParams();
  if (!characterId) {
    throw new Error('characterId is required');
  }
  const character = useAtomValue(characterAtomFamily(Number(characterId)));
  const navigate = useNavigate();
  if (!character) {
    close();
    return null;
  }

  function close() {
    navigate(`/admin`);
  }

  return (
    <CharacterForm
      defaultValues={character}
      onClose={() => close()}
      onSubmit={async (data) => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { actions, ...withoutActions } = character;
        await NpcSurpriseApi.updateCharacter({ ...withoutActions, ...data });
        close();
      }}
    />
  );
}
