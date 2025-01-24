import { UserRegisterFormType, UserRegisterRequest, UserResponse } from "@/interfaces/User";
import axios from "@/lib/axios/axios.config";
import { useMutation } from "@tanstack/react-query";

const setupUserFn = async (user: UserRegisterFormType) => {
    const payload: UserRegisterRequest = {
        fullName: user.fullName,
        email: user.email,
        password: user.password,
    }

    const response = await axios.post<UserResponse>('/setup/v1/user', payload);

    return response.data;
}

const useSetupUser = () => {
    return useMutation({
        mutationFn: setupUserFn,
    })
}

export default useSetupUser;