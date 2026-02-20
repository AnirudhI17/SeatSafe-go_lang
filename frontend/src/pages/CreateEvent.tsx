import { type FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Button } from '../components/Button'
import { useCreateEvent } from '../api/events'

export function CreateEventPage() {
  const navigate = useNavigate()
  const createMutation = useCreateEvent()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [location, setLocation] = useState('')
  const [startsAt, setStartsAt] = useState('')
  const [endsAt, setEndsAt] = useState('')
  const [capacity, setCapacity] = useState(50)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    
    // Validate dates
    const start = new Date(startsAt)
    const end = new Date(endsAt)
    
    if (end <= start) {
      setError('End date must be after start date')
      return
    }
    
    try {
      const payload = {
        title,
        description,
        location,
        is_online: false,
        online_url: '',
        starts_at: start.toISOString(),
        ends_at: end.toISOString(),
        capacity,
        price_cents: 0,
        banner_url: '',
      }
      const created = await createMutation.mutateAsync(payload)
      navigate(`/events/${created.id}`)
    } catch (err) {
      console.error(err)
      setError('Failed to create event. Please check all fields and try again.')
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      className="mx-auto w-full max-w-2xl space-y-6 rounded-2xl border border-slate-800/70 bg-slate-900/70 p-6 shadow-xl"
    >
      <header className="space-y-1">
        <h1 className="text-xl font-semibold text-slate-50">Create event</h1>
        <p className="text-xs text-slate-400">
          Publish a new event and SeatSafe will handle concurrency-safe registrations.
        </p>
      </header>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="title">
            Title
          </label>
          <input
            id="title"
            required
            minLength={3}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="description">
            Description
          </label>
          <textarea
            id="description"
            rows={4}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="location">
            Location
          </label>
          <input
            id="location"
            required
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-1 text-sm">
            <label className="block text-slate-200" htmlFor="starts_at">
              Starts at
            </label>
            <input
              id="starts_at"
              type="datetime-local"
              required
              value={startsAt}
              onChange={(e) => setStartsAt(e.target.value)}
              className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
          <div className="space-y-1 text-sm">
            <label className="block text-slate-200" htmlFor="ends_at">
              Ends at
            </label>
            <input
              id="ends_at"
              type="datetime-local"
              required
              value={endsAt}
              onChange={(e) => setEndsAt(e.target.value)}
              className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
        </div>
        <div className="space-y-1 text-sm">
          <label className="block text-slate-200" htmlFor="capacity">
            Capacity
          </label>
          <input
            id="capacity"
            type="number"
            min={1}
            required
            value={capacity}
            onChange={(e) => setCapacity(Number(e.target.value) || 1)}
            className="w-32 rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        {error && (
          <div className="rounded-lg bg-rose-500/10 border border-rose-500/20 px-4 py-3 text-sm text-rose-300">
            {error}
          </div>
        )}
        <Button type="submit" loading={createMutation.isPending} className="w-full">
          Create event
        </Button>
      </form>
    </motion.div>
  )
}

