import { Link } from "@nextui-org/link";

export interface BackButtonProps {
  href: string;
}

export default function BackButton({ href }: BackButtonProps) {
  return <Link href={href} className="mb-3 text-sm">← Back</Link>;
}
