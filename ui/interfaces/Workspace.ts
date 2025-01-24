import { z } from "zod";

export interface Workspace {
    id: string;
    name: string;
    members: WorkspaceMember[];
    createdBy: string
    createdAt: Date;
    updatedAt: Date
}

export interface WorkspaceMember {
    userId: string;
    name: string;
    role: string;
    joinedAt: Date;
}

export const CreateWorkspaceFormSchema = z.object({
    name: z.string().min(1),
});

export type CreateWorkspaceFormType = z.infer<typeof CreateWorkspaceFormSchema>;

export interface CreateWorkspaceRequest {
    name: string;
}

export interface CreateWorkspaceResponse {
    id: string;
    name: string;
    members: WorkspaceMember[];
    createdBy: string;
    createdAt: string;
    updatedAt: string;
}

