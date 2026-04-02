export interface Deal {
  id: string;
  created_at: string;
  status: string;
  property: Property;
  deck_type: string;
  brand: Brand;
  analysis?: FinancialAnalysis;
  deck?: Deck;
  comps?: Comp[];
  market_data?: MarketData;
  review?: Review;
  figma_file_key?: string;
  figma_file_url?: string;
}

export interface Property {
  name: string;
  address: Address;
  asset_class: string;
  units?: number;
  sq_ft?: number;
  year_built?: number;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip: string;
}

export interface Brand {
  id: string;
  name: string;
  primary_color: string;
  secondary_color: string;
  accent_color: string;
}

export interface FinancialAnalysis {
  noi: number;
  occupancy_rate: number;
  expense_ratio: number;
  avg_monthly_rent: number;
  total_units: number;
  occupied_units: number;
}

export interface Deck {
  html: string;
  generated_at: string;
  sections: Section[];
}

export interface Section {
  type: string;
  title: string;
  content: string;
}

export interface Comp {
  address: Address;
  sale_date: string;
  sale_price: number;
  cap_rate: number;
}

export interface MarketData {
  population: number;
  median_income: number;
  vacancy_rate: number;
}

export interface Review {
  id: string;
  deal_id: string;
  reviewer_id: string;
  status: string;
  started_at?: string;
  completed_at?: string;
  notes?: string;
}
