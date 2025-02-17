"use client";

import BackButton from "@/components/ui/BackButton";
import Header from "@/components/ui/Header";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Chip } from "@heroui/chip";
import { Avatar } from "@heroui/avatar";
import { Input } from "@heroui/input";
import { Select, SelectItem } from "@heroui/select";

export default function TaskDetailPage() {
  return (
    <div>
      <BackButton href="/tasks" />
      <Header>
        <Chip
          color="primary"
          variant="flat"
        >
          SR-1
        </Chip>{" "}
        Create a new design
      </Header>
      <div className="flex flex-col gap-3">
        <Card className="shadow-none border-1 border-gray-300">
          <CardBody className="flex gap-2">
            <div>
              <Chip className="mb-2">Dev</Chip>
              <div className="flex items-center gap-5 mb-2">
                <div className="flex gap-2 items-center">
                  <Avatar
                    src="https://i.pravatar.cc/150?img=1"
                    size="sm"
                  />
                  <div>Jane Doe</div>
                </div>
                <div>
                  <Input
                    size="sm"
                    placeholder="Point"
                    endContent={<>pts</>}
                    value={"5"}
                  />
                </div>
              </div>
            </div>
            <div>
              <Chip className="mb-2">QA</Chip>
              <div className="flex items-center gap-5 mb-2">
                <div className="flex gap-2 items-center">
                  <Avatar
                    src="https://i.pravatar.cc/150?img=2"
                    size="sm"
                  />
                  <div>Janny Doe</div>
                </div>
                <div>
                  <Input
                    size="sm"
                    placeholder="Point"
                    endContent={<>pts</>}
                    value={"3"}
                  />
                </div>
              </div>
            </div>
            <div className="flex gap-2 items-center">
              <Chip>Priority</Chip>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                <p className="text-sm">High</p>
              </div>
            </div>
          </CardBody>
        </Card>
        <Card className="shadow-none border-1 border-gray-300">
          <CardHeader>Description</CardHeader>
          <CardBody>
              <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla accumsan, metus ultrices
                eleifend gravida, nulla nunc varius lectus, nec rutrum justo nibh eu lectus. Ut vulputate
                semper dui. Fusce erat odio, sollicitudin vel erat vel, interdum mattis neque.
              </p>
          </CardBody>
        </Card>
        <Card className="shadow-none border-1 border-gray-300">
          <CardHeader>Comment</CardHeader>
          <CardBody>
            <div className="flex gap-3 items-center">
              <Avatar
                src="https://avatars.githubusercontent.com/u/86820985?v=4"
                size="sm"
              />
              <Input
                placeholder="Write a comment..."
              />
            </div>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}
