import { ErrorApiResponse } from "@/interfaces/Common";
import axios, { AxiosError } from "axios";

export const getApiErrorMessage = (error: any) => {
    if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError<ErrorApiResponse>;

        return axiosError.response?.data.message;
    }

    return 'An error occurred';
} 