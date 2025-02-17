"use client";

import { motion } from "framer-motion";
import { UserRegisterFormType } from "@/interfaces/User";
import RegisterForm from "@/components/auth/RegisterForm";

export default function RegisterPage() {
  const onRegister = async (data: UserRegisterFormType) => {
    console.log(data);
  };

  return (
    <div className="lg:w-[500px]">
      <motion.div
        className="flex flex-col gap-5 px-1"
        initial={{ opacity: 0, x: -100 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      >
        <h1 className="text-2xl mb-5">Register</h1>
      </motion.div>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.5 }}
      >
        <RegisterForm
          onRegister={onRegister}
          isLoading={false}
        />
      </motion.div>
    </div>
  );
}
