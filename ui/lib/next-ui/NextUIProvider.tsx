"use client";

import { ReactNode } from "react";
import { NextUIProvider as NUP } from "@nextui-org/system";
import { useRouter } from "next/navigation";

export default function NextUIProvider({ children }: { children: ReactNode }) {
  const router = useRouter();

  return <NUP navigate={router.push}>{children}</NUP>;
}
