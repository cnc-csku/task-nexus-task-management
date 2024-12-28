export default function MainContainer({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <div className="px-6 py-5">{children}</div>;
}
