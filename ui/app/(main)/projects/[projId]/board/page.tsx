"use client";
import Header from "@/components/ui/Header";
import TaskCard from "./_components/TaskCard";
import { Select, SelectItem } from "@nextui-org/select";
import { Button } from "@nextui-org/button";
import { useDisclosure } from "@nextui-org/modal";
import SprintGoalModal from "./_components/SprintGoalModal";

export default function BoardPage() {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  return (
    <div>
      <Header>Board</Header>
      <div className="flex mb-2 justify-between items-center">
        <div>Sprint: SR-SP-1</div>
        <div>
          <Button
            variant="flat"
            color="primary"
            onPress={onOpen}
          >
            Sprint Goal
          </Button>
        </div>
      </div>
      <div>
        <Select placeholder="Filter Assignee">
          <SelectItem value="1">Tanaroeg O-Charoen</SelectItem>
          <SelectItem value="2">John Doe</SelectItem>
          <SelectItem value="3">Jane Doe</SelectItem>
        </Select>
      </div>
      <div className="grid grid-cols-3 gap-6 mt-5">
        <div>
          <div className="bg-white rounded-md p-4">
            <h1 className="text-xl font-semibold">To Do</h1>
            <div className="mt-4 flex flex-col gap-2">
              <TaskCard />
              <TaskCard />
              <TaskCard />
            </div>
          </div>
        </div>
        <div>
          <div className="bg-white rounded-md p-4">
            <h1 className="text-xl font-semibold">In Progress</h1>
            <div className="mt-4 flex flex-col gap-2">
              <TaskCard />
              <TaskCard />
              <TaskCard />
            </div>
          </div>
        </div>
        <div>
          <div className="bg-white rounded-md p-4">
            <h1 className="text-xl font-semibold">Complete</h1>
            <div className="mt-4 flex flex-col gap-2">
              <TaskCard />
              <TaskCard />
              <TaskCard />
            </div>
          </div>
        </div>
      </div>
      <SprintGoalModal
        isOpen={isOpen}
        onOpenChange={onOpenChange}
      />
    </div>
  );
}
