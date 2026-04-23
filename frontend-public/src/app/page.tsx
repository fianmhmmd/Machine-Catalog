import Link from "next/link";
import { ArrowRight, Settings, Shield, Zap } from "lucide-react";
import { getLatestProducts, getCategories } from "@/lib/api";
import ProductCard from "@/components/ui/ProductCard";

// ISR: Revalidate every 5 minutes
export const revalidate = 300;

export default async function Home() {
  const [latestProducts, categories] = await Promise.all([
    getLatestProducts(6),
    getCategories(),
  ]);

  return (
    <div className="flex flex-col gap-16 pb-16">
      {/* Hero Section */}
      <section className="relative bg-blue-900 py-24 sm:py-32 overflow-hidden">
        <div className="absolute inset-0 opacity-20">
          <div className="absolute inset-0 bg-gradient-to-r from-blue-900 to-transparent" />
        </div>
        <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center sm:text-left">
          <div className="max-w-2xl">
            <h1 className="text-4xl font-extrabold tracking-tight text-white sm:text-6xl">
              Solusi Mesin Manufaktur <span className="text-blue-400">Terbaik</span> Untuk Bisnis Anda
            </h1>
            <p className="mt-6 text-lg text-blue-100 leading-8">
              Tingkatkan produktivitas industri Anda dengan pilihan mesin manufaktur berkualitas tinggi. 
              Dukungan teknis terpercaya dan pengiriman ke seluruh Indonesia.
            </p>
            <div className="mt-10 flex flex-wrap gap-4 justify-center sm:justify-start">
              <Link
                href="/katalog"
                className="rounded-md bg-blue-500 px-6 py-3 text-lg font-semibold text-white shadow-sm hover:bg-blue-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-400 transition-all"
              >
                Lihat Katalog
              </Link>
              <Link
                href="#kategori"
                className="rounded-md bg-white/10 px-6 py-3 text-lg font-semibold text-white hover:bg-white/20 transition-all"
              >
                Telusuri Kategori
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="flex flex-col items-center p-6 bg-white rounded-xl shadow-sm border border-gray-100 text-center">
            <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
              <Shield className="h-6 w-6 text-blue-600" />
            </div>
            <h3 className="text-xl font-bold mb-2">Kualitas Terjamin</h3>
            <p className="text-gray-500">Seluruh mesin kami telah melewati uji kualitas yang ketat sebelum sampai ke tangan Anda.</p>
          </div>
          <div className="flex flex-col items-center p-6 bg-white rounded-xl shadow-sm border border-gray-100 text-center">
            <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
              <Zap className="h-6 w-6 text-blue-600" />
            </div>
            <h3 className="text-xl font-bold mb-2">Performa Tinggi</h3>
            <p className="text-gray-500">Didesain untuk efisiensi maksimal guna mendukung percepatan target produksi industri Anda.</p>
          </div>
          <div className="flex flex-col items-center p-6 bg-white rounded-xl shadow-sm border border-gray-100 text-center">
            <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
              <Settings className="h-6 w-6 text-blue-600" />
            </div>
            <h3 className="text-xl font-bold mb-2">Dukungan Teknis</h3>
            <p className="text-gray-500">Tim ahli kami siap membantu pemasangan dan pemeliharaan mesin secara berkala.</p>
          </div>
        </div>
      </section>

      {/* Categories Grid */}
      <section id="kategori" className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 scroll-mt-24">
        <div className="flex items-center justify-between mb-8">
          <h2 className="text-3xl font-bold text-gray-900">Kategori Produk</h2>
          <Link href="/katalog" className="text-blue-600 font-medium flex items-center gap-1 hover:underline">
            Lihat Semua <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
          {categories.map((cat) => (
            <Link
              key={cat.id}
              href={`/katalog?category=${cat.slug}`}
              className="group relative h-40 flex flex-col items-center justify-center bg-white rounded-lg border border-gray-200 overflow-hidden hover:border-blue-500 transition-all text-center p-4"
            >
              <div className="absolute inset-0 bg-blue-500 opacity-0 group-hover:opacity-5 transition-opacity" />
              <h3 className="text-lg font-semibold text-gray-800 group-hover:text-blue-600 transition-colors">
                {cat.name}
              </h3>
              <p className="text-xs text-gray-500 mt-2">Lihat Produk</p>
            </Link>
          ))}
        </div>
      </section>

      {/* Latest Products */}
      <section className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between mb-8">
          <h2 className="text-3xl font-bold text-gray-900">Produk Terbaru</h2>
          <Link href="/katalog" className="text-blue-600 font-medium flex items-center gap-1 hover:underline">
            Semua Produk <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {latestProducts.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
          {latestProducts.length === 0 && (
            <div className="col-span-full py-12 text-center text-gray-500 bg-gray-100 rounded-lg">
              Belum ada produk yang ditambahkan.
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
