import type { Metadata } from "next";
import { Plus_Jakarta_Sans } from "next/font/google";

import { StoreProvider } from "@/store/provider";

import "./globals.css";

const plusJakarta = Plus_Jakarta_Sans({
  subsets: ["latin"],
  variable: "--font-plus-jakarta",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "FlashShare",
    template: "%s · FlashShare",
  },
  description: "Temporary file sharing — send big files in a flash.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${plusJakarta.variable} antialiased`}>
        <StoreProvider>{children}</StoreProvider>
      </body>
    </html>
  );
}
