"use client";
import { Navbar as NextNav, NavbarBrand, NavbarContent, NavbarItem, NavbarMenuToggle, NavbarMenu, NavbarMenuItem } from "@nextui-org/navbar";
import { Dropdown, DropdownItem, DropdownMenu, DropdownTrigger } from "@nextui-org/dropdown";
import { Avatar } from "@nextui-org/avatar";
import { User } from "@nextui-org/user";
import { IoMdNotificationsOutline } from "react-icons/io";
import { Button } from "@nextui-org/button";
import { Badge } from "@nextui-org/badge";

export default function Navbar() {
  return (
    <NextNav
      maxWidth="full"
      className="bg-white h-10 shadow-sm"
    >
      <NavbarBrand>TaskNexus</NavbarBrand>
      <NavbarContent justify="end">
        <NavbarItem>
          <Dropdown placement="bottom-end">
            <DropdownTrigger>
              <Button
                variant="light"
                size="sm"
                isIconOnly
              >
                <Badge
                  color="danger"
                  content="3"
                  shape="circle"
                  placement="bottom-right"
                  size="sm"
                >
                  <IoMdNotificationsOutline className="text-xl cursor-pointer" />
                </Badge>
              </Button>
            </DropdownTrigger>
            <DropdownMenu
              aria-label="Notifications"
              variant="flat"
            >
              <DropdownItem key="notification1">Notification 1</DropdownItem>
              <DropdownItem key="notification2">Notification 2</DropdownItem>
              <DropdownItem key="notification3">Notification 3</DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </NavbarItem>
        <NavbarItem>
          <Dropdown placement="bottom-end">
            <DropdownTrigger>
              <Avatar
                isBordered
                as="button"
                className="transition-transform"
                size="sm"
                src="https://avatars.githubusercontent.com/u/86820985?v=4"
              />
            </DropdownTrigger>
            <DropdownMenu
              aria-label="Profile Actions"
              variant="flat"
            >
              <DropdownItem
                key="profile"
                isReadOnly
              >
                <User
                  avatarProps={{
                    size: "sm",
                    src: "https://avatars.githubusercontent.com/u/86820985?v=4",
                  }}
                  classNames={{
                    name: "text-default-600",
                    description: "text-default-500",
                  }}
                  name="Tanaroeg O-Charoen"
                  description="tanaroeg.o@ku.th"
                />
              </DropdownItem>
              <DropdownItem key="editProfile">Edit Profile</DropdownItem>
              <DropdownItem key="settings">My Settings</DropdownItem>
              <DropdownItem
                key="logout"
                color="danger"
              >
                Log Out
              </DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </NavbarItem>
      </NavbarContent>
    </NextNav>
  );
}
