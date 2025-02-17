import { twMerge } from "tailwind-merge";

export default function MainContainer({
  className,
  children,
}: Readonly<{
  className?: string;
  children: React.ReactNode;
}>) {
  return <div className={twMerge("px-5 py-2", className)}>{children}</div>;
}
