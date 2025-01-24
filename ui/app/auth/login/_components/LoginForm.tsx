"use client";

import { Input } from "@nextui-org/input";
import { Button } from "@nextui-org/button";
import { UserLoginFormSchema, UserLoginFormType } from "@/interfaces/User";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

export interface LoginFormProps {
  onLogin: (data: UserLoginFormType) => void;
  isLoading: boolean;
}

export default function LoginForm({ onLogin, isLoading }: LoginFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<UserLoginFormType>({ resolver: zodResolver(UserLoginFormSchema) });

  return (
    <form
      onSubmit={handleSubmit(onLogin)}
      className="flex flex-col space-y-4"
    >
      <Input
        label="Email"
        type="email"
        isRequired
        isInvalid={!!errors.email}
        errorMessage={errors.email?.message}
        {...register("email")}
      />
      <Input
        label="Password"
        type="password"
        isRequired
        isInvalid={!!errors.password}
        errorMessage={errors.password?.message}
        {...register("password")}
      />
      <Button
        color="primary"
        size="lg"
        isLoading={isLoading}
        type="submit"
      >
        Login
      </Button>
    </form>
  );
}
