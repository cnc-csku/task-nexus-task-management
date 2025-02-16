import { Card, CardBody, CardHeader } from "@heroui/card";
import { Avatar } from "@heroui/avatar";
import { Badge } from "@heroui/badge";
import { Chip } from "@heroui/chip";

export default function TaskCard() {
  return (
    <Card className="shadow-none border-1 border-gray-300">
      <CardHeader className="flex items-center gap-2">
        <Chip
          color="primary"
          size="sm"
        >
          Task
        </Chip>
         Create a new design
      </CardHeader>
      <CardBody>
        <p className="text-sm text-gray-500">Task description goes here...</p>

        <div className="flex items-center gap-2 justify-between mt-2">
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
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-red-500 rounded-full"></div>
            <p className="text-sm">High</p>
          </div>
        </div>
      </CardBody>
    </Card>
  );
}
