import { Workspace } from "@/interfaces/Workspace";
import Link from "next/link";
import { SlArrowRight } from "react-icons/sl";

export interface WorkspaceItemProps {
  workspace: Workspace;
}

export default function WorkspaceItem({ workspace }: WorkspaceItemProps) {
  return (
    <Link
      href={"/"}
      className="flex flex-col border border-gray-200 rounded-2xl hover:bg-gray-200 transition-colors"
    >
      <div className="px-5 py-4">
        <div className="flex items-center justify-between">
          <div>{workspace.name}</div>
          <div>
            <SlArrowRight className="text-md" />
          </div>
        </div>
      </div>
    </Link>
  );
}
