export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-red-500">NOVIX</h1>
          <p className="text-gray-400 mt-1 text-sm">Where stories come alive.</p>
        </div>
        {children}
      </div>
    </div>
  );
}
