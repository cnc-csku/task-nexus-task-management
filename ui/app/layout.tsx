import type { Metadata } from "next";
import "./globals.css";
import { Montserrat, Noto_Sans_Thai } from "next/font/google";
import ToastContainer from "@/lib/react-toastify/ToastContainer";
import NextUIProvider from "@/lib/next-ui/NextUIProvider";
import QueryClientProvider from "@/lib/react-query/QueryClientProvider";

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

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="th">
      <body className={monserat.className}>
        <NextUIProvider>
          <QueryClientProvider>
            <ToastContainer />
            {children}
          </QueryClientProvider>
        </NextUIProvider>
      </body>
    </html>
  );
}
