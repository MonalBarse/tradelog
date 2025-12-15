"use client"

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Shield, ShieldAlert, Lock } from "lucide-react"
import { useRouter } from "next/navigation"
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

// Schema
const promoteSchema = z.object({
  secret: z.string().min(1, "Secret key is required"),
})

export function PromoteDialog() {
  const [open, setOpen] = useState(false)
  const router = useRouter()

  const form = useForm<z.infer<typeof promoteSchema>>({
    resolver: zodResolver(promoteSchema),
    defaultValues: { secret: "" },
  })

  async function onSubmit(data: z.infer<typeof promoteSchema>) {
    try {
      await api.post("/auth/promote", data)

      toast.success("Privileges Escalated", {
        description: "You are now an Admin. Please log in again to access the tower.",
        duration: 4000,
      })

      // Force Logout to get new Token with "admin" role
      setOpen(false)
      setTimeout(async () => {
        await api.post("/auth/logout")
        localStorage.removeItem("access_token")
        router.push("/login")
      }, 1500)

    } catch (error: any) {
      toast.error("Access Denied", {
        description: error.response?.data?.error || "Invalid Administration Secret",
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon" className="text-gray-400 hover:text-gray-900 hover:bg-gray-100">
          <Shield className="h-5 w-5" />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-red-600">
            <ShieldAlert className="h-5 w-5" />
            Elevate Privileges
          </DialogTitle>
          <DialogDescription>
            Enter the root secret to enable God Mode.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">

            <FormField
              control={form.control}
              name="secret"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Admin Secret Key</FormLabel>
                  <FormControl>
                    <div className="relative">
                        <Lock className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
                        <Input type="password" placeholder="••••••••••••" className="pl-9" {...field} />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" variant="destructive" className="w-full font-bold">
               Verify & Promote
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}