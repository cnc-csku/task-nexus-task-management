import { Button, ButtonProps } from "@nextui-org/button";
import Image from "next/image";

export interface GoogleLoginButtonProps extends ButtonProps {}

export default function GoogleLoginButton({ ...props }: GoogleLoginButtonProps) {
  return (
    <Button
      className="flex items-center justify-between space-x-2 w-full"
      variant="flat"
      startContent={
        <Image
          src="https://static.cdnlogo.com/logos/g/23/goolge-icon.png"
          alt="Google"
          width={24}
          height={24}
        />
      }
      {...props}
    >
      <span>Continue with Google</span>
      <div></div>
    </Button>
  );
}
