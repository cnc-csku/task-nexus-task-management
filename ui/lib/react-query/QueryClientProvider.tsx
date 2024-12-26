"use client";

import { ReactNode } from "react";
import { QueryClient, QueryClientProvider as QP } from "@tanstack/react-query";

const queryClient = new QueryClient();

export default function QueryClientProvider({ children }: { children: ReactNode }) {
  return <QP client={queryClient}>{children}</QP>;
}
