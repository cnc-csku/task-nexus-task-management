import MainContainer from "@/components/templates/MainContainer";
import SideBar from "@/components/ui/SideBar";

export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="grid grid-cols-12">
      <div className="col-span-2">
        <SideBar />
      </div>
      <div className="col-span-10">
        <MainContainer>{children}</MainContainer>
      </div>
    </div>
  );
}
