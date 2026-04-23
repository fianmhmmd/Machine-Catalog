import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatImageUrl(url: string | undefined) {
  if (!url) return "/placeholder-image.jpg";
  if (url.startsWith("http")) return url;
  
  const minioUrl = process.env.NEXT_PUBLIC_MINIO_URL || "http://localhost:9000";
  return `${minioUrl}${url}`;
}
