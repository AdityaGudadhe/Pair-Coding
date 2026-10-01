import { useState } from 'react'

const initialForm = {
  title: '',
  problemBody: '',
  timeLimit: '',
  memoryLimit: '',
  inputBody: '',
  outputBody: '',
}

type UploadMode = 'checker' | 'output' | null

function CreateNewProblemPage() {
  const [form, setForm] = useState(initialForm)
  const [uploadMode, setUploadMode] = useState<UploadMode>(null)
  const [uploadedFiles, setUploadedFiles] = useState<File[]>([])

  const handleChange = (field: keyof typeof initialForm, value: string) => {
    setForm((previous) => ({ ...previous, [field]: value }))
  }

  const handleFiles = (files: FileList | null) => {
    if (!files) return
    setUploadedFiles((previous) => [...previous, ...Array.from(files)])
  }

  return (
    <section className="rounded-2xl border border-slate-800 bg-slate-900 p-6 shadow-lg shadow-slate-950/40">
      <p className="mb-3 text-sm uppercase tracking-[0.2em] text-cyan-400">Admin</p>
      <h1 className="text-3xl font-bold text-white">Create new problem</h1>

      <div className="mt-6 grid gap-5">
        <label className="grid gap-2 text-sm text-slate-200">
          <span>Title</span>
          <input
            value={form.title}
            onChange={(event) => handleChange('title', event.target.value)}
            className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </label>

        <label className="grid gap-2 text-sm text-slate-200">
          <span>Problem body</span>
          <textarea
            value={form.problemBody}
            onChange={(event) => handleChange('problemBody', event.target.value)}
            rows={5}
            className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </label>

        <div className="grid gap-5 md:grid-cols-2">
          <label className="grid gap-2 text-sm text-slate-200">
            <span>Time limit</span>
            <input
              value={form.timeLimit}
              onChange={(event) => handleChange('timeLimit', event.target.value)}
              className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
            />
          </label>

          <label className="grid gap-2 text-sm text-slate-200">
            <span>Memory limit</span>
            <input
              value={form.memoryLimit}
              onChange={(event) => handleChange('memoryLimit', event.target.value)}
              className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
            />
          </label>
        </div>

        <label className="grid gap-2 text-sm text-slate-200">
          <span>Input body</span>
          <textarea
            value={form.inputBody}
            onChange={(event) => handleChange('inputBody', event.target.value)}
            rows={5}
            className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </label>

        <label className="grid gap-2 text-sm text-slate-200">
          <span>Output body</span>
          <textarea
            value={form.outputBody}
            onChange={(event) => handleChange('outputBody', event.target.value)}
            rows={5}
            className="rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100 outline-none transition focus:border-cyan-400"
          />
        </label>

        <div className="rounded-xl border border-slate-800 bg-slate-950 p-4">
          <p className="mb-3 text-sm font-medium text-slate-200">Choose validation type</p>
          <div className="flex gap-3">
            <button
              type="button"
              onClick={() => setUploadMode('checker')}
              className={`rounded-lg px-4 py-2 font-medium transition ${
                uploadMode === 'checker'
                  ? 'bg-cyan-500 text-slate-950'
                  : 'border border-slate-700 bg-slate-900 text-slate-100 hover:border-slate-500'
              }`}
            >
              Checker
            </button>
            <button
              type="button"
              onClick={() => setUploadMode('output')}
              className={`rounded-lg px-4 py-2 font-medium transition ${
                uploadMode === 'output'
                  ? 'bg-cyan-500 text-slate-950'
                  : 'border border-slate-700 bg-slate-900 text-slate-100 hover:border-slate-500'
              }`}
            >
              Output
            </button>
          </div>

          {uploadMode && (
            <div className="mt-4">
              <label className="block rounded-lg border border-dashed border-slate-700 bg-slate-900 p-4 text-sm text-slate-300">
                <span className="mb-2 block font-medium text-slate-200">
                  Upload {uploadMode === 'checker' ? 'checker' : 'output'} file(s)
                </span>
                <input
                  type="file"
                  multiple
                  onChange={(event) => handleFiles(event.target.files)}
                  className="block w-full text-sm text-slate-300 file:mr-4 file:rounded-md file:border-0 file:bg-cyan-500 file:px-3 file:py-2 file:font-medium file:text-slate-950"
                />
              </label>

              {uploadedFiles.length > 0 && (
                <ul className="mt-3 space-y-2 text-sm text-slate-300">
                  {uploadedFiles.map((file, index) => (
                    <li key={`${file.name}-${index}`} className="rounded-md border border-slate-800 bg-slate-950 px-3 py-2">
                      {file.name}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}
        </div>

        <button
          type="button"
          className="w-full rounded-lg bg-cyan-500 px-4 py-2.5 font-semibold text-slate-950 transition hover:bg-cyan-400"
        >
          Save Problem
        </button>
      </div>
    </section>
  )
}

export default CreateNewProblemPage
