function AccountPage() {
  const user = {
    name: 'John Doe',
    email: 'john@example.com',
    password: '********',
  }

  return (
    <section className="mx-auto max-w-xl rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-lg shadow-slate-950/40">
      <p className="mb-4 text-sm uppercase tracking-[0.2em] text-cyan-400">Account</p>
      <h1 className="text-3xl font-bold text-white">Your account</h1>

      <div className="mt-6 space-y-4 rounded-xl border border-slate-800 bg-slate-950 p-4">
        <div>
          <p className="text-sm text-slate-400">Name</p>
          <p className="text-lg font-medium text-slate-100">{user.name}</p>
        </div>

        <div>
          <p className="text-sm text-slate-400">Email</p>
          <p className="text-lg font-medium text-slate-100">{user.email}</p>
        </div>

        <div>
          <p className="text-sm text-slate-400">Password</p>
          <p className="text-lg font-medium text-slate-100">{user.password}</p>
        </div>
      </div>
    </section>
  )
}

export default AccountPage
