"use client";

import { ReactNode } from "react";
import { HeroUIProvider as NUP } from "@heroui/system";
import { useRouter } from "next/navigation";

export default function HeroUIProvider({ children }: { children: ReactNode }) {
  const router = useRouter();

  return <NUP navigate={router.push}>{children}</NUP>;
}
