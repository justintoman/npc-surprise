import { Label } from '~/components/ui/label';

type Props = {
  label: string;
  value: string | undefined;
};

export function CharacterField({ label, value }: Props) {
  return (
    <div className="flex flex-col">
      <Label className="text-sm font-bold">{label}</Label>
      <p className="text-sm text-gray-500">{value || 'hidden'}</p>
    </div>
  );
}
