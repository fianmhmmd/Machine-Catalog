export interface Category {
  id: string;
  name: string;
  slug: string;
}

export interface ProductImage {
  id: string;
  product_id: string;
  image_url: string;
  is_primary: boolean;
  sort_order: number;
}

export interface Product {
  id: string;
  category_id: string;
  category: Category;
  name: string;
  slug: string;
  description: string;
  specifications: Record<string, any>;
  contact_phone: string;
  contact_name: string;
  is_published: boolean;
  images: ProductImage[];
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    page: number;
    limit: number;
  };
}
