import { Drama, Speech } from 'lucide-react';
import Markdown from 'react-markdown';

function CharacterSpeech({
  children,
  ...props
}: React.ClassAttributes<HTMLQuoteElement> &
  React.BlockquoteHTMLAttributes<HTMLQuoteElement>) {
  return (
    <blockquote
      className="flex items-start gap-2 rounded-lg bg-yellow-100 p-4 font-serif text-black"
      {...props}
    >
      <Speech className="h-5 w-5 shrink-0" />
      {children}
    </blockquote>
  );
}

function StageDirection({
  children,
  ...props
}: React.ClassAttributes<HTMLHeadingElement> &
  React.BlockquoteHTMLAttributes<HTMLHeadingElement>) {
  return (
    <h2
      className="flex items-start gap-2 rounded-lg bg-purple-100 p-4 font-sans text-sm text-black"
      {...props}
    >
      <Drama className="h-5 w-5 shrink-0" />
      {children}
    </h2>
  );
}

function Title({
  children,
  ...props
}: React.ClassAttributes<HTMLHeadingElement> &
  React.BlockquoteHTMLAttributes<HTMLHeadingElement>) {
  return (
    <h1 className="text-sm italic" {...props}>
      {children}
    </h1>
  );
}

export function ActionMarkdown({ children }: { children: string }) {
  return (
    <div className="space-y-5">
      <Markdown
        components={{
          blockquote: CharacterSpeech,
          h2: StageDirection,
          h1: Title,
        }}
      >
        {children}
      </Markdown>
    </div>
  );
}
