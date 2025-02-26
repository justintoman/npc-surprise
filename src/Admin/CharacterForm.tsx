import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { Button } from '~/components/ui/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { Input } from '~/components/ui/input';
import { Textarea } from '~/components/ui/textarea';

type Props = {
  defaultValues: CharacterFormInput;
  onClose(): void;
  onSubmit(data: CharacterFormInput): void;
};

export function CharacterForm({ defaultValues, onClose, onSubmit }: Props) {
  const methods = useForm<CharacterFormInput>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  return (
    <Form {...methods}>
      <form
        onSubmit={methods.handleSubmit(onSubmit)}
        className="mx-auto max-w-lg space-y-4"
      >
        <FormField
          control={methods.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="race"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Race</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="gender"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Gender</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="age"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Age</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="appearance"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Appearance</FormLabel>
              <FormControl>
                <Textarea {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div className="mt-8 flex justify-end space-x-4">
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit">Submit</Button>
        </div>
      </form>
    </Form>
  );
}

type CharacterFormInput = z.infer<typeof schema>;

const schema = z.object({
  name: z.string(),
  race: z.string(),
  gender: z.string(),
  age: z.string(),
  description: z.string(),
  appearance: z.string(),
});
