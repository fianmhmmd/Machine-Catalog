"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, Send } from "lucide-react";
import { submitInquiry } from "@/lib/api";
import { toast } from "react-hot-toast";

const inquirySchema = z.object({
  customer_name: z.string().min(2, "Nama minimal 2 karakter"),
  customer_email: z.string().email("Email tidak valid"),
  customer_phone: z.string().min(10, "Nomor telepon minimal 10 digit"),
  message: z.string().min(10, "Pesan minimal 10 karakter"),
});

type InquiryValues = z.infer<typeof inquirySchema>;

interface InquiryFormProps {
  productId: string;
  productName: string;
}

export default function InquiryForm({ productId, productName }: InquiryFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<InquiryValues>({
    resolver: zodResolver(inquirySchema),
    defaultValues: {
      message: `Halo, saya tertarik dengan produk ${productName}. Mohon info lebih lanjut.`,
    },
  });

  const onSubmit = async (values: InquiryValues) => {
    setIsSubmitting(true);
    try {
      await submitInquiry({
        product_id: productId,
        ...values,
      });
      toast.success("Pertanyaan Anda telah terkirim!");
      reset();
    } catch (error: any) {
      toast.error(error || "Gagal mengirim pertanyaan");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
      <h3 className="text-xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        Tanyakan Produk
      </h3>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
          <input
            {...register("customer_name")}
            type="text"
            className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
            placeholder="Contoh: John Doe"
          />
          {errors.customer_name && (
            <p className="mt-1 text-xs text-red-500">{errors.customer_name.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            {...register("customer_email")}
            type="email"
            className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
            placeholder="johndoe@email.com"
          />
          {errors.customer_email && (
            <p className="mt-1 text-xs text-red-500">{errors.customer_email.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Nomor WhatsApp</label>
          <input
            {...register("customer_phone")}
            type="tel"
            className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
            placeholder="081234567890"
          />
          {errors.customer_phone && (
            <p className="mt-1 text-xs text-red-500">{errors.customer_phone.message}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Pesan</label>
          <textarea
            {...register("message")}
            rows={4}
            className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all resize-none"
          />
          {errors.message && (
            <p className="mt-1 text-xs text-red-500">{errors.message.message}</p>
          )}
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded-lg transition-colors disabled:opacity-50"
        >
          {isSubmitting ? (
            <Loader2 className="h-5 w-5 animate-spin" />
          ) : (
            <>
              <Send className="h-5 w-5" /> Kirim Pertanyaan
            </>
          )}
        </button>
      </form>
    </div>
  );
}
