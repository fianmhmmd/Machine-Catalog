"use client";

import { useEffect } from "react";
import { Product } from "@/types";
import { trackView } from "@/lib/api";
import Gallery from "@/components/ui/Gallery";
import WhatsAppButton from "@/components/ui/WhatsAppButton";
import InquiryForm from "@/components/ui/InquiryForm";
import ProductCard from "@/components/ui/ProductCard";
import { Toaster } from "react-hot-toast";

interface ProductDetailProps {
  product: Product;
  relatedProducts: Product[];
}

export default function ProductDetail({ product, relatedProducts }: ProductDetailProps) {
  useEffect(() => {
    trackView(product.id);
  }, [product.id]);

  return (
    <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <Toaster position="top-right" />
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        {/* Left: Gallery */}
        <Gallery images={product.images} />

        {/* Right: Info */}
        <div className="flex flex-col gap-8">
          <div>
            <p className="text-sm font-medium text-blue-600 uppercase tracking-widest mb-2">
              {product.category.name}
            </p>
            <h1 className="text-3xl font-extrabold text-gray-900 sm:text-4xl">
              {product.name}
            </h1>
          </div>

          <div className="prose prose-blue max-w-none text-gray-600">
            <h3 className="text-lg font-bold text-gray-900 mb-2">Deskripsi Produk</h3>
            <p className="whitespace-pre-line">{product.description}</p>
          </div>

          {product.specifications && Object.keys(product.specifications).length > 0 && (
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <div className="bg-gray-50 px-6 py-4 border-b border-gray-200">
                <h3 className="font-bold text-gray-900">Spesifikasi Teknik</h3>
              </div>
              <table className="w-full text-sm">
                <tbody className="divide-y divide-gray-200">
                  {Object.entries(product.specifications).map(([key, value]) => (
                    <tr key={key}>
                      <td className="px-6 py-4 font-medium text-gray-500 w-1/3 bg-gray-50/50 capitalize">
                        {key.replace(/_/g, " ")}
                      </td>
                      <td className="px-6 py-4 text-gray-900">
                        {String(value)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <div className="flex flex-col gap-4 mt-4">
            <WhatsAppButton 
              productId={product.id}
              phoneNumber={product.contact_phone}
              productName={product.name}
            />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-12 mt-16">
        <div className="lg:col-span-2">
          {/* Related Products */}
          <section>
            <h2 className="text-2xl font-bold text-gray-900 mb-8">Produk Terkait</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
              {relatedProducts.map((p) => (
                <ProductCard key={p.id} product={p} />
              ))}
              {relatedProducts.length === 0 && (
                <p className="text-gray-500 italic">Tidak ada produk terkait.</p>
              )}
            </div>
          </section>
        </div>

        <div className="lg:col-span-1">
          {/* Inquiry Form */}
          <InquiryForm 
            productId={product.id}
            productName={product.name}
          />
        </div>
      </div>
    </div>
  );
}
