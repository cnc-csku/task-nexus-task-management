import MainContainer from "@/components/templates/MainContainer";
import Navbar from "@/components/ui/Navbar";
import SideBar from "@/components/ui/SideBar";

export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="bg-white h-screen">
      <Navbar />
      {children}
    </div>
  );
}
