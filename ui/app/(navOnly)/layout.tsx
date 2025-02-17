import Navbar from "@/components/ui/Navbar";

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
