"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { ShieldAlert, ArrowLeft, RefreshCw, LogOut } from "lucide-react"
import { toast } from "sonner"
import api from "@/lib/api"
import { Trade } from "@/types"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

export default function AdminDashboard() {
  const router = useRouter()
  const [trades, setTrades] = useState<Trade[]>([])
  const [loading, setLoading] = useState(true)
  const [isAdmin, setIsAdmin] = useState(false)

  const fetchSystemData = async () => {
    setLoading(true)
    try {
      const { data } = await api.get("/admin/trades")
      setTrades(data.data || [])
      setIsAdmin(true)
    } catch (error: any) {
      console.error(error)
      // If the backend says 403 Forbidden, we know they aren't an admin
      if (error.response?.status === 403) {
        toast.error("Access Denied", {
          description: "You do not have administrative privileges.",
        })
        router.push("/dashboard")
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSystemData()
  }, [])

  const handleLogout = async () => {
    try {
        await api.post("/auth/logout")
        localStorage.removeItem("access_token")
        router.push("/login")
      } catch (e) {
        console.error(e)
      }
  }

  if (!isAdmin && loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-slate-900 border-t-transparent mx-auto"></div>
          <p className="mt-4 text-slate-500">Verifying Clearance...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-slate-950 p-8 text-slate-50">
      {/* Top Bar */}
      <div className="mb-8 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="rounded-full bg-red-500/10 p-2">
            <ShieldAlert className="h-6 w-6 text-red-500" />
          </div>
          <div>
            <h1 className="text-3xl font-bold">System Watchtower</h1>
            <p className="text-slate-400">Restricted Access • Level 5 Clearance</p>
          </div>
        </div>
        <div className="flex gap-2">
           <Button variant="secondary" onClick={() => router.push("/dashboard")}>
            <ArrowLeft className="mr-2 h-4 w-4" /> Exit to Terminal
          </Button>
          <Button variant="outline" className="bg-slate-900 text-white hover:bg-slate-800 border-slate-700" onClick={fetchSystemData}>
            <RefreshCw className={`mr-2 h-4 w-4 ${loading ? 'animate-spin' : ''}`} /> Refresh
          </Button>
           <Button variant="destructive" onClick={handleLogout}>
              <LogOut className="mr-2 h-4 w-4" /> Logout
            </Button>
        </div>
      </div>

      <div className="grid gap-6">
        {/* Stats Row (Optional simple stats) */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
            <Card className="bg-slate-900 border-slate-800 text-slate-50">
                <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-medium text-slate-400">Total Volume</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{trades.length} Trades</div>
                </CardContent>
            </Card>
             {/* You could add more stats here like "Total Users" if the API supported it */}
        </div>

        {/* The Master Ledger */}
        <Card className="bg-slate-900 border-slate-800 text-slate-50">
          <CardHeader>
            <CardTitle>System-Wide Trade Ledger</CardTitle>
            <CardDescription className="text-slate-400">
              Real-time monitoring of all user activities.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border border-slate-800">
              <Table>
                <TableHeader className="bg-slate-950">
                  <TableRow className="border-slate-800 hover:bg-slate-900">
                    <TableHead className="text-slate-400">ID</TableHead>
                    <TableHead className="text-slate-400">User ID</TableHead>
                    <TableHead className="text-slate-400">Time</TableHead>
                    <TableHead className="text-slate-400">Symbol</TableHead>
                    <TableHead className="text-slate-400">Side</TableHead>
                    <TableHead className="text-right text-slate-400">Price</TableHead>
                    <TableHead className="text-right text-slate-400">Qty</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {trades.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="text-center text-slate-500 h-24">
                        System is idle. No trades recorded.
                      </TableCell>
                    </TableRow>
                  ) : (
                    trades.map((trade) => (
                      <TableRow key={trade.id} className="border-slate-800 hover:bg-slate-800/50">
                        <TableCell className="font-mono text-xs text-slate-500">{trade.id}</TableCell>
                        <TableCell>
                            <span className="inline-flex items-center rounded-md bg-slate-800 px-2 py-1 text-xs font-medium text-slate-300 ring-1 ring-inset ring-slate-700/10">
                                UID-{trade.user_id}
                            </span>
                        </TableCell>
                        <TableCell className="text-xs text-slate-400">
                          {new Date(trade.executed_at).toLocaleString()}
                        </TableCell>
                        <TableCell className="font-medium text-slate-200">{trade.symbol}</TableCell>
                        <TableCell>
                          <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-bold ${
                            trade.type === 'BUY'
                            ? 'bg-green-500/10 text-green-400'
                            : 'bg-red-500/10 text-red-400'
                          }`}>
                            {trade.type}
                          </span>
                        </TableCell>
                        <TableCell className="text-right font-mono text-slate-300">${trade.price}</TableCell>
                        <TableCell className="text-right font-mono text-slate-300">{trade.quantity}</TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}