"use client";

import MainContainer from "@/components/templates/MainContainer";
import Navbar from "@/components/ui/Navbar";
import SideBar from "@/components/ui/SideBar";
import { useState } from "react";

export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [isSideBarOpen, setIsSideBarOpen] = useState<boolean>(true);

  const sideBarToggle = () => {
    setIsSideBarOpen((prev) => !prev);
  };

  return (
    <div className="grid grid-cols-12">
      <div className="col-span-2">
        <SideBar isOpen={isSideBarOpen} />
      </div>
      <div className={isSideBarOpen ? "col-span-10" : "col-span-12"}>
        <Navbar sideBarToggle={sideBarToggle} />
        <MainContainer>{children}</MainContainer>
      </div>
    </div>
  );
}
