import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Link } from "@nextui-org/link";
import LoginForm from "./RegisterForm";
import TextDivider from "@/components/ui/TextDivider";
import SocialLoginForm from "@/app/auth/_components/SocialLoginForm";
import { twMerge } from "tailwind-merge";

export interface RegisterBoxProps {
  className?: string;
}

export default function RegisterBox({ className }: RegisterBoxProps) {
  return (
    <Card className={twMerge("px-3 py-5", className)}>
      <CardHeader>
        <h1 className="text-2xl font-bold">Register</h1>
      </CardHeader>
      <CardBody>
        <LoginForm />
        <TextDivider
          text="Or"
          className="my-4"
        />
        <SocialLoginForm />
        <div className="mt-5">
          Already have an account? <Link href="/auth/login">Login</Link>
        </div>
      </CardBody>
    </Card>
  );
}
