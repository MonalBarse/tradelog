'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { LogOut, RefreshCw } from 'lucide-react';
import api from '@/lib/api';
import { Trade, PortfolioItem } from '@/types';
import { TradeDialog } from '@/components/trade-dialog';
import { PromoteDialog } from '@/components/promote-dialog';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useAtomValue } from 'jotai';
import { authState } from '@/store/auth';
import { ShieldCheck } from 'lucide-react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'; // You need to add 'table' component

export default function DashboardPage() {
  const router = useRouter();
  const [portfolio, setPortfolio] = useState<PortfolioItem[]>([]);
  const [trades, setTrades] = useState<Trade[]>([]);
  const [loading, setLoading] = useState(true);

  const auth = useAtomValue(authState);

  const fetchData = async () => {
    setLoading(true);
    try {
      // Fetch both in parallel for speed
      const [portfolioRes, tradesRes] = await Promise.all([
        api.get('/portfolio'),
        api.get('/trades'),
      ]);
      setPortfolio(portfolioRes.data.data || []);
      setTrades(tradesRes.data.data || []);
    } catch (error) {
      console.error('Failed to fetch dashboard data', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleLogout = async () => {
    try {
      await api.post('/auth/logout');
      localStorage.removeItem('access_token');
      router.push('/login');
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* Top Bar */}
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Trader Dashboard</h1>
          <p className="text-gray-500">Welcome back, Trader.</p>
        </div>
        <div className="flex gap-2">
          {auth.user?.role === 'admin' ? (
            <Button
              className="bg-red-600 hover:bg-red-700 text-white border border-red-800"
              onClick={() => router.push('/admin')}
            >
              <ShieldCheck className="mr-2 h-4 w-4" />
              Admin Mode
            </Button>
          ) : (
            <PromoteDialog />
          )}
          <Button variant="outline" onClick={fetchData}>
            <RefreshCw
              className={`mr-2 h-4 w-4 ${loading ? 'animate-spin' : ''}`}
            />{' '}
            Refresh
          </Button>
          <Button variant="destructive" onClick={handleLogout}>
            <LogOut className="mr-2 h-4 w-4" /> Logout
          </Button>
        </div>
      </div>

      <div className="grid gap-8 md:grid-cols-3">
        {/* Left Column: Portfolio & Actions */}
        <div className="space-y-8 md:col-span-1">
          {/* Action Station */}
          <Card className="bg-slate-900 text-white border-none">
            <CardHeader>
              <CardTitle className="text-lg">Quick Actions</CardTitle>
            </CardHeader>
            <CardContent>
              <TradeDialog onSuccess={fetchData} />
            </CardContent>
          </Card>

          {/* Portfolio Summary */}
          <Card>
            <CardHeader>
              <CardTitle>Your Holdings</CardTitle>
            </CardHeader>
            <CardContent>
              {portfolio.length === 0 ? (
                <p className="text-sm text-gray-500">No assets owned.</p>
              ) : (
                <div className="space-y-4">
                  {portfolio.map((item) => (
                    <div
                      key={item.symbol}
                      className="flex items-center justify-between border-b pb-2 last:border-0"
                    >
                      <div>
                        <p className="font-bold">{item.symbol}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-mono text-lg">{item.quantity}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Right Column: Ledger */}
        <div className="md:col-span-2">
          <Card className="h-full">
            <CardHeader>
              <CardTitle>Trade History</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Time</TableHead>
                    <TableHead>Symbol</TableHead>
                    <TableHead>Side</TableHead>
                    <TableHead className="text-right">Price</TableHead>
                    <TableHead className="text-right">Qty</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {trades.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-center text-gray-500"
                      >
                        No trades recorded yet.
                      </TableCell>
                    </TableRow>
                  ) : (
                    trades.map((trade) => (
                      <TableRow key={trade.id}>
                        <TableCell className="text-xs text-gray-500">
                          {new Date(trade.executed_at).toLocaleString()}
                        </TableCell>
                        <TableCell className="font-medium">
                          {trade.symbol}
                        </TableCell>
                        <TableCell>
                          <span
                            className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                              trade.type === 'BUY'
                                ? 'bg-green-100 text-green-700'
                                : 'bg-red-100 text-red-700'
                            }`}
                          >
                            {trade.type}
                          </span>
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          ${trade.price}
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          {trade.quantity}
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
