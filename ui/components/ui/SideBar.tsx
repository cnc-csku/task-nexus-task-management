"use client";

import { MdKeyboardArrowDown, MdAnalytics, MdSettings, MdTimeline, MdDashboard } from "react-icons/md";
import { Dropdown, DropdownItem, DropdownMenu, DropdownTrigger } from "@heroui/dropdown";
import { GrPowerCycle } from "react-icons/gr";
import { GoTasklist } from "react-icons/go";
import SideBarItem from "./SideBarItem";
import { Button } from "@heroui/button";
import { useParams } from "next/navigation";

export default function SideBar() {
  const { projId } = useParams<{ projId: string }>();

  return (
    <div className="w-full bg-white border-r-1 border-gray-200 h-[calc(100vh-36px)] px-5 py-6 flex flex-col justify-between">
      <div>
        <Dropdown>
          <DropdownTrigger>
            <Button
              variant="ghost"
              className="border border-gray-300 justify-between"
              fullWidth
              endContent={<MdKeyboardArrowDown />}
            >
              Senior Project
            </Button>
          </DropdownTrigger>
          <DropdownMenu>
            <DropdownItem key="pj1">Project 1</DropdownItem>
            <DropdownItem key="pj2">Project 2</DropdownItem>
            <DropdownItem key="pj3">Project 3</DropdownItem>
            <DropdownItem key="all">All Project</DropdownItem>
          </DropdownMenu>
        </Dropdown>
        <div className="mt-3 flex flex-col gap-2">
          <SideBarItem
            name="Board"
            href={`/projects/${projId}/board`}
            startIcon={<MdDashboard className="text-lg" />}
          />
          <SideBarItem
            name="Tasks"
            href={`/projects/${projId}/tasks`}
            startIcon={<GoTasklist className="text-lg" />}
          />
          <SideBarItem
            name="Sprint"
            href={`/projects/${projId}/sprints`}
            startIcon={<GrPowerCycle className="text-lg" />}
          />
          <SideBarItem
            name="Timeline"
            href={`/projects/${projId}/timeline`}
            startIcon={<MdTimeline className="text-lg" />}
          />
          <SideBarItem
            name="Report"
            href={`/projects/${projId}/report`}
            startIcon={<MdAnalytics className="text-lg" />}
          />
        </div>
      </div>
      <div>
        <SideBarItem
          name="Project Setting"
          href="/customers"
          startIcon={<MdSettings className="text-lg" />}
        />
      </div>
    </div>
  );
}
