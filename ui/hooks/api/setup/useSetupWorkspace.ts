import { CreateWorkspaceFormType, CreateWorkspaceRequest, CreateWorkspaceResponse } from "@/interfaces/Workspace";
import { useMutation } from "@tanstack/react-query";
import axios from "@/lib/axios/axios.config";
import { getSession } from "next-auth/react";

const setupWorkspaceFn = async (data: CreateWorkspaceFormType) => {
    const session = await getSession();

    const payload: CreateWorkspaceRequest = {
        name: data.name,
    }

    const response = await axios.post<CreateWorkspaceResponse>('/setup/v1/workspace', payload, {
        headers: {
            Authorization: `Bearer ${session?.user?.token}`,
        }
    });

    return response.data;
}

const useSetupWorkspace = () => {
    return useMutation({
        mutationFn: setupWorkspaceFn,
    });
}

export default useSetupWorkspace;