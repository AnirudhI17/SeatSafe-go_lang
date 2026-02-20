import type { ReactNode } from 'react'
import { Navbar } from '../components/Navbar'

export function MainLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950">
      <Navbar />
      <main className="mx-auto flex w-full max-w-[1400px] flex-1 px-8 py-12">
        {children}
      </main>
      <footer className="border-t border-slate-800/50 bg-slate-950/50 py-8">
        <div className="mx-auto max-w-[1400px] px-8">
          <div className="flex items-center justify-between">
            <p className="text-sm text-slate-400">
              © 2024 SeatSafe. All rights reserved.
            </p>
            <div className="flex items-center gap-6">
              <a href="#" className="text-sm text-slate-400 hover:text-white transition-colors duration-200">
                Privacy
              </a>
              <a href="#" className="text-sm text-slate-400 hover:text-white transition-colors duration-200">
                Terms
              </a>
              <a href="#" className="text-sm text-slate-400 hover:text-white transition-colors duration-200">
                Support
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}

