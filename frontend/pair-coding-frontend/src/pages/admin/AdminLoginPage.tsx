import { NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'

function AdminLoginPage() {
  const navigate = useNavigate()
  const adminLogin = useAuthStore((state) => state.adminLogin)

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    adminLogin()
    navigate('/admin/home')
  }

  return (
    <section className="mx-auto max-w-md rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-lg shadow-slate-950/40">
      <p className="mb-4 text-sm uppercase tracking-[0.2em] text-cyan-400">Admin</p>
      <h1 className="text-3xl font-bold text-white">Admin log in</h1>

      <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
        <div>
          <label className="mb-2 block text-sm font-medium text-slate-200">Email</label>
          <input
            type="email"
            placeholder="admin@example.com"
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-slate-200">Password</label>
          <input
            type="password"
            placeholder="••••••••"
            className="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </div>

        <button
          type="submit"
          className="w-full rounded-lg bg-cyan-500 px-4 py-2.5 font-semibold text-slate-950 transition hover:bg-cyan-400"
        >
          Log in as admin
        </button>
      </form>

      <p className="mt-5 text-center text-sm text-slate-300">
        <span>Return to app? </span>
        <NavLink to="/login" className="font-medium text-cyan-400 transition hover:text-cyan-300">
          User login
        </NavLink>
      </p>
    </section>
  )
}

export default AdminLoginPage
