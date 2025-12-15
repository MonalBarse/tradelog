"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { TrendingUp, TrendingDown, Bitcoin } from "lucide-react"
import { toast } from "sonner"
import api from "@/lib/api"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

const TRADING_PAIRS = [
  { value: "BTC/USD", label: "Bitcoin (BTC/USD)" },
  { value: "ETH/USD", label: "Ethereum (ETH/USD)" },
  { value: "SOL/USD", label: "Solana (SOL/USD)" },
  { value: "DOGE/USD", label: "Dogecoin (DOGE/USD)" },
  { value: "AAPL/USD", label: "Apple (AAPL/USD)" },
] as const

const tradeSchema = z.object({
  symbol: z.enum(["BTC/USD", "ETH/USD", "SOL/USD", "DOGE/USD", "AAPL/USD"]),
  type: z.enum(["BUY", "SELL"]),
  price: z.coerce.number().positive({ message: "Price must be positive" }),
  quantity: z.coerce.number().positive({ message: "Quantity must be positive" }),
})

type TradeFormValues = z.infer<typeof tradeSchema>

export function TradeDialog({ onSuccess }: { onSuccess: () => void }) {
  const [open, setOpen] = useState(false)

  // No generic type passed to useForm to prevent Resolver conflicts
  const form = useForm({
    resolver: zodResolver(tradeSchema),
    defaultValues: {
      symbol: "BTC/USD",
      type: "BUY",
      price: 0,
      quantity: 0,
    },
  })

  async function onSubmit(data: TradeFormValues) {
    try {
      await api.post("/trades", data)
      toast.success("Trade Executed", {
        description: `${data.type} ${data.quantity} ${data.symbol} @ $${data.price}`,
      })
      setOpen(false)
      form.reset()
      onSuccess()
    } catch (error: any) {
      toast.error("Trade Failed", {
        description: error.response?.data?.error || "Could not execute trade",
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="lg" className="w-full bg-blue-600 hover:bg-blue-700 shadow-lg shadow-blue-900/20">
          <Bitcoin className="mr-2 h-4 w-4" />
          Execute Trade
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>New Market Order</DialogTitle>
          <DialogDescription>
            Select an asset and execute a trade immediately.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pt-2">

            <div className="grid grid-cols-1 gap-4">
              <FormField
                control={form.control}
                name="symbol"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Asset Pair</FormLabel>
                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select a pair" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {TRADING_PAIRS.map((pair) => (
                          <SelectItem key={pair.value} value={pair.value}>
                            {pair.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Action</FormLabel>
                    <div className="flex gap-2">
                        <div
                          className={`flex-1 cursor-pointer rounded-md border p-3 text-center font-bold transition-all ${
                            field.value === 'BUY'
                              ? 'border-green-500 bg-green-50 text-green-700 ring-1 ring-green-500'
                              : 'border-gray-200 hover:bg-gray-50'
                          }`}
                          onClick={() => field.onChange('BUY')}
                        >
                          BUY
                        </div>
                        <div
                          className={`flex-1 cursor-pointer rounded-md border p-3 text-center font-bold transition-all ${
                            field.value === 'SELL'
                              ? 'border-red-500 bg-red-50 text-red-700 ring-1 ring-red-500'
                              : 'border-gray-200 hover:bg-gray-50'
                          }`}
                          onClick={() => field.onChange('SELL')}
                        >
                          SELL
                        </div>
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="price"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Price ($)</FormLabel>
                    <FormControl>
                      {/* FIXED: Added explicit cast to 'any' to satisfy strict Typescript build */}
                      <Input
                        type="number"
                        step="0.01"
                        {...field}
                        value={field.value as any}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="quantity"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Quantity</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        step="0.000001"
                        {...field}
                        value={field.value as any}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="rounded-lg bg-slate-100 p-3 text-sm text-slate-600 flex justify-between">
                <span>Total Estimate:</span>
                <span className="font-mono font-bold text-slate-900">
                    ${(Number(form.watch("price") || 0) * Number(form.watch("quantity") || 0)).toLocaleString()}
                </span>
            </div>

            <Button
              type="submit"
              className={`w-full font-bold ${
                form.watch("type") === "BUY" ? "bg-green-600 hover:bg-green-700" : "bg-red-600 hover:bg-red-700"
              }`}
            >
               {form.watch("type") === "BUY" ? <TrendingUp className="mr-2 h-4 w-4" /> : <TrendingDown className="mr-2 h-4 w-4" />}
               Confirm {form.watch("type")}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}