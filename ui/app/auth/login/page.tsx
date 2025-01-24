"use client";

import { motion } from "framer-motion";
import LoginForm from "./_components/LoginForm";
import { UserLoginFormType } from "@/interfaces/User";
import { useState } from "react";
import { signIn } from "next-auth/react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const router = useRouter();

  const onLogin = async (data: UserLoginFormType) => {
    setIsLoading(true);

    const response = await signIn("credentials", {
      email: data.email,
      password: data.password,
      redirect: false,
    });

    if (response?.ok) {
      router.push("/workspaces");
      return;
    } else {
      toast.error("Invalid email or password");
    }

    setIsLoading(false);
  };

  return (
    <div className="lg:w-[500px]">
      <motion.div
        className="flex flex-col gap-5 px-1"
        initial={{ opacity: 0, x: -100 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      >
        <h1 className="text-2xl mb-5">Login</h1>
      </motion.div>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.5 }}
      >
        <LoginForm
          onLogin={onLogin}
          isLoading={isLoading}
        />
      </motion.div>
    </div>
  );
}
