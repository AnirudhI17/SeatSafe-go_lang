import { useParams } from 'react-router-dom'
import { useState } from 'react'
import { motion } from 'framer-motion'
import { useEvent } from '../api/events'
import { useRegisterForEvent } from '../api/registrations'
import { Button } from '../components/Button'
import { Skeleton } from '../components/LoadingSkeleton'
import { useAuth } from '../context/AuthContext'

export function EventDetailsPage() {
  const { id } = useParams()
  const { user } = useAuth()
  const { data: event, isLoading, isError } = useEvent(id)
  const [quantity, setQuantity] = useState(1)
  const registerMutation = useRegisterForEvent(id)

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-1/3" />
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-48 w-full" />
      </div>
    )
  }

  if (isError || !event) {
    return (
      <p className="text-sm text-rose-300">
        Failed to load event details. Please try again.
      </p>
    )
  }

  const remaining = Math.max(event.capacity - event.registered_count, 0)
  const isFull = remaining === 0

  const handleRegister = () => {
    registerMutation.mutate({ quantity })
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      className="mx-auto flex w-full max-w-3xl flex-col gap-6"
    >
      <header className="space-y-3">
        <h1 className="text-3xl font-semibold tracking-tight text-slate-50">
          {event.title}
        </h1>
        <p className="text-sm text-slate-400">{event.location}</p>
        <p className="text-xs text-slate-500">
          {new Date(event.starts_at).toLocaleString()}
        </p>
      </header>

      <section className="grid gap-6 md:grid-cols-[2fr,1fr]">
        <article className="space-y-4 rounded-2xl border border-slate-800/70 bg-slate-900/60 p-5">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-400">
            About this event
          </h2>
          <p className="text-sm leading-relaxed text-slate-200">
            {event.description || 'No description provided.'}
          </p>
        </article>

        <aside className="space-y-4 rounded-2xl border border-slate-800/70 bg-slate-900/60 p-5">
          <div className="space-y-2 text-sm text-slate-300">
            <div className="flex items-center justify-between">
              <span>Capacity</span>
              <span className="font-semibold">
                {event.registered_count}/{event.capacity} booked
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span>Status</span>
              <span className="rounded-full bg-emerald-500/10 px-2 py-0.5 text-xs font-semibold text-emerald-300">
                {isFull ? 'Sold out' : 'Open'}
              </span>
            </div>
          </div>

          {user ? (
            <div className="space-y-3 border-t border-slate-800 pt-3">
              <label className="flex items-center justify-between text-sm text-slate-200">
                Tickets
                <input
                  type="number"
                  min={1}
                  max={Math.min(10, remaining)}
                  value={quantity}
                  onChange={(e) =>
                    setQuantity(Math.max(1, Math.min(10, Number(e.target.value) || 1)))
                  }
                  className="w-20 rounded-lg border border-slate-700 bg-slate-900 px-2 py-1 text-right text-sm text-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </label>
              <Button
                className="w-full"
                onClick={handleRegister}
                loading={registerMutation.isPending}
                disabled={isFull}
              >
                {isFull ? 'Event is sold out' : 'Register now'}
              </Button>
              {registerMutation.isError && (
                <p className="text-xs text-rose-300">
                  Booking failed. Please try again.
                </p>
              )}
              {registerMutation.isSuccess && (
                <p className="text-xs text-emerald-300">
                  You&apos;re registered. A ticket has been generated for you.
                </p>
              )}
            </div>
          ) : (
            <p className="border-t border-slate-800 pt-3 text-xs text-slate-400">
              Log in to register for this event.
            </p>
          )}
        </aside>
      </section>
    </motion.div>
  )
}

