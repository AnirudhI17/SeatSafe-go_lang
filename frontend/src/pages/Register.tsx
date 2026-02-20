import { type FormEvent, useState } from 'react'
import { motion } from 'framer-motion'
import { Button } from '../components/Button'
import { useAuth } from '../context/AuthContext'

export function RegisterPage() {
  const { register } = useAuth()
  const [fullName, setFullName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await register({ email, password, full_name: fullName })
    } catch (err) {
      console.error(err)
      setError('Registration failed. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      className="mx-auto flex w-full max-w-md flex-col gap-6 rounded-2xl border border-slate-800/70 bg-slate-900/70 p-6 shadow-xl"
    >
      <header className="space-y-2">
        <h1 className="text-2xl font-bold text-white">
          Create your SeatSafe account
        </h1>
        <p className="text-sm text-slate-400">
          Book events in seconds and keep your tickets in one place.
        </p>
      </header>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="full_name">
            Full name
          </label>
          <input
            id="full_name"
            type="text"
            required
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-purple-500 transition-all duration-200"
          />
        </div>
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="email">
            Email
          </label>
          <input
            id="email"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-purple-500 transition-all duration-200"
          />
        </div>
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="password">
            Password
          </label>
          <input
            id="password"
            type="password"
            required
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-purple-500 transition-all duration-200"
          />
        </div>
        {error && <p className="text-xs text-rose-300">{error}</p>}
        <Button type="submit" className="w-full" loading={loading}>
          Sign up
        </Button>
      </form>
    </motion.div>
  )
}

