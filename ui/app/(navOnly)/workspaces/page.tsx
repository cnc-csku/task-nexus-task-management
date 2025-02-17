"use client";

import { motion } from "framer-motion";
import React from "react";
import WorkspaceList from "../../../components/workspace/WorkspaceList";
import useMyWorkspaces from "@/hooks/api/workspace/useMyWorkspaces";
import LoadingScreen from "@/components/ui/LoadingScreen";

export default function WorkspacesListPage() {
  const { data: workspaces, isPending, error } = useMyWorkspaces();

  if (isPending) {
    return <LoadingScreen />;
  }

  if (error) {
    return <p>{error.message}</p>;
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-64px)]">
      <div className="lg:w-[500px]">
        <motion.div
          className="flex flex-col gap-5 px-1"
          initial={{ opacity: 0, x: -100 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8, ease: "easeOut" }}
        >
          <h1 className="text-2xl mb-5">Workspaces</h1>
        </motion.div>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <WorkspaceList workspaces={workspaces} />
        </motion.div>
      </div>
    </div>
  );
}
