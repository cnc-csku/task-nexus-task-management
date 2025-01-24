const workspaceQueryKeys = {
    all: ['Workspaces'],
    my: (userId: string) => ['Workspaces', 'My', userId],
    byId: (id: string) => [...workspaceQueryKeys.all, id],
};

export default workspaceQueryKeys;