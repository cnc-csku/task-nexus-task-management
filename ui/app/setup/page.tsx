"use client";
import { Link } from "@nextui-org/link";
import { useState } from "react";
import { BsArrowRightCircle } from "react-icons/bs";
import { motion } from "framer-motion";

export default function SetupHomePage() {
  return (
    <div>
      <motion.div
        className="flex flex-col gap-5"
        initial={{ opacity: 0, x: -100 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      >
        <h2 className="text-5xl">Welcome to</h2>
        <h1 className="text-8xl font-semibold">TaskNexus</h1>
        <h3 className="text-primary-500 text-lg">An agile task management tool for everyone.</h3>
      </motion.div>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.5 }}
      >
        <Link
          href="/setup/register"
          className="mt-10 text-2xl lets-start-btn"
        >
          Let's start <BsArrowRightCircle className="transition-all inline-block text-3xl ml-2 lets-start-icon" />
        </Link>
      </motion.div>
    </div>
  );
}
