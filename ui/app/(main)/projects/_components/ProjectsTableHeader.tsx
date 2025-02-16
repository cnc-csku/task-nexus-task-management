"use client";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { IoMdAdd } from "react-icons/io";
import { MdSearch } from "react-icons/md";

export default function ProjectsTableHeader() {
  return (
    <div className="flex justify-between items-center">
      <div>
        <Input
          label="Search"
          size="sm"
          className="md:w-80"
          startContent={<MdSearch />}
        />
      </div>
      <div>
        <Button
          startContent={<IoMdAdd />}
          color="primary"
        >
          Add Project
        </Button>
      </div>
    </div>
  );
}
