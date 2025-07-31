import logo from "./assets/logo-1.svg"

function App() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <img src={logo} alt="Spooler Logo" className="h-20 mb-6" />
      <p className="text-lg text-gray-700 mb-6 max-w-xl text-center">
        Spooler is a web platform for managing and submitting 3D print jobs.
      </p>
      <div className="flex flex-col sm:flex-row gap-4 mb-6">
        <a
          href="/login"
          className="bg-spooler-orange hover:bg-spooler-orange-light text-white px-6 py-2 rounded text-center transition"
        >
          Login
        </a>
        <a
          href="/register"
          className="bg-white border border-spooler-orange text-spooler-orange hover:bg-spooler-orange hover:text-white px-6 py-2 rounded text-center transition"
        >
          Register
        </a>
        <a
          href="https://github.com/North-Hall-High-School-Engineering/spooler"
          target="_blank"
          rel="noopener noreferrer"
          className="bg-gray-800 hover:bg-gray-900 text-white px-6 py-2 rounded text-center transition"
        >
          GitHub
        </a>
      </div>
      <p className="text-sm text-gray-500">
        &copy; {new Date().getFullYear()} Spooler Project
      </p>
    </div>
  )
}

export default App;