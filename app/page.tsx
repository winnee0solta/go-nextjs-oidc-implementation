export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)] bg-black text-white">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <a
            href="http://localhost:8080/oidc/login"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-full border border-white border-solid transition-colors flex items-center justify-center bg-white text-black gap-2 hover:bg-gray-200 font-medium text-sm sm:text-base h-10 sm:h-12 px-6 sm:px-8 sm:w-auto select-none"
          >
            Login with OIDC
          </a>
        </div>
      </main>
    </div>
  );
}
