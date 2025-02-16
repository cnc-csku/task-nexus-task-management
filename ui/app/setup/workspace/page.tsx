"use client";

import { Divider } from "@heroui/divider";
import { CreateWorkspaceFormType } from "@/interfaces/Workspace";
import { motion } from "framer-motion";
import SetupWorkspaceForm from "./_components/SetupWorkspaceForm";
import useSetupWorkspace from "@/hooks/api/setup/useSetupWorkspace";
import { toast } from "sonner";
import { getApiErrorMessage } from "@/utils/errutils";
import { useRouter } from "next/navigation";

export default function SetupWorkspacePage() {
  const router = useRouter();
  const { isPending, mutateAsync } = useSetupWorkspace();

  const onCreate = async (data: CreateWorkspaceFormType) => {
    try {
      await mutateAsync(data);
    } catch (error) {
      toast.error(getApiErrorMessage(error));
      return;
    }

    router.push("/workspaces");
  };

  return (
    <div className="lg:w-[500px]">
      <motion.div
        className="flex flex-col gap-5 px-1"
        initial={{ opacity: 0, x: -100 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      >
        <div className="flex items-center gap-2">
          <Divider className="inline-block w-5" />
          <h3>Step 2</h3>
        </div>
        <h1 className="text-2xl mb-5">Create a workspace</h1>
      </motion.div>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.5 }}
      >
        <SetupWorkspaceForm
          onCreate={onCreate}
          isLoading={isPending}
        />
      </motion.div>
    </div>
  );
}
