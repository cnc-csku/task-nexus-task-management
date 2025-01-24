export default function SetupLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <div className="flex flex-col items-center justify-center min-h-screen">{children}</div>;
}
