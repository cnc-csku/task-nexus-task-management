import axios from "@/lib/axios/axios.config";
import { useQuery } from "@tanstack/react-query";
import workspaceQueryKeys from "./workspaceQueryKeys";
import { getSession } from "next-auth/react";
import { Workspace } from "@/interfaces/Workspace";


const fetchMyWorkspacesFn = async () => {
    const session = await getSession();

    const response = await axios.get<Workspace[]>('/workspaces/v1/own-workspaces', {
        headers: {
            Authorization: `Bearer ${session?.user?.token}`,
        }
    });

    return response.data;
}

const useMyWorkspaces = () => {
    
    return useQuery({
        queryKey: workspaceQueryKeys.my('m'),
        queryFn: fetchMyWorkspacesFn,
    });
}

export default useMyWorkspaces