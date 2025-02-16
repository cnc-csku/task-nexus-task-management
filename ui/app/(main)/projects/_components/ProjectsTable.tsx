"use client";

import { Table, TableBody, TableCell, TableColumn, TableHeader, TableRow } from "@heroui/table";
import ProjectsTableHeader from "./ProjectsTableHeader";
import { Button } from "@heroui/button";
import { Dropdown, DropdownTrigger, DropdownItem, DropdownMenu } from "@heroui/dropdown";
import { IoMdMore } from "react-icons/io";
import { Link } from "@heroui/link";
import { Avatar } from "@heroui/avatar";
import { Chip } from "@heroui/chip";
import { useRouter } from "next/navigation";
import { Key } from "react";

export default function ProjectsTable() {
  const router = useRouter();

  const rowAction = (key: Key) => {
    router.push("/projects/1/board");
  };

  return (
    <Table
      topContent={<ProjectsTableHeader />}
      onRowAction={rowAction}
      selectionMode="single"
      aria-label="Projects Table"
    >
      <TableHeader>
        <TableColumn>Name</TableColumn>
        <TableColumn>Prefix</TableColumn>
        <TableColumn>Owner</TableColumn>
        <TableColumn className="text-center">Status</TableColumn>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell>Senior Project</TableCell>
          <TableCell>SP</TableCell>
          <TableCell>
            <Link className="text-sm">
              <Avatar
                isBordered
                className="transition-transform mr-1"
                size="sm"
                src="https://avatars.githubusercontent.com/u/86820985?v=4"
              />{" "}
              Tanaroeg O-Charoen
            </Link>
          </TableCell>
          <TableCell className="text-center">
            <Chip
              color="success"
              variant="flat"
            >
              Active
            </Chip>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
}
