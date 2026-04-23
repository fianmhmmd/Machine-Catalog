import { Metadata } from "next";
import { getProducts, getCategories } from "@/lib/api";
import CatalogGrid from "@/components/layout/CatalogGrid";

export const metadata: Metadata = {
  title: "Katalog Produk",
  description: "Telusuri berbagai macam mesin manufaktur berkualitas tinggi mulai dari mesin CNC, Bubut, hingga Welding.",
};

export default async function CatalogPage({
  searchParams,
}: {
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  const page = Number(searchParams.page) || 1;
  const category = searchParams.category as string | undefined;
  const search = searchParams.search as string | undefined;

  const [initialData, categories] = await Promise.all([
    getProducts({ page, category, search, limit: 12 }),
    getCategories(),
  ]);

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <div className="flex flex-col gap-4 mb-12">
        <h1 className="text-4xl font-bold text-gray-900 tracking-tight">
          Katalog Mesin Industri
        </h1>
        <p className="text-lg text-gray-500 max-w-2xl">
          Temukan solusi mesin manufaktur terbaik untuk mengoptimalkan operasional bisnis Anda. 
          Gunakan filter untuk menemukan spesifikasi yang tepat.
        </p>
      </div>

      <CatalogGrid 
        initialData={initialData} 
        categories={categories} 
      />
    </div>
  );
}
