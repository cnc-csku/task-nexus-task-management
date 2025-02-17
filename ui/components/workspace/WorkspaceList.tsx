import { Workspace } from "@/interfaces/Workspace";
import WorkspaceItem from "./WorkspaceItem";

interface WorkspaceListProps {
  workspaces: Workspace[];
}

export default function WorkspaceList({ workspaces }: WorkspaceListProps) {
  return (
    <div className="flex flex-col gap-3">
        {workspaces.map((workspace) => {
          return (
            <WorkspaceItem
              key={workspace.id}
              workspace={workspace}
            />
          );
        })}
    </div>
  );
}
