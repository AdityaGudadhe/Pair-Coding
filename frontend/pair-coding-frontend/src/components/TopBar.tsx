import { NavLink } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'

function TopBar() {
  const isLoggedIn = useAuthStore((state) => state.isLoggedIn)

  return (
    <header className="border-b border-slate-800 bg-slate-900/90 backdrop-blur-sm">
      <div className="flex w-full items-center justify-between px-4 py-3 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3">
          <NavLink to="/" className="flex items-center gap-3 rounded-md px-1 py-1.5 transition hover:bg-slate-800">
            <div className="flex h-8 w-8 items-center justify-center rounded-md border border-slate-700 bg-slate-800 text-xs font-bold text-cyan-300">
              FC
            </div>
            <span className="text-base font-semibold text-white">FodeCorces</span>
          </NavLink>

          <NavLink
            to="/problems"
            className={({ isActive }) =>
              `rounded-md px-3 py-2 text-sm font-medium transition ${
                isActive ? 'bg-slate-800 text-white' : 'text-slate-300 hover:bg-slate-800 hover:text-white'
              }`
            }
          >
            Problems
          </NavLink>
        </div>

        <NavLink
          to={isLoggedIn ? '/account' : '/login'}
          className={({ isActive }) =>
            `inline-flex h-10 w-10 items-center justify-center rounded-full border text-sm font-semibold transition ${
              isActive
                ? 'border-cyan-400 bg-cyan-500/10 text-cyan-300'
                : 'border-slate-700 bg-slate-800 text-slate-200 hover:border-slate-500 hover:text-white'
            }`
          }
          aria-label="Account"
          title="Account"
        >
          A
        </NavLink>
      </div>
    </header>
  )
}

export default TopBar
