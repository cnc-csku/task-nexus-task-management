"use client";

import { Table, TableBody, TableCell, TableColumn, TableHeader, TableRow } from "@nextui-org/table";
import { Link } from "@nextui-org/link";
import React from "react";
import { Avatar } from "@nextui-org/avatar";
import { Chip } from "@nextui-org/chip";
import { Badge } from "@nextui-org/badge";
import { useRouter } from "next/navigation";

export default function TasksTable() {
  const router = useRouter();
  return (
    <Table
      onRowAction={(key) => router.push("/tasks/1")}
      selectionMode="single"
    >
      <TableHeader>
        <TableColumn>ID</TableColumn>
        <TableColumn>Type</TableColumn>
        <TableColumn>Name</TableColumn>
        <TableColumn>Sprint</TableColumn>
        <TableColumn>Assignee</TableColumn>
        <TableColumn className="text-center">Status</TableColumn>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell>
            SR-1
          </TableCell>
          <TableCell>
            <Chip
              color="primary"
              variant="flat"
            >
              Task
            </Chip>
          </TableCell>
          <TableCell>Create a new design</TableCell>
          <TableCell>
            <Link href="#">SR-SP-1</Link>
          </TableCell>
          <TableCell>
            <div className="flex items-center gap-2 mb-2">
              <Badge
                color="primary"
                content="2"
                placement="bottom-right"
              >
                <Avatar
                  src="https://i.pravatar.cc/150?img=1"
                  size="sm"
                />
              </Badge>
              <Badge
                color="primary"
                content="3"
                placement="bottom-right"
              >
                <Avatar
                  src="https://i.pravatar.cc/150?img=2"
                  size="sm"
                />
              </Badge>
            </div>
          </TableCell>
          <TableCell className="text-center">
            <Chip
              color="default"
              variant="flat"
            >
              Todo
            </Chip>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
}
