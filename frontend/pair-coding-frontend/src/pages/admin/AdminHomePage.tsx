import { Link } from 'react-router-dom'

const problems = [
  'Two Sum',
  'Longest Substring',
  'Merge Intervals',
  'Binary Search',
]

function AdminHomePage() {
  return (
    <section className="rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-lg shadow-slate-950/40">
      <div className="mb-6 flex items-center justify-between gap-4">
        <div>
          <p className="mb-2 text-sm uppercase tracking-[0.2em] text-cyan-400">Admin</p>
          <h1 className="text-3xl font-bold text-white">Problem list</h1>
        </div>

        <Link
          to="/admin/createnew"
          className="rounded-lg bg-cyan-500 px-4 py-2.5 font-semibold text-slate-950 transition hover:bg-cyan-400"
        >
          Create New Problem
        </Link>
      </div>

      <div className="space-y-3">
        {problems.map((problem) => (
          <div key={problem} className="rounded-xl border border-slate-800 bg-slate-950 px-4 py-3 text-slate-200">
            {problem}
          </div>
        ))}
      </div>
    </section>
  )
}

export default AdminHomePage
