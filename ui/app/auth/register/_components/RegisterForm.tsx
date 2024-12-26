import { Input } from "@nextui-org/input";
import { Button } from "@nextui-org/button";

export default function RegisterForm() {
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
      <Input
        label="Confirm Password"
        type="password"
      />
      <Input label="Name" />
      <Button color="primary">Register</Button>
    </form>
  );
}
