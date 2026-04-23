"use client";

import { useState } from "react";
import Image from "next/image";
import { ProductImage } from "@/types";
import { formatImageUrl } from "@/lib/utils";
import { cn } from "@/lib/utils";

interface GalleryProps {
  images: ProductImage[];
}

export default function Gallery({ images }: GalleryProps) {
  const sortedImages = [...images].sort((a, b) => a.sort_order - b.sort_order);
  const [activeImage, setActiveImage] = useState(
    sortedImages.find((img) => img.is_primary) || sortedImages[0]
  );

  if (!sortedImages.length) {
    return (
      <div className="aspect-square bg-gray-100 rounded-lg flex items-center justify-center text-gray-400">
        No image available
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="relative aspect-square overflow-hidden rounded-xl bg-gray-100 border border-gray-200">
        <Image
          src={formatImageUrl(activeImage?.image_url)}
          alt="Product image"
          fill
          className="object-cover"
          sizes="(max-width: 768px) 100vw, 50vw"
          priority
        />
      </div>
      
      {sortedImages.length > 1 && (
        <div className="grid grid-cols-5 gap-2">
          {sortedImages.map((img) => (
            <button
              key={img.id}
              onClick={() => setActiveImage(img)}
              className={cn(
                "relative aspect-square overflow-hidden rounded-md border-2 transition-all",
                activeImage?.id === img.id ? "border-blue-600 shadow-sm" : "border-transparent hover:border-gray-300"
              )}
            >
              <Image
                src={formatImageUrl(img.image_url)}
                alt="Thumbnail"
                fill
                className="object-cover"
                sizes="100px"
              />
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
