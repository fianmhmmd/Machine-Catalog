import Link from "next/link";
import { Hammer, Mail, Phone, MapPin } from "lucide-react";

export default function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-300">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="col-span-1 md:col-span-2">
            <Link href="/" className="flex items-center gap-2 text-white">
              <Hammer className="h-8 w-8 text-blue-500" />
              <span className="text-xl font-bold">Machine Katalog</span>
            </Link>
            <p className="mt-4 text-sm max-w-md">
              Solusi terpercaya untuk kebutuhan mesin manufaktur berkualitas tinggi. 
              Kami menyediakan berbagai jenis mesin industri untuk menunjang produktivitas bisnis Anda.
            </p>
          </div>

          <div>
            <h3 className="text-white font-semibold mb-4 text-lg">Navigasi</h3>
            <ul className="space-y-2">
              <li><Link href="/" className="hover:text-white transition-colors">Home</Link></li>
              <li><Link href="/katalog" className="hover:text-white transition-colors">Katalog</Link></li>
            </ul>
          </div>

          <div>
            <h3 className="text-white font-semibold mb-4 text-lg">Kontak</h3>
            <ul className="space-y-3 text-sm">
              <li className="flex items-center gap-2">
                <MapPin className="h-4 w-4 text-blue-500" />
                Jl. Industri No. 123, Jakarta
              </li>
              <li className="flex items-center gap-2">
                <Phone className="h-4 w-4 text-blue-500" />
                +62 812-3456-7890
              </li>
              <li className="flex items-center gap-2">
                <Mail className="h-4 w-4 text-blue-500" />
                info@machinekatalog.id
              </li>
            </ul>
          </div>
        </div>
        
        <div className="mt-12 pt-8 border-t border-gray-800 text-center text-sm">
          <p>© {new Date().getFullYear()} Machine Katalog. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
