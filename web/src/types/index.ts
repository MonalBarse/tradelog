export interface Trade {
  id: number;
  user_id: number;
  symbol: string;
  type: 'BUY' | 'SELL';
  price: string; // Backend sends decimals as strings or numbers, safer to handle as string for display
  quantity: string;
  executed_at: string;
}

export interface PortfolioItem {
  symbol: string;
  quantity: string;
  value: string;
}
