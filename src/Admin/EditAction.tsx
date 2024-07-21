import { useAtomValue } from 'jotai';
import { useNavigate, useParams } from 'react-router-dom';
import { ActionForm } from './ActionForm';
import { NpcSurpriseApi } from '~/api';
import { actionAtomFamily } from '~/state';

export function EditAction() {
  const { actionId } = useParams();
  if (!actionId) {
    throw new Error('ActionId is required');
  }
  const action = useAtomValue(actionAtomFamily(Number(actionId)));
  const navigate = useNavigate();
  if (!action) {
    close();
    return null;
  }

  function close() {
    navigate(`/admin`);
  }

  return (
    <ActionForm
      defaultValues={action}
      onClose={() => close()}
      onSubmit={async (data) => {
        await NpcSurpriseApi.updateAction({ ...action, ...data });
        close();
      }}
    />
  );
}
