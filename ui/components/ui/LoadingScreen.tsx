import { Spinner } from "@nextui-org/spinner";

export default function LoadingScreen() {
  return (
    <div className="fixed top-0 left-0 w-full h-full flex justify-center items-center">
      <Spinner />
    </div>
  );
}
