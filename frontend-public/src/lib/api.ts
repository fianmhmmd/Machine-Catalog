import axios from "axios";
import { Product, Category, PaginatedResponse } from "@/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export async function getLatestProducts(limit = 6): Promise<Product[]> {
  try {
    const res = await axios.get<PaginatedResponse<Product>>(`${API_URL}/products?limit=${limit}`);
    return res.data.data;
  } catch (error) {
    console.error("Error fetching latest products:", error);
    return [];
  }
}

export async function getCategories(): Promise<Category[]> {
  try {
    const res = await axios.get<Category[]>(`${API_URL}/categories`);
    return res.data;
  } catch (error) {
    console.error("Error fetching categories:", error);
    return [];
  }
}

export async function getProducts(params: {
  category?: string;
  search?: string;
  page?: number;
  limit?: number;
}): Promise<PaginatedResponse<Product>> {
  try {
    const res = await axios.get<PaginatedResponse<Product>>(`${API_URL}/products`, { params });
    return res.data;
  } catch (error) {
    console.error("Error fetching products:", error);
    return { data: [], meta: { total: 0, page: 1, limit: 10 } };
  }
}
