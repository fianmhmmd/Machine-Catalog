"use client";

import { useState, useEffect, useCallback } from "react";
import { Search, Filter, ChevronLeft, ChevronRight, Loader2 } from "lucide-react";
import { Product, Category, PaginatedResponse } from "@/types";
import ProductCard from "@/components/ui/ProductCard";
import { getProducts } from "@/lib/api";
import { useRouter, useSearchParams } from "next/navigation";
import { cn } from "@/lib/utils";

interface CatalogGridProps {
  initialData: PaginatedResponse<Product>;
  categories: Category[];
}

export default function CatalogGrid({ initialData, categories }: CatalogGridProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  
  const [products, setProducts] = useState<Product[]>(initialData.data);
  const [meta, setMeta] = useState(initialData.meta);
  const [isLoading, setIsLoading] = useState(false);
  const [search, setSearch] = useState(searchParams.get("search") || "");
  const [category, setCategory] = useState(searchParams.get("category") || "");
  const [page, setPage] = useState(Number(searchParams.get("page")) || 1);

  const fetchProducts = useCallback(async () => {
    setIsLoading(true);
    const data = await getProducts({
      search,
      category,
      page,
      limit: 12,
    });
    setProducts(data.data);
    setMeta(data.meta);
    setIsLoading(false);
  }, [search, category, page]);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      setPage(1);
      updateUrl({ search, category, page: 1 });
      fetchProducts();
    }, 500);
    return () => clearTimeout(timer);
  }, [search]);

  // Filter and Page changes
  useEffect(() => {
    updateUrl({ search, category, page });
    fetchProducts();
  }, [category, page]);

  const updateUrl = (params: { search?: string; category?: string; page?: number }) => {
    const newParams = new URLSearchParams();
    if (params.search) newParams.set("search", params.search);
    if (params.category) newParams.set("category", params.category);
    if (params.page && params.page > 1) newParams.set("page", params.page.toString());
    
    router.push(`/katalog?${newParams.toString()}`, { scroll: false });
  };

  return (
    <div className="flex flex-col gap-8">
      {/* Search & Filter Bar */}
      <div className="flex flex-col md:flex-row gap-4">
        <div className="relative flex-grow">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
          <input
            type="text"
            placeholder="Cari mesin manufaktur..."
            className="w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="relative min-w-[200px]">
          <Filter className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
          <select
            className="w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none appearance-none bg-white transition-all"
            value={category}
            onChange={(e) => {
              setCategory(e.target.value);
              setPage(1);
            }}
          >
            <option value="">Semua Kategori</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.slug}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Grid */}
      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-24 gap-4">
          <Loader2 className="h-12 w-12 text-blue-600 animate-spin" />
          <p className="text-gray-500 font-medium">Memuat produk...</p>
        </div>
      ) : products.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      ) : (
        <div className="text-center py-24 bg-white rounded-xl border border-dashed border-gray-300">
          <div className="mx-auto h-12 w-12 text-gray-400 mb-4">
            <Search className="h-12 w-12" />
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Produk tidak ditemukan</h3>
          <p className="text-gray-500 mt-1">Coba ubah kata kunci atau filter kategori Anda.</p>
          <button
            onClick={() => {
              setSearch("");
              setCategory("");
            }}
            className="mt-4 text-blue-600 font-medium hover:underline"
          >
            Reset Semua Filter
          </button>
        </div>
      )}

      {/* Pagination */}
      {meta.total > meta.limit && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1 || isLoading}
            className="p-2 rounded-md border border-gray-200 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            <ChevronLeft className="h-5 w-5" />
          </button>
          <div className="flex gap-1">
            {Array.from({ length: Math.ceil(meta.total / meta.limit) }).map((_, i) => (
              <button
                key={i}
                onClick={() => setPage(i + 1)}
                className={cn(
                  "h-10 w-10 rounded-md text-sm font-medium transition-colors",
                  page === i + 1
                    ? "bg-blue-600 text-white"
                    : "hover:bg-gray-100 text-gray-700"
                )}
              >
                {i + 1}
              </button>
            ))}
          </div>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={page >= Math.ceil(meta.total / meta.limit) || isLoading}
            className="p-2 rounded-md border border-gray-200 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            <ChevronRight className="h-5 w-5" />
          </button>
        </div>
      )}
    </div>
  );
}
