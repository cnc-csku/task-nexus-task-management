"use client";

import { MdAnalytics, MdSettings, MdTimeline, MdDashboard } from "react-icons/md";
import { GrPowerCycle } from "react-icons/gr";
import { GoTasklist } from "react-icons/go";
import SideBarItem from "./SideBarItem";
import { Button } from "@heroui/button";
import { useParams } from "next/navigation";
import { IoMdArrowRoundBack } from "react-icons/io";
import Link from "next/link";

export interface SideBarProps {
  isOpen: boolean;
}
export default function SideBar({ isOpen }: SideBarProps) {
  const { projId } = useParams<{ projId: string }>();

  return (
    <>
      {isOpen && (
        <div className="w-full bg-white border-r-1 border-gray-200 h-screen px-5 py-6 flex flex-col justify-between">
          <div>
            {/* Add back to all project button */}
            <Button
              variant="light"
              isIconOnly
              as={Link}
              href="/workspaces"
            >
              <IoMdArrowRoundBack className="text-lg" />
            </Button>
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
      )}
    </>
  );
}
