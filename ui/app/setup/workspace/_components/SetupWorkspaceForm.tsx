"use client";

import { CreateWorkspaceFormSchema, CreateWorkspaceFormType } from "@/interfaces/Workspace";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { useForm } from "react-hook-form";

export interface SetupWorkspaceFormProps {
  onCreate: (data: CreateWorkspaceFormType) => void;
  isLoading: boolean;
}

export default function SetupWorkspaceForm({ isLoading, onCreate }: SetupWorkspaceFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateWorkspaceFormType>({
    resolver: zodResolver(CreateWorkspaceFormSchema),
  });

  return (
    <form
      onSubmit={handleSubmit(onCreate)}
      className="flex flex-col gap-4"
    >
      <Input
        label="Name"
        type="text"
        isRequired
        isInvalid={!!errors.name}
        errorMessage={errors.name?.message}
        {...register("name")}
      />
      <Button
        color="primary"
        size="lg"
        type="submit"
        isLoading={isLoading}
      >
        Create
      </Button>
    </form>
  );
}
