import { redirect } from 'next/navigation';

export default function Home() {
  // Automatically redirect all traffic from "/" to "/login"
  redirect('/login');
}