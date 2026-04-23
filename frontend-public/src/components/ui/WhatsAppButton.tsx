"use client";

import { MessageCircle } from "lucide-react";
import { trackClick } from "@/lib/api";

interface WhatsAppButtonProps {
  productId: string;
  phoneNumber: string;
  productName: string;
}

export default function WhatsAppButton({ productId, phoneNumber, productName }: WhatsAppButtonProps) {
  const handleClick = async () => {
    // Track click to API
    await trackClick(productId);

    // Format phone number (remove leading 0 or +)
    let formattedPhone = phoneNumber.replace(/[^0-9]/g, "");
    if (formattedPhone.startsWith("0")) {
      formattedPhone = "62" + formattedPhone.slice(1);
    } else if (!formattedPhone.startsWith("62")) {
      formattedPhone = "62" + formattedPhone;
    }

    const message = encodeURIComponent(
      `Halo, saya ingin bertanya tentang produk *${productName}* yang ada di katalog.`
    );
    const waUrl = `https://wa.me/${formattedPhone}?text=${message}`;
    
    window.open(waUrl, "_blank");
  };

  return (
    <button
      onClick={handleClick}
      className="flex items-center justify-center gap-2 bg-green-500 hover:bg-green-600 text-white font-bold py-4 px-6 rounded-xl transition-all shadow-md hover:shadow-lg w-full"
    >
      <MessageCircle className="h-6 w-6" />
      Hubungi via WhatsApp
    </button>
  );
}
