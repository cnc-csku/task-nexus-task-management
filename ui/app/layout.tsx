import type { Metadata } from "next";
import "./globals.css";
import { Montserrat, Noto_Sans_Thai } from "next/font/google";
import HeroUIProvider from "@/lib/hero-ui/HeroUIProvider";
import QueryClientProvider from "@/lib/react-query/QueryClientProvider";
import NextAuthProvider from "@/lib/next-auth/NextAuthProvider";
import { getServerSession } from "next-auth";
import { Toaster } from "sonner";
import { authOptions } from "@/lib/next-auth/auth";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const notoSansThai = Noto_Sans_Thai({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  preload: true,
  style: ["normal"],
  subsets: ["latin", "latin-ext", "thai"],
});

const monserat = Montserrat({
  weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
  preload: true,
  style: ["normal"],
  subsets: ["latin", "latin-ext"],
});

export const metadata: Metadata = {
  title: "TaskNexus",
  description: "An agile task management for everyone.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession(authOptions);

  return (
    <html lang="en">
      <body className={monserat.className + " " + notoSansThai.className}>
        <HeroUIProvider>
          <QueryClientProvider>
            <NextAuthProvider session={session}>
              <ReactQueryDevtools initialIsOpen={false} />
              <Toaster
                position="top-right"
                closeButton
                richColors
              />
              <div className="bg-white">{children}</div>
            </NextAuthProvider>
          </QueryClientProvider>
        </HeroUIProvider>
      </body>
    </html>
  );
}
