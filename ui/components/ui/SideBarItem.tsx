"use client";

import { Button } from "@heroui/button";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ReactNode } from "react";

interface SideBarItemProps {
  name: string;
  href: string;
  color?: "primary" | "default" | "secondary" | "success" | "warning" | "danger";
  startIcon?: ReactNode;
}

export default function SideBarItem({ name, href, color, startIcon }: SideBarItemProps) {
  const path = usePathname();

  return (
    <Button
      className="w-full justify-start text-[15px]"
      as={Link}
      href={href}
      variant={path === href ? "flat" : "light"}
      color={color || "default"}
      startContent={startIcon}
    >
      {name}
    </Button>
  );
}
