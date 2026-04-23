import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: {
    default: "Machine Katalog - Katalog Mesin Manufaktur Terlengkap",
    template: "%s | Machine Katalog"
  },
  description: "Temukan berbagai mesin manufaktur berkualitas tinggi untuk kebutuhan industri Anda.",
  keywords: ["mesin manufaktur", "mesin industri", "katalog mesin", "jual mesin industri"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body className={`${inter.className} flex flex-col min-h-screen bg-gray-50`}>
        <Header />
        <main className="flex-grow">
          {children}
        </main>
        <Footer />
      </body>
    </html>
  );
}
