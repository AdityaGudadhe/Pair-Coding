import { Link } from 'react-router-dom'

function HomePage() {
  return (
    <section className="flex min-h-[70vh] items-center justify-center rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-lg shadow-slate-950/40">
      <div className="text-center">
        <p className="mb-3 text-sm uppercase tracking-[0.2em] text-cyan-400">Home</p>
        <h1 className="text-3xl font-bold text-white md:text-4xl">Welcome home</h1>
        <p className="mt-4 max-w-2xl text-slate-300">
          This is just a placeholder landing page for now. Add content here later.
        </p>

        <Link
          to="/problems"
          className="mt-8 inline-flex items-center justify-center rounded-lg bg-cyan-500 px-6 py-3 font-semibold text-slate-950 transition hover:bg-cyan-400"
        >
          Go to Problems
        </Link>
      </div>
    </section>
  )
}

export default HomePage
