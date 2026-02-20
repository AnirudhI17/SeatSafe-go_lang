import { motion } from 'framer-motion'
import { useEvents } from '../api/events'
import { EventCard } from '../components/EventCard'
import { EventCardSkeletonGrid } from '../components/LoadingSkeleton'

export function HomePage() {
  const { data, isLoading, isError } = useEvents()

  if (isLoading) {
    return <EventCardSkeletonGrid />
  }

  if (isError) {
    return (
      <div className="flex h-full flex-1 items-center justify-center">
        <div className="text-center">
          <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-full bg-rose-500/10 text-rose-400">
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-white mb-2">Failed to load events</h3>
          <p className="text-sm text-slate-400">Please try again later</p>
        </div>
      </div>
    )
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex h-full flex-1 flex-col items-center justify-center text-center">
        <div className="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-full bg-slate-800/50">
          <svg className="h-8 w-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <h2 className="mb-2 text-xl font-semibold text-white">
          No events yet
        </h2>
        <p className="max-w-md text-sm text-slate-400">
          Once organisers publish events, they'll appear here for attendees
          to discover and book seats in real time.
        </p>
      </div>
    )
  }

  return (
    <div className="w-full space-y-8">
      {/* Header */}
      <div className="space-y-2">
        <h1 className="text-3xl font-bold text-white tracking-tight">
          Featured events
        </h1>
        <p className="text-slate-400 max-w-2xl">
          Book seats with confidence — SeatSafe prevents overbooking even under heavy load.
        </p>
      </div>

      {/* Events grid */}
      <motion.div
        initial="hidden"
        animate="visible"
        variants={{
          visible: {
            transition: {
              staggerChildren: 0.05
            }
          }
        }}
        className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
      >
        {data.map((event) => (
          <EventCard key={event.id} {...event} />
        ))}
      </motion.div>
    </div>
  )
}

