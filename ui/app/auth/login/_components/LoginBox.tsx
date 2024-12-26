import { Card, CardBody, CardHeader } from "@nextui-org/card";
import LoginForm from "./LoginForm";
import TextDivider from "@/components/ui/TextDivider";
import SocialLoginForm from "@/app/auth/_components/SocialLoginForm";
import { twMerge } from "tailwind-merge";
import { Link } from "@nextui-org/link";

export interface LoginBoxProps {
  className?: string;
}

export default function LoginBox({ className }: LoginBoxProps) {
  return (
    <Card className={twMerge("px-3 py-5", className)}>
      <CardHeader>
        <h1 className="text-2xl font-bold">Login</h1>
      </CardHeader>
      <CardBody>
        <LoginForm />
        <TextDivider
          text="Or"
          className="my-4"
        />
        <SocialLoginForm />
        <div className="mt-5">
          No account? <Link href="/auth/register">Register</Link>
        </div>
      </CardBody>
    </Card>
  );
}
