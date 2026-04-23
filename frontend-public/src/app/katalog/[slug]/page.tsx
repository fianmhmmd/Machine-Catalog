import { Metadata } from "next";
import { notFound } from "next/navigation";
import { getProductBySlug, getRelatedProducts, getProducts } from "@/lib/api";
import ProductDetail from "@/components/layout/ProductDetail";
import { formatImageUrl } from "@/lib/utils";

// ISR: Revalidate every 5 minutes
export const revalidate = 300;

interface PageProps {
  params: {
    slug: string;
  };
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const product = await getProductBySlug(params.slug);
  
  if (!product) {
    return {
      title: "Produk Tidak Ditemukan",
    };
  }

  const primaryImage = product.images?.find(img => img.is_primary) || product.images?.[0];
  const imageUrl = formatImageUrl(primaryImage?.image_url);

  return {
    title: product.name,
    description: product.description.substring(0, 160),
    openGraph: {
      title: product.name,
      description: product.description.substring(0, 160),
      images: [
        {
          url: imageUrl,
          width: 800,
          height: 800,
          alt: product.name,
        },
      ],
    },
  };
}

// Generate static paths for popular products
export async function generateStaticParams() {
  const productsResponse = await getProducts({ limit: 100 });
  return productsResponse.data.map((product) => ({
    slug: product.slug,
  }));
}

export default async function ProductPage({ params }: PageProps) {
  const [product, relatedProducts] = await Promise.all([
    getProductBySlug(params.slug),
    getRelatedProducts(params.slug),
  ]);

  if (!product) {
    notFound();
  }

  return <ProductDetail product={product} relatedProducts={relatedProducts} />;
}
