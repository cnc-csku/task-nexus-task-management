import { Input } from "@nextui-org/input";
import { Button } from "@nextui-org/button";

export default function LoginForm() {
  return (
    <form className="flex flex-col space-y-4">
      <Input
        label="Email"
        type="email"
      />
      <Input
        label="Password"
        type="password"
      />
      <Button color="primary">Login</Button>
    </form>
  );
}
