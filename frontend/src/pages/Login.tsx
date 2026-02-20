import { type FormEvent, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Button } from '../components/Button'
import { useAuth } from '../context/AuthContext'

export function LoginPage() {
  const { login } = useAuth()
  const location = useLocation()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login(email, password)
    } catch (err) {
      console.error(err)
      setError('Invalid email or password.')
    } finally {
      setLoading(false)
    }
  }

  const from = (location.state as { from?: Location })?.from

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      className="mx-auto flex w-full max-w-md flex-col gap-6 rounded-2xl border border-slate-800/70 bg-slate-900/70 p-6 shadow-xl"
    >
      <header className="space-y-1">
        <h1 className="text-xl font-semibold text-slate-50">Welcome back</h1>
        <p className="text-xs text-slate-400">
          {from ? 'Log in to continue where you left off.' : 'Log in to manage your tickets and events.'}
        </p>
      </header>
      <form onSubmit={handleSubmit} className="space-y-4">
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
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
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
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        {error && <p className="text-xs text-rose-300">{error}</p>}
        <Button type="submit" className="w-full" loading={loading}>
          Log in
        </Button>
      </form>
    </motion.div>
  )
}

