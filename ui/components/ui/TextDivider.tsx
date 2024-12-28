import { twMerge } from "tailwind-merge";

export interface TextDividerProps {
  text: string;
  className?: string;
}

export default function TextDivider({ text, className }: TextDividerProps) {
  return (
    <div className={twMerge("flex items-center w-full space-x-5", className)}>
      <hr className="w-full" />
      <span className="mx-2">{text}</span>
      <hr className="w-full" />
    </div>
  );
}
