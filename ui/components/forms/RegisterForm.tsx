"use client";

import { UserRegisterFormSchema, UserRegisterFormType } from "@/interfaces/User";
import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

export interface RegisterFormProps {
  onRegister: (data: UserRegisterFormType) => void;
  isLoading: boolean;
}

export default function RegisterForm({ onRegister, isLoading }: RegisterFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<UserRegisterFormType>({
    resolver: zodResolver(UserRegisterFormSchema),
  });

  return (
    <form
      onSubmit={handleSubmit(onRegister)}
      className="flex flex-col gap-4"
    >
      <Input
        label="Name"
        type="text"
        isRequired
        isInvalid={!!errors.fullName}
        errorMessage={errors.fullName?.message}
        {...register("fullName")}
      />
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
        description="Minimum 8 characters"
        {...register("password")}
      />
      <Input
        label="Confirm Password"
        type="password"
        isRequired
        isInvalid={!!errors.confirmPassword}
        errorMessage={errors.confirmPassword?.message}
        description="Must match the password"
        {...register("confirmPassword")}
      />
      <Button
        isLoading={isLoading}
        color="primary"
        size="lg"
        type="submit"
      >
        Register
      </Button>
    </form>
  );
}
