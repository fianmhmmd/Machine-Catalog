import Link from "next/link";
import Image from "next/image";
import { Product } from "@/types";
import { formatImageUrl } from "@/lib/utils";

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const primaryImage = product.images?.find(img => img.is_primary) || product.images?.[0];
  const imageUrl = formatImageUrl(primaryImage?.image_url);

  return (
    <Link 
      href={`/katalog/${product.slug}`}
      className="group bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow"
    >
      <div className="relative aspect-square overflow-hidden bg-gray-100">
        <Image
          src={imageUrl}
          alt={product.name}
          fill
          className="object-cover group-hover:scale-105 transition-transform duration-300"
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
        />
      </div>
      <div className="p-4">
        <p className="text-xs font-medium text-blue-600 mb-1 uppercase tracking-wider">
          {product.category?.name}
        </p>
        <h3 className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-2">
          {product.name}
        </h3>
        <p className="mt-2 text-sm text-gray-500 line-clamp-2">
          {product.description}
        </p>
      </div>
    </Link>
  );
}
