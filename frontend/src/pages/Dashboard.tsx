import { motion } from 'framer-motion'
import { useMyRegistrations, useMyTickets } from '../api/registrations'
import { Skeleton } from '../components/LoadingSkeleton'

export function DashboardPage() {
  const { data: regs, isLoading: regsLoading } = useMyRegistrations()
  const { data: tickets, isLoading: ticketsLoading } = useMyTickets()

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex w-full flex-col gap-8"
    >
      <section>
        <h2 className="mb-3 text-lg font-semibold text-slate-50">
          My registrations
        </h2>
        {regsLoading ? (
          <Skeleton className="h-24 w-full" />
        ) : !regs || regs.length === 0 ? (
          <p className="text-sm text-slate-400">
            You haven&apos;t registered for any events yet.
          </p>
        ) : (
          <ul className="space-y-2 text-sm text-slate-200">
            {regs.map((reg) => (
              <li
                key={reg.id}
                className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-900/70 px-3 py-2"
              >
                <span>
                  Event {reg.event_id.slice(0, 8)} • {reg.quantity} ticket(s)
                </span>
                <span className="text-xs uppercase text-slate-400">
                  {reg.status}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h2 className="mb-3 text-lg font-semibold text-slate-50">My tickets</h2>
        {ticketsLoading ? (
          <Skeleton className="h-24 w-full" />
        ) : !tickets || tickets.length === 0 ? (
          <p className="text-sm text-slate-400">
            Your tickets will appear here after successful registrations.
          </p>
        ) : (
          <ul className="space-y-2 text-sm text-slate-200">
            {tickets.map((ticket) => (
              <li
                key={ticket.id}
                className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-900/70 px-3 py-2"
              >
                <span>Ticket {ticket.ticket_code}</span>
                <span className="text-xs text-slate-400">
                  Event {ticket.event_id.slice(0, 8)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>
    </motion.div>
  )
}

