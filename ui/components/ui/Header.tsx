export default function Header({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <h1 className="text-xl font-normal mb-3">{children}</h1>;
}
