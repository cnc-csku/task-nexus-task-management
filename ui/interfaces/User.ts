import { z } from 'zod';

export interface UserResponse {
    id: string;
    email: string;
    fullName: string;
    displayName: string;
    createdAt: string;
    updatedAt: string;
    token: string;
    tokenExpireAt: string;
}

export interface UserRegisterRequest {
    fullName: string;
    email: string;
    password: string;
}

export const UserRegisterFormSchema = z
    .object({
        fullName: z.string(),
        email: z.string().email(),
        password: z.string().min(8),
        confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
        path: ["confirmPassword"], // Specifies the field that triggered the error
        message: "Password and Confirm Password must match",
    });

export type UserRegisterFormType = z.infer<typeof UserRegisterFormSchema>;

export const UserLoginFormSchema = z.object({
    email: z.string(),
    password: z.string(),
});

export type UserLoginFormType = z.infer<typeof UserLoginFormSchema>;

export interface UserLoginRequest {
    email: string;
    password: string;
}

export interface UserLoginResponse {
    id: string;
    email: string;
    fullName: string;
    displayName: string;
    profileUrl: string;
    token: string;
    tokenExpireAt: string;
    createdAt: string;
    updatedAt: string;
}
