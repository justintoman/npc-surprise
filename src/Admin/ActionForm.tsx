import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { ActionMarkdown } from '~/components/ActionMarkdown';
import { Button } from '~/components/ui/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { Textarea } from '~/components/ui/textarea';

type Props = {
  defaultValues: ActionFormInput;
  onClose(): void;
  onSubmit(data: ActionFormInput): void;
};

export function ActionForm({ defaultValues, onClose, onSubmit }: Props) {
  const methods = useForm<ActionFormInput>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  const content = methods.watch('content');

  return (
    <div className="grid grid-cols-2 gap-10">
      <Form {...methods}>
        <form
          onSubmit={methods.handleSubmit(onSubmit)}
          className="space-y-8 p-4"
        >
          <FormField
            control={methods.control}
            name="content"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Content</FormLabel>
                <FormControl>
                  <Textarea rows={15} className="font-mono" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="space-x-4">
            <Button variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">Submit</Button>
          </div>
        </form>
      </Form>
      <div>
        <h3 className="mb-8">Preview</h3>
        <div className="max-w-sm bg-secondary p-4">
          <ActionMarkdown>{content}</ActionMarkdown>
        </div>
      </div>
    </div>
  );
}

type ActionFormInput = z.infer<typeof schema>;

const schema = z.object({
  content: z.string().min(1, 'Content is required'),
});
